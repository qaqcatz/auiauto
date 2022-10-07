package pahttp

import (
	"auiauto/perrorx"
)

// 务必在每次测试前调用一下init
func GetInit(avd string) (string, *perrorx.ErrorX) {
	ans, err := AntranceRequest("GET", avd, "init", nil)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}
