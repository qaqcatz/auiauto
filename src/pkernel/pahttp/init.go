package pahttp

import (
	"auiauto/perrorx"
)

// 应用生成日志后会禁止后续ui tree的获取, 需要init
func GetInit(avd string) *perrorx.ErrorX {
	_, err := AntranceRequest("GET", avd, "init", nil)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}
