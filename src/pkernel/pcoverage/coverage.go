package pcoverage

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/pahttp"
	"auiauto/pkernel/pcfg"
	"auiauto/pkernel/pstmtlog"
	"encoding/json"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"
)

// 覆盖率
type Coverage struct {
	// 是否崩溃
	MStatus       string       `json:"status"`
	// 覆盖的类
	MCoverClasses []CoverClass `json:"classes"`
}

// 从projectId/testcases/caseName/cover.json中读取Coverage
func ReadCoverageStd(projectId string, caseName string) (*Coverage, *perrorx.ErrorX) {
	coverPath := pdba.DBURLProjectIdTestcaseCover(projectId, caseName)
	coverage, err := ReadCoverage(coverPath)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	return coverage, nil
}

// 从path中加载Coverage
func ReadCoverage(path string) (*Coverage, *perrorx.ErrorX) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, perrorx.NewErrorXReadFile(path, err.Error(), nil)
	}
	var coverage Coverage
	err = json.Unmarshal(jsonData, &coverage)
	if err != nil {
		return nil, perrorx.NewErrorXUnmarshal(err.Error(), nil)
	}
	return &coverage, nil
}

// 其实是个单链dfs, 分支的判断交给了外层迭代
// @stmtLog: stmtLog.Crashed(
// @jsonPath: stmtLog.Crashed(jsonPath,
// @cfgMethod: cfgMethod.MStmts[jid]
func lineCoverageDFS(stmtLog *pstmtlog.StmtLog, jsonPath string, cfgMethod *pcfg.CFGMethod,
	jidMap map[int]int64, jid int, eventId int64) *perrorx.ErrorX {

	stmt := cfgMethod.MStmts[jid]

	// 记忆化, 这一步一定要先做, 不然used可能会遗漏一些语句
	newEventId := eventId
	if oldEventId, ok := jidMap[jid]; ok {
		if (oldEventId | newEventId) == oldEventId {
			return nil
		} else {
			newEventId = newEventId | oldEventId
		}
	}
	jidMap[jid] = newEventId

	// 判断语句是否崩溃
	if stmtLog.Crashed(jsonPath, stmt.MSid) {
		return nil
	}
	// 禁止走br, 因为br的逻辑是在外层迭代做的
	if stmt.MType == "br" {
		return nil
	}
	if stmt.MType == "gt" {
		// goto跳转
		err := lineCoverageDFS(stmtLog, jsonPath, cfgMethod, jidMap, stmt.MTargets[0], newEventId)
		if err != nil {
			return err
		}
	} else if stmt.MType == "n" {
		// n语句通过fallThrough判断是否要继续向下一条语句走
		if stmt.MFallThrough {
			err := lineCoverageDFS(stmtLog, jsonPath, cfgMethod, jidMap, jid+1, newEventId)
			if err != nil {
				return err
			}
		}
	} else {
		return perrorx.NewErrorXLineCoverageDFS("unknown type", nil)
	}
	return nil
}

// 计算覆盖语句
func LineCoverage(stmtLog *pstmtlog.StmtLog, projectId string) (*Coverage, *perrorx.ErrorX) {
	// class name -> lines
	classesMap := make(map[string]map[int]int64)

	// 遍历stmtLog中的每个class, 在cfg中寻找对应的class
	for classJsonPath, stmtLogMethods := range stmtLog.MClasses {
		// 将classJsonPath还原成className
		className := strings.ReplaceAll(classJsonPath[0:len(classJsonPath)-5], "/", ".")
		// test$1 -> test
		if strings.Contains(className, "$") {
			className = className[0:strings.Index(className, "$")]
		}

		var coverLinesMap map[int]int64 = nil
		exist := false
		if coverLinesMap, exist = classesMap[className]; !exist {
			coverLinesMap = make(map[int]int64)
			classesMap[className] = coverLinesMap
		}

		// 根据stmtLog的class信息加载对应文件的cfg
		cfg, err := pcfg.ReadCFG(path.Join(pdba.DBURLProjectIdCFG(projectId), classJsonPath))
		if err != nil {
			return nil, perrorx.TransErrorX(err)
		}

		// 对于stmtLog.MClasses中的每个函数, 在cfg中寻找对应的cfgMethod
		// 对于每个cfgMethod, 开启单链dfs, 根据stmtLogMethod中的信息决定每个分支点怎么走, 从而得到语句覆盖
		for _, stmtLogMethod := range *stmtLogMethods {
			cfgMethod := cfg.Search(stmtLogMethod.MethodSig)
			if cfgMethod == nil {
				return nil, perrorx.NewErrorXCalLineCoverage("smtLog's method and cfg's method do not match: " + stmtLogMethod.MethodSig, nil)
			}

			// 函数内的jid记忆化
			jidMap := make(map[int]int64)

			// 先从函数头开始走一下
			methodEventId, err := strconv.ParseInt(stmtLogMethod.MethodEventId, 10, 64)
			if err != nil {
				return nil, perrorx.NewErrorXParseInt(stmtLogMethod.MethodEventId, nil)
			}
			err_ := lineCoverageDFS(stmtLog, classJsonPath, cfgMethod, jidMap, 0, methodEventId)
			if err_ != nil {
				return nil, perrorx.TransErrorX(err_)
			}
			// 根据每个分支信息走一下
			for _, stmt := range stmtLogMethod.Stmts {
				if stmt.MType == "br" {
					eventId, err := strconv.ParseInt(stmt.MEventId, 10, 64)
					if err != nil {
						return nil, perrorx.NewErrorXParseInt(stmt.MEventId, nil)
					}
					// 一定会在开始断掉, 用来做下记忆化以及崩溃判断
					err_ := lineCoverageDFS(stmtLog, classJsonPath, cfgMethod, jidMap, stmt.MJid, eventId)
					if err_ != nil {
						return nil, perrorx.TransErrorX(err_)
					}
					// 搜索对应的分支
					// 这里cfgMethod.MStmts[stmt.MJid].MTargets[stmt.MValue]表示根据stmtLogStmt的jid获取cfgStmt,
					// 然后以stmtLogStmt的Value作为分支下标在cfgStmt.Targets中获取目标语句
					err_ = lineCoverageDFS(stmtLog, classJsonPath, cfgMethod, jidMap,
						cfgMethod.MStmts[stmt.MJid].MTargets[stmt.MValue], eventId)
					if err_ != nil {
						return nil, perrorx.TransErrorX(err_)
					}
				} else {
					return nil, perrorx.NewErrorXCalLineCoverage("unknown type", nil)
				}
			}

			// 现在used中记录了这个函数里走过的jid, 我们将它还原成sid, add it to coverLines
			for jid, eventId := range jidMap {
				sid := cfgMethod.MStmts[jid].MSid
				if sid != 0 {
					if oldEventId, exist := coverLinesMap[sid]; exist {
						coverLinesMap[sid] = oldEventId | eventId
					} else {
						coverLinesMap[sid] = eventId
					}
				} // end of if MSid != 0
			} // end of for jidMap
		} // end of for stmtLogMethod
	} // end of for stmtLog.MClasses

	// change ans to Coverage
	coverage := &Coverage{MStatus: stmtLog.MStatus, MCoverClasses: make([]CoverClass, 0)}
	for className, coverLinesMap := range classesMap {
		if len(coverLinesMap) == 0 {
			continue
		}
		coverLines := make(SortPairs, 0)
		for sid, eventId := range coverLinesMap {
			coverLines = append(coverLines, SortPair{sid, strconv.FormatInt(eventId, 10)})
		}
		// 为了用户阅读友好, 给coverLines排序
		sort.Sort(coverLines)
		tLines := make([]int, coverLines.Len())
		tEventIds := make([]string, coverLines.Len())
		for i := 0; i < coverLines.Len(); i++ {
			tLines[i] = coverLines[i].MSid
			tEventIds[i] = coverLines[i].MEventId
		}
		coverClass := CoverClass{MClassName: className, MLines: tLines, MEventIds: tEventIds}
		coverage.MCoverClasses = append(coverage.MCoverClasses, coverClass)
	}

	return coverage, nil
}

