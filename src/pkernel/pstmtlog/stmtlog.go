package pstmtlog

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/putils"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 程序运行中产生的语句执行日志, 通过解析antrance产生的stmtlog.json文件而来, 用于结合CFG计算覆盖率
// 另外, 由于我们用了数组标记的方式收集语句日志, 还需要根据soot产生的logIdSig.txt将数组id还原成实际语句
// stmtlog.json格式:
// { "projectId"(当前程序的项目id, 用户指定):"com.example.debugapp",
//   "status"(程序正常/崩溃):"true/false",
//   "stmts"(程序运行过程中执行的语句id):[0, 3, 8001, 10234],
//   "eventids"(语句关联的eventids, 1<<0表示init覆盖):["9","3","1","2"],
//   "stackTrace"(status为false时表示出现了uncaught exception, 需记录栈调用信息, status true时为空): [
//     "类@语句在文件中的源码行"
//   ],
//   "stackTraceOrigin": "原始栈信息"
// }
// logIdSig.txt格式:
// 语句对应的数组下标 语句sig
// 语句sig格式:
// methodSig@jid@sid@type@value"(jid:语句在函数中的字节码id, sid:语句在文件中的源码行,
// type目前只考虑了branch(br), value 0表示false/default分支, >=1表示true/各个case分支),
// 特别地, 对于每个函数的入口插桩, 只记录methodSig.
type StmtLog struct {
	// 项目id
	MProjectId        string   `json:"projectId"`
	// 本次运行是否崩溃
	MStatus           string   `json:"status"`
	// 语句数组, MStmts[i] != 0表示语句i被执行过
	MStmts            []int    `json:"stmts"`
	// 语句被哪些动作覆盖过, 二进制表示
	MEventIds         []string `json:"eventids"`
	// 按行拆分并处理后的stack trace
	MStackTraceStr    []string `json:"stackTrace"`
	// 原始的stack trace
	MStackTraceOrigin string   `json:"stackTraceOrigin"`

	// stmt根据idsigmap还原后根据methodsig归类
	MStmtLogMethods []*StmtLogMethod `json:"-"`

	// 根据class将各个StmtLogMethod分组, 为了方便与cfg对应, 这里的class直接使用项目目录下的json路径, 例如:
	// a.b.c$d->a/b/c$d.json
	MClasses map[string]*[]*StmtLogMethod `json:"-"`

	// 将stack trace中的class转换为jsonPath形式(与stmtLog.Classes同理)
	MStackTrace map[string]bool `json:"-"`
}

// 根据logIdSig, 将Stmts聚合成Methods的形式
// 根据class将各个StmtLogMethod分组, 为了方便与cfg对应, 这里的class直接使用项目目录下的json路径, 例如:
// a.b.c->a/b/c.json
// 将stack trace中的class转换为jsonPath形式(与stmtLog.Classes同理)
func (stmtLog *StmtLog) Parse(logIdSigPath string) *perrorx.ErrorX {
	// 读取logIdSig, 拆分成id:sig存储在map中
	idSigMap := make(map[int]string)
	if !putils.FileExist(logIdSigPath) {
		return perrorx.NewErrorXFileNotFound(logIdSigPath, nil)
	}
	if file, err := os.Open(logIdSigPath); err != nil {
		return perrorx.NewErrorXOpen(logIdSigPath, nil)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			idSig := scanner.Text()
			if idSig == "" {
				continue
			}
			// 注意可能有多个空格, 只分离第一个
			sp := strings.SplitN(idSig, " ", 2)
			if len(sp) != 2 {
				return perrorx.NewErrorXSplitN(len(sp), 2, nil)
			}
			id, err := strconv.Atoi(sp[0])
			if err != nil {
				return perrorx.NewErrorXAtoI(sp[0], nil)
			}
			idSigMap[id] = sp[1]
		}
	}
	// 根据idSigMap将stmtLog.Stmts解析成methodSig:*[]stmtstr的形式, 存储在methodStmtsMap中
	// 注意要考虑到eventids, 作为@后缀存储在stmtstr中
	methodStmtsMap := make(map[string]*[]string)
	// 特殊处理methodSig -> eventid
	methodEventIdMap := make(map[string]string)
	for i, id := range stmtLog.MStmts {
		if sig, exist := idSigMap[id]; exist {
			// 可能有多个@, 只分离第一个
			sp := strings.SplitN(sig, "@", 2)
			if len(sp) == 1 {
				// 进入method时只记录methodSig, 不记录语句, 因此特殊判断没有@的情况
				if _, exist = methodStmtsMap[sp[0]]; !exist {
					temp := make([]string, 0)
					methodStmtsMap[sp[0]] = &temp
				}
				methodEventIdMap[sp[0]] = stmtLog.MEventIds[i]
			} else {
				if len(sp) != 2 {
					return perrorx.NewErrorXSplitN(len(sp), 2, nil)
				}
				if stmts, exist := methodStmtsMap[sp[0]]; exist {
					*stmts = append(*stmts, sp[1]+"@"+stmtLog.MEventIds[i])
				} else {
					temp := make([]string, 0)
					temp = append(temp, sp[1]+"@"+stmtLog.MEventIds[i])
					methodStmtsMap[sp[0]] = &temp
				}
			}
		} else {
			return perrorx.NewErrorXStmtlogParse("stmtLog.MStmts has unknown id", nil)
		}
	}
	// 将methodStmtsMap解析成stmtLog.MStmtLogMethods
	stmtLog.MStmtLogMethods = make([]*StmtLogMethod, 0)
	for method, stmts := range methodStmtsMap {
		stmtMethod := StmtLogMethod{}
		stmtMethod.MethodSig = method
		stmtMethod.MethodEventId = methodEventIdMap[method]
		stmtMethod.StmtStrs = make([]string, len(*stmts))
		for i, stmt := range *stmts {
			stmtMethod.StmtStrs[i] = stmt
		}
		// 调用StmtLogMethod的parse函数将StmtStrs解析成真正的StmtLogStmt
		err_ := stmtMethod.Parse()
		if err_ != nil {
			return perrorx.TransErrorX(err_)
		}
		stmtLog.MStmtLogMethods = append(stmtLog.MStmtLogMethods, &stmtMethod)
	}

	// 根据class将各个StmtLogMethod分组, 为了方便与cfg对应, 这里的class直接使用项目目录下的json路径, 例如:
	// a.b.c$d->a/b/c$d.json
	stmtLog.MClasses = make(map[string]*[]*StmtLogMethod)
	for _, stmtMethod := range stmtLog.MStmtLogMethods {
		sig := stmtMethod.MethodSig
		className := strings.Split(sig[1:], ":")[0]
		classJsonPath := strings.ReplaceAll(className, ".", "/") + ".json"
		if _, exist := stmtLog.MClasses[classJsonPath]; !exist {
			newMethods := make([]*StmtLogMethod, 0)
			stmtLog.MClasses[classJsonPath] = &newMethods
		}
		methods := stmtLog.MClasses[classJsonPath]
		*methods = append(*methods, stmtMethod)
	}
	// 将stack trace中的class转换为classJsonPath形式(与stmtLog.Classes同理)
	stmtLog.MStackTrace = make(map[string]bool)
	for _, stackElementStr := range stmtLog.MStackTraceStr {
		stackElementStrSP := strings.Split(stackElementStr, "@")
		className := stackElementStrSP[0]
		classJsonPath := strings.ReplaceAll(className, ".", "/") + ".json"
		sid := stackElementStrSP[1]
		stmtLog.MStackTrace[classJsonPath+"@"+sid] = true
	}
	return nil
}

