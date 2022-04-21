package padb

import (
	"auiauto/pconfig"
	"auiauto/perrorx"
	"auiauto/putils"
)

// 执行adb -s avd push src(android) dst(pc)
func AdbPush(avd string, src string, dst string) (string, *perrorx.ErrorX) {
	ans, err := putils.Sh(pconfig.GConfig.MAdb +" -s " + avd + " push " + src + " " + dst)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}

// 执行adb -s avd pull src(android) dst(pc)
func AdbPull(avd string, src string, dst string) (string, *perrorx.ErrorX) {
	ans, err := putils.Sh(pconfig.GConfig.MAdb +" -s " + avd + " pull " + src + " " + dst)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}
