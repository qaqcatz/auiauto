package pahttp

import (
	"auiauto/perrorx"
)

// 获取ui tree
// -1: 当前应用崩溃或已经调用过getStmtLog完成了一个流程
func GetUITree(avd string) (string, *perrorx.ErrorX) {
	ans, err := AntranceRequest("GET", avd, "uitree", nil)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	if ans == "-1" {
		return "", perrorx.NewErrorXGetUITree("a stmtlog here, you can not get ui tree", nil)
	}
	return ans, nil
}