func (stmtLog *StmtLog) Debug() {
	for cls, methods := range stmtLog.MClasses {
		fmt.Println("==================================================")
		fmt.Println(cls)
		fmt.Println("==================================================")
		for _, method := range *methods {
			fmt.Println(method.MethodSig)
			for _, stmt := range method.Stmts {
				fmt.Println(stmt.MJid, stmt.MSid, stmt.MType, stmt.MValue)
			}
		}
	}
}

// 判断指定的源码行有没有在崩溃栈中出现过
func (stmtLog *StmtLog) Crashed(classJsonPath string, sid int) bool {
	if _, ok := stmtLog.MStackTrace[classJsonPath+"@"+strconv.Itoa(sid)]; ok {
		return true
	}
	return false
}

// 读取projectId/testcases/caseName下的stmtlog.json, 转化为StmtLog, 轻量级读取
func ReadStmtLogStd(projectId string, caseName string) (*StmtLog, *perrorx.ErrorX) {
	stmtLogPath := pdba.DBURLProjectIdTestcaseStmtlog(projectId, caseName)
	stmtLog, err := ReadStmtLog(stmtLogPath)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	return stmtLog, nil
}

// 读取path, 转化为StmtLog, 轻量级读取
func ReadStmtLog(path string) (*StmtLog, *perrorx.ErrorX) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, perrorx.NewErrorXReadFile(path, err.Error(), nil)
	}
	var stmtLog StmtLog
	err = json.Unmarshal(jsonData, &stmtLog)
	if err != nil {
		return nil, perrorx.NewErrorXUnmarshal(err.Error(), nil)
	}
	return &stmtLog, nil
}

// 读取projectId/testcases/caseName下的stmtlog.json, 解析成StmtLog
func ParseStmtLogStd(projectId string, caseName string) (*StmtLog, *perrorx.ErrorX) {
	stmtlogPath := pdba.DBURLProjectIdTestcaseStmtlog(projectId, caseName)
	logIdSigPath := pdba.DBURLProjectIdLogidsig(projectId)
	stmtLog, err := ParseStmtLog(stmtlogPath, logIdSigPath)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	return stmtLog, nil
}

// 读取path, 解析成StmtLog
func ParseStmtLog(path string, logIdSigPath string) (*StmtLog, *perrorx.ErrorX) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, perrorx.NewErrorXReadFile(path, err.Error(), nil)
	}
	var stmtLog StmtLog
	err = json.Unmarshal(jsonData, &stmtLog)
	if err != nil {
		return nil, perrorx.NewErrorXUnmarshal(err.Error(), nil)
	}
	err_ := stmtLog.Parse(logIdSigPath)
	if err_ != nil {
		return nil, perrorx.TransErrorX(err_)
	}
	return &stmtLog, nil
}