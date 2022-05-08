package psuscode

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/putils"
	"path"
	"strconv"
	"strings"
)

// 可疑语句切片
type SusCodes []SusCode

func (s SusCodes) Len() int           { return len(s) }
func (s SusCodes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SusCodes) Less(i, j int) bool { return s[i].MValue > s[j].MValue }

// 读取projectId/testcases/analyzeFile, 根据susCodes计算全局可疑语句切片GSusCodesSlice
// 注意切片中的语句顺序和analyzeFile中的语句顺序是一致的
func CreateSlice(projectId string, analyzeFile string, susCodes SusCodes) (SusCodes, *perrorx.ErrorX) {
	var susCodesSlice SusCodes = nil

	analyzeFilePath := path.Join(pdba.DBURLProjectIdTestcases(projectId), analyzeFile+".json")
	if !putils.FileExist(analyzeFilePath) {
		// 待分析文件找不到的话默认分析top100
		susCodesSlice = make(SusCodes, 0)
		cnt := 0
		for i := 0; i < len(susCodes); i++ {
			// 这个判断是为了方便统计纯栈信息的效果的, 有时也会用来过滤手动异常捕获类MyCrash, 平时不会触发
			if susCodes[i].MValue < 0 {
				continue
			}
			susCodesSlice = append(susCodesSlice, susCodes[i])
			cnt++
			if cnt == 100 {
				break
			}
		}
	} else {
		// 若待分析文件存在, 则在susCodes中寻找对应语句的排名, 生成GSusCodesSlice
		analyzeClasses, err := ReadAnalyzeFile(projectId, analyzeFile)
		if err != nil {
			return nil, perrorx.TransErrorX(err)
		}
		// class:line->id, id为当前行在analyzeFile中的下标(将class和line展开后的整体下标)
		analyzeLineToIdMap := make(map[string]int)
		// 统计class和line展开后的id总数
		sumIdx := 0
		// 初始化, 默认rank为0, 表示未被覆盖
		susCodesSlice = make(SusCodes, 0)
		for _, analyzeClass := range analyzeClasses.MAnalyzeClasses {
			className := analyzeClass.MClassName
			for _, line := range analyzeClass.MLines {
				analyzeLineToIdMap[className+":"+strconv.Itoa(line)] = sumIdx
				sumIdx++
				shortClassName := className
				if strings.Contains(shortClassName, ".") {
					shortClassName = shortClassName[strings.LastIndex(shortClassName, ".")+1:]
				}
				susCodesSlice = append(susCodesSlice, SusCode{nil, 0, shortClassName,
					className, line, 0, 0, 0, 0, 0, 0})
			}
		}
		// 遍历suscode, 若可疑语句在analyzeLineToIdMap中存在, 则根据id将susCode设置到对应的susCodesSlice项
		for i := 0; i < susCodes.Len(); i++ {
			// 这个判断是为了方便统计纯栈信息的效果的, 有时也会用来过滤手动异常捕获类MyCrash, 平时不会触发
			if susCodes[i].MValue < 0 {
				continue
			}
			classLine := susCodes[i].MClassName + ":" + strconv.Itoa(susCodes[i].MLine)
			if _, exist := analyzeLineToIdMap[classLine]; exist {
				susCodesSlice[analyzeLineToIdMap[classLine]] = susCodes[i]
			}
		}
	} // end of else(!utils.FileExist(analyzeFilePath))
	return susCodesSlice, nil
}

// 初始化susCodes
func DoInit(susCodes SusCodes, totalPassed int, totalFailed int) {
	for i := 0; i < len(susCodes); i++ {
		susCodes[i].MClassShortName = susCodes[i].MOriginNode.MName
		susCodes[i].MClassName = susCodes[i].MOriginNode.MFullName
		susCodes[i].MLine = susCodes[i].MIdx + 1
		susCodes[i].MA11 = susCodes[i].MOriginNode.MFailed[susCodes[i].MIdx]
		susCodes[i].MA10 = susCodes[i].MOriginNode.MPassed[susCodes[i].MIdx]
		susCodes[i].MA01 = totalFailed-susCodes[i].MA11
		susCodes[i].MA00 = totalPassed-susCodes[i].MA10
	}
}