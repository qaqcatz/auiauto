package panalyze

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/psrctree"
	"auiauto/pkernel/pstmtlog"
	"auiauto/pkernel/psuscode"
	"auiauto/putils"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

// 普通分析, 在指定目录下生成分析文件analyze_casePrefix_analyzeFile_factor.txt
// 格式:
// 第一行一个数, 表示聚焦语句数量s
// 接下来一行s个数, 代表聚焦语句文件中每个聚焦语句的排名

// 根据projectId/src创建源码树, 读取所有以casePrefix开头的testcases下的cover.json进行源码树覆盖, 通过factor进行差异分析, 根据analyzeFile进行分析结果过滤
// 将分析结果保存在projectId/testcases/的analyze_casePrefix_analyzeFile_factor.txt
func ReadSourceTreeAndNormalAnalyzeStd(projectId string, casePrefix string, analyzeFile string, factor string) (*psrctree.SourceTree,
	psuscode.SusCodes, *perrorx.ErrorX) {
	testcasesPath := pdba.DBURLProjectIdTestcases(projectId)
	caseNames, err := putils.GetDirsStartWith(casePrefix, testcasesPath)
	if err != nil {
		return nil, nil, perrorx.TransErrorX(err)
	}
	sourceTree, suscodeSlice, err := readSourceTreeAndAnalyzeStd(projectId, caseNames, analyzeFile, factor)
	if err != nil {
		return nil, nil, perrorx.TransErrorX(err)
	}

	// 将分析结果保存在projectId/testcases/的analyze_casePrefix_analyzeFile_factor.txt
	analyzePath := path.Join(testcasesPath, "analyze_"+casePrefix+"_"+analyzeFile+"_"+factor+".txt")
	ans := strconv.Itoa(suscodeSlice.Len()) + "\n"
	for i := 0; i < suscodeSlice.Len(); i++ {
		ans += strconv.Itoa(suscodeSlice[i].MRank) + " "
	}
	ans += "\n"
	err_ := ioutil.WriteFile(analyzePath, []byte(ans), 0777)
	if err_ != nil {
		return nil, nil, perrorx.NewErrorXWriteFile(analyzePath, err_.Error(), nil)
	}
	return sourceTree, suscodeSlice, nil
}

// 根据projectId/src创建源码树, 读取projectId/testcases/caseNames列表对应的caseName下的cover.json进行源码树覆盖,
// 通过factor进行差异分析, 根据analyzeFile进行分析结果过滤
func readSourceTreeAndAnalyzeStd(projectId string, caseNames []string, analyzeFile string, factor string) (*psrctree.SourceTree,
	psuscode.SusCodes, *perrorx.ErrorX) {
	cases := make([]string, 0)
	for _, caseName := range caseNames {
		cases = append(cases, pdba.DBURLProjectIdTestcaseCover(projectId, caseName))
	}
	sourceTree, susCodeSlice, err := readSourceTreeAndAnalyze(projectId, cases, analyzeFile, factor)
	if err != nil {
		return nil, nil, perrorx.TransErrorX(err)
	}
	return sourceTree, susCodeSlice, nil
}

// 实验性功能开关, 在Sim2S时启用, 增加栈信息权重
var gSim2S = false

// 根据projectId/src创建源码树, 读取cases列表的路径进行源码树覆盖, 通过factor进行差异分析, 根据analyzeFile进行分析结果过滤
func readSourceTreeAndAnalyze(projectId string, cases []string, analyzeFile string, factor string) (*psrctree.SourceTree,
	psuscode.SusCodes, *perrorx.ErrorX) {
	sourceTree, err := psrctree.ReadSourceTree(projectId)
	if err != nil {
		return nil, nil, perrorx.TransErrorX(err)
	}
	for _, casePath := range cases {
		err = psrctree.CoverSourceTree(sourceTree, casePath, "")
		if err != nil {
			return nil, nil, perrorx.TransErrorX(err)
		}
	}
	psrctree.CalCoverNumDFS(sourceTree.MRoot)

	// calculate ranking
	susCodes := make(psuscode.SusCodes, 0)
	sourceTree.Foreach(func (cur *psrctree.SourceNode) {
		// 不是源文件(叶结点)的话忽略, 否则遍历MTotalNum会出事
		if len(cur.MChildren) != 0 {
			return
		}
		for i := 0; i < cur.MTotalNum; i++ {
			// 用cur.MPassed[i] + cur.MFailed[i] > 0判断这句代码是否被覆盖过
			if cur.MPassed[i] + cur.MFailed[i] > 0 {
				susCodes = append(susCodes, psuscode.SusCode{MOriginNode: cur, MIdx: i})
			}
		}
	})

	// 通过栈信息增强定位效果
	// 格式: 源文件的dotClassPath@line->idx(在栈中的顺序, 之统计源文件中出现过的),
	// 即去除$, 通过srctree的findSrcNode获取的src node对应的MFullName
	st := make(map[string]int, 0)
	if gSim2S {
		// 这里先定死读origin_crash
		stmtLog, err := pstmtlog.ReadStmtLogStd(projectId, "origin_crash")
		if err != nil {
			return nil, nil, perrorx.TransErrorX(err)
		}
		idx := 0
		for _, element := range stmtLog.MStackTraceStr {
			sp := strings.Split(element, "@")
			if len(sp) != 2 {
				continue
			}
			dotClassPath := sp[0]
			if strings.Contains(dotClassPath, "$") {
				dotClassPath = dotClassPath[0:strings.Index(dotClassPath, "$")]
			}
			if srcNode := sourceTree.FindSrcNode(dotClassPath); srcNode != nil {
				idx++
				st[srcNode.MFullName+"@"+sp[1]] = idx
			}
		}
	}

	psuscode.DoInit(susCodes, sourceTree.MTotalPassed, sourceTree.MTotalFailed)
	err = psuscode.DoSusFormula(factor, susCodes, st)
	if err != nil {
		return nil, nil, perrorx.TransErrorX(err)
	}
	// 注意分析完后要把ranking信息反馈给srctree
	for i := 0; i < susCodes.Len(); i++ {
		susCodes[i].MOriginNode.MRanks[susCodes[i].MIdx] = susCodes[i].MRank
	}
	susCodeSlice, err := psuscode.CreateSlice(projectId, analyzeFile, susCodes)
	if err != nil {
		return nil, nil, perrorx.TransErrorX(err)
	}

	return sourceTree, susCodeSlice, nil
}