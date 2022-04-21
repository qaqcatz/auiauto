package panalyze

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/psuscode"
	"auiauto/putils"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

// 组合分析, 在指定目录下生成分析文件combineanalyze_casePrefix_analyzeFile_factor.txt
// 格式:
// 第一行两个部分, 表示用例总数n(组合总数2^(n-1)), 聚焦语句数量s
// 第二行空格打印所用的用例
// 接下来2^(n-1)个组合, 每个组合:
// 	第一行一个01字符串(长度为n,原始错误用例一定为1), 二进制形式表示所选的用例, le存储
// 	  接下来s个数, 代表聚焦语句文件中每个聚焦语句的排名

// 2^n组合分析, 在projectId/testcases/下生成combineanalyze_casePrefix_analyzeFile_factor.txt分析文件
// 切片中的语句顺序和analyzeFile的顺序是一致的, 有需要的话可以自己按顺序还原到各个类的各个语句上
// 更新: 组合分析需要默认包含原始错误用例, 规定以crash为后缀的为原始错误用例
func ReadSourceTreeAndCombineAnalyzeStd(projectId string, casePrefix string, analyzeFile string, factor string) *perrorx.ErrorX {
	testcasesPath := pdba.DBURLProjectIdTestcases(projectId)
	caseNames, err := putils.GetDirsStartWith(casePrefix, testcasesPath)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	if len(caseNames) > 10 {
		return perrorx.NewErrorXReadSourceTreeAndCombineAnalyze("> 2^10!", nil)
	}
	// 读取analyzeFile, 获取聚焦语句总数
	analyzeClasses, err := psuscode.ReadAnalyzeFile(projectId, analyzeFile)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	totalLineNum := analyzeClasses.GetTotalLineNum()
	// 获取原始错误用例在caseNames中的下标
	originCrashIdx := 0
	for i, caseName := range caseNames {
		if strings.HasSuffix(caseName, "crash") {
			originCrashIdx = i
		}
	}
	defer func() {
		fmt.Println()
	} ()
	ans := ""
	// 第二行三个部分, 表示用例总数n(组合总数2^n-1, 全不选的情况不考虑), 聚焦语句数量s
	ans += strconv.Itoa(len(caseNames)) + " " + strconv.Itoa(totalLineNum) + "\n"
	// 第二行空格打印所用的用例, 原始错误用例标*
	for i := 0; i < len(caseNames); i++ {
		if i == originCrashIdx {
			ans += "*" + caseNames[i] + " "
		} else {
			ans += caseNames[i] + " "
		}
	}
	ans += " \n"
	// 接下来2^(n-1)个组合, 除去不包含originCrashIdx的情况
	total := 1<<(len(caseNames))
	cur := 1
	for i := 0; i < total; i++ {
		cur++
		fmt.Print("\r" + strconv.Itoa(cur) + "/" + strconv.Itoa(total))
		if i&(1<<originCrashIdx) == 0 {
			continue
		}
		// 第一行一个01字符串, 二进制形式表示所选的用例
		s01 := ""
		each := make([]string, 0)
		for j := 0; j < len(caseNames); j++ {
			if i&(1<<j) != 0 {
				s01 += "1"
				each = append(each, caseNames[j])
			} else {
				s01 += "0"
			}
		}
		ans += s01 + "\n"
		// 调用readSourceTreeAndAnalyzeStd进行分析
		_, susCodeSlice, err := readSourceTreeAndAnalyzeStd(projectId, each, analyzeFile, factor)
		if err != nil {
			return perrorx.TransErrorX(err)
		}
		// 接下来s个数, 代表聚焦语句文件中每个聚焦语句的排名
		for _, susCode := range susCodeSlice {
			ans += strconv.Itoa(susCode.MRank) + " "
		}
		ans += "\n"
	}
	// 最后在projectId/testcases/下生成combineanalyze_casePrefix_analyzeFile_factor.txt分析文件
	combineAnalyzePath := path.Join(pdba.DBURLProjectIdTestcases(projectId),
		"combineanalyze_"+casePrefix+"_"+analyzeFile+"_"+factor+".txt")
	err_ := ioutil.WriteFile(combineAnalyzePath, []byte(ans), 0777)
	if err_ != nil {
		return perrorx.NewErrorXWriteFile(combineAnalyzePath, err_.Error(), nil)
	}
	return nil
}

