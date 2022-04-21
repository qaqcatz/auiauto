package psuscode

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"encoding/json"
	"io/ioutil"
	"path"
)

// 解析需要分析的文件(fix/root cause)
type AnalyzeClasses struct {
	MAnalyzeClasses []AnalyzeClass `json:"classes"`
}

// 获取聚焦语句总数
func (analyzeClasses *AnalyzeClasses)GetTotalLineNum() int {
	ans := 0
	for _, analyzeClass := range analyzeClasses.MAnalyzeClasses {
		ans += len(analyzeClass.MLines)
	}
	return ans
}

// analyzeFile是文件名, 需要做路径的拼接
func ReadAnalyzeFile(projectId string, analyzeFile string) (*AnalyzeClasses, *perrorx.ErrorX){
	analyzeFilePath := path.Join(pdba.DBURLProjectIdTestcases(projectId), analyzeFile+".json")
	data, err := ioutil.ReadFile(analyzeFilePath)
	if err != nil {
		return nil, perrorx.NewErrorXReadFile(analyzeFilePath, err.Error(), nil)
	}
	var analyzeClasses AnalyzeClasses
	err = json.Unmarshal(data, &analyzeClasses)
	if err != nil {
		return nil, perrorx.NewErrorXUnmarshal(err.Error(), nil)
	}
	return &analyzeClasses, nil
}
