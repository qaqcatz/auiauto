package padb

import (
	"auiauto/pconfig"
	"auiauto/perrorx"
	"auiauto/putils"
	"strings"
)

// 判断当前是否有设备offline
func adbOffline() bool {
	str, _ := putils.Sh(pconfig.GConfig.MAdb + " devices")
	// 对adb执行结果做一些格式处理
	devices := strings.Split(str, "\n")
	for i := 1; i < len(devices); i++ {
		if devices[i] == "" {
			continue
		}
		sp := strings.Split(devices[i], "\t")
		if len(sp) != 2 {
			continue
		}
		if sp[1] == "offline" {
			return true
		}
	}
	return false
}

// adb reconnect offline
func AdbReconnectOffline() (string, *perrorx.ErrorX) {
	if adbOffline() {
		ans, err := putils.Sh(pconfig.GConfig.MAdb + " reconnect offline")
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		return ans, nil
	}
	return "", nil
}