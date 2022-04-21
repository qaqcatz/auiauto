package pahttp

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"io/ioutil"
)

// 获取日志
// -1: antrance访问不到antrance ins
// -2: 未知错误, 一般是用户没有按照规定的流程操作造成的
func getStmtLog(avd string) (string, *perrorx.ErrorX) {
	ans, err := AntranceRequest("GET", avd, "stmtlog", nil)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	if ans == "-1" {
		return "", perrorx.NewErrorXGetStmtLog("antrance can not access antrance ins", nil)
	}
	if ans == "-2" {
		return "", perrorx.NewErrorXGetStmtLog("unknown error", nil)
	}
	return ans, nil
}

// 获取语句日志, 保存在DBURLProjectIdTestcaseStmtlog下
func GetStmtLogAndSaveStd(avd string, projectId string, caseName string) *perrorx.ErrorX {
	stmtLogPath := pdba.DBURLProjectIdTestcaseStmtlog(projectId, caseName)
	err := GetStmtLogAndSave(avd, stmtLogPath)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}

// 获取语句日志, 保存在path下
func GetStmtLogAndSave(avd string, path string) *perrorx.ErrorX {
	ans, err := getStmtLog(avd)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	err_ := ioutil.WriteFile(path, []byte(ans), 0777)
	if err_ != nil {
		return perrorx.NewErrorXWriteFile(path, err_.Error(), nil)
	}
	return nil
}
