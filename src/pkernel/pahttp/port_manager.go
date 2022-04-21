package pahttp

import (
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"strconv"
	"strings"
)

// 传入avdId, 通过avdId获取avd端口, 比如emulator-5554的端口为5554, 如果不是这个格式则报错
// 接着把端口号+1000作为forward的端口, 开启端口转发
// 返回端口号
func connectAvd(avd string) (string, *perrorx.ErrorX) {
	if strings.Contains(avd, "-") {
		sp := strings.SplitN(avd, "-", 2)
		port, err := strconv.Atoi(sp[1])
		if err != nil {
			return "", perrorx.NewErrorXParseInt(sp[1], nil)
		}
		newPort := port + 1000
		newPortStr := strconv.Itoa(newPort)
		_, _ = padb.AdbForward(avd, newPortStr, "8624")
		return newPortStr, nil
	} else {
		return "", perrorx.NewErrorXConnectAvd(avd + " format error", nil)
	}
}
