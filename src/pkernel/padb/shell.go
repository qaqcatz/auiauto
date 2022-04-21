package padb

import (
	"auiauto/pconfig"
	"auiauto/perrorx"
	"auiauto/putils"
)

// 调用adb shell -s avd执行adbCmd, 禁止直接使用adb shell(会造成阻塞)
func AdbShell(avd string, adbCmd string) (string, *perrorx.ErrorX) {
	_, err_ := AdbReconnectOffline()
	if err_ != nil {
		return "", perrorx.TransErrorX(err_)
	}
	// 防止直接使用shell命令
	if adbCmd == "shell" {
		return "", perrorx.NewErrorXADBShellBlock(nil)
	}
	ans, err := putils.Sh(pconfig.GConfig.MAdb + " -s " + avd + " shell " + adbCmd)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}
