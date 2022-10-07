package pahttp

import (
	"auiauto/perrorx"
)

// 获取ui tree
func GetUITree(avd string) (string, *perrorx.ErrorX) {
	ans, err := AntranceRequest("GET", avd, "uitree", nil)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}