// 将覆盖语句保存在projectId/testcases/caseName/cover.json下
func SaveLineCoverageStd(projectId string, caseName string, coverage *Coverage) *perrorx.ErrorX {
	coverPath := pdba.DBURLProjectIdTestcaseCover(projectId, caseName)
	err := SaveLineCoverage(coverPath, coverage)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}

// 将覆盖语句保存在path下
func SaveLineCoverage(path string, coverage *Coverage) *perrorx.ErrorX {
	jsonData, err := json.Marshal(coverage)
	if err != nil {
		return perrorx.NewErrorXMarshal(err.Error(), nil)
	}
	err = ioutil.WriteFile(path, jsonData, 0777)
	if err != nil {
		return perrorx.NewErrorXWriteFile(path, err.Error(), nil)
	}
	return nil
}

// 从antrance获取stmtlog, 保存在projectId/testcases/caseName/stmtlog.json,
// 接着读取stmtlog.json, 调用ParseStmtLog解析出StmtLog,
// 然后使用LineCoverage计算覆盖率, 最后使用SaveLineCoverage保存在projectId/testcases/caseName/cover.json下
// 返回StmtLog, 方便做一些状态判断
func SaveStmtLogAndLineCoverageStd(avd string, projectId string, caseName string) (*pstmtlog.StmtLog, *perrorx.ErrorX) {
	dir := pdba.DBURLProjectIdTestcase(projectId, caseName)
	stmtLog, err := SaveStmtLogAndLineCoverage(avd, dir, projectId)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	return stmtLog, err
}

// 从antrance获取stmtlog, 保存在dir/stmtlog.json, 接着读取dir/stmtlog.json, 调用ParseStmtLog解析出StmtLog,
// 然后使用LineCoverage计算覆盖率, 最后使用SaveLineCoverage保存到dir/cover.json下
// 返回StmtLog, 方便做一些状态判断
func SaveStmtLogAndLineCoverage(avd string, dir string, projectId string) (*pstmtlog.StmtLog,
	*perrorx.ErrorX) {
	stmtLogPath := path.Join(dir, "stmtlog.json")
	logIdSigPath :=pdba.DBURLProjectIdLogidsig(projectId)
	coverPath := path.Join(dir, "cover.json")
	// 获取stmtlog并保存到dir/stmtlog.json下
	err := pahttp.GetStmtLogAndSave(avd, stmtLogPath)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	// 解析stmtlog.json
	stmtLog, err := pstmtlog.ParseStmtLog(stmtLogPath, logIdSigPath)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	// 校验申请方的projectId和antrance的实际projectId是否一致
	if projectId != stmtLog.MProjectId {
		return nil, perrorx.NewErrorXSaveStmtLogAndLineCoverage("projectId(" + projectId +
			") != stmtLog.MProjectId(" + stmtLog.MProjectId + ")", nil)
	}
	// 计算cover.json
	coverage, err := LineCoverage(stmtLog, projectId)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	// 保存cover.json
	err = SaveLineCoverage(coverPath, coverage)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	return stmtLog, nil
}
