package padb

import (
	"auiauto/pconfig"
	"auiauto/perrorx"
	"auiauto/putils"
	"net"
	"strings"
)

// 获取avd的ip, 若avd以emulator开头则视为模拟器, 返回127.0.0.1
func AdbWlanIp(avd string) (string, *perrorx.ErrorX) {
	if strings.HasPrefix(avd, "emulator") {
		return "127.0.0.1", nil
	}
	ans, err := AdbShell(avd, "ip addr show wlan0 | grep 'inet\\s' | awk '{print $2}' | awk -F '/' '{print $1}'")
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	// 格式化
	ans = strings.TrimSpace(ans)
	// 检查地址格式是否合法
	address := net.ParseIP(ans)
	if address == nil {
		return "", perrorx.NewErrorXInvalidAddress(ans, nil)
	}
	return ans, nil
}

// 执行adb -s avd forward tcp:pcPort tcp:androidPort做端口映射
func AdbForward(avd string, pcPort string, androidPort string) (string, *perrorx.ErrorX) {
	ans, err := putils.Sh(pconfig.GConfig.MAdb + " -s " + avd + " forward tcp:" + pcPort + " tcp:" + androidPort)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}

// 执行adb -s avd forward --remove-all
func AdbForwardRemoveAll(avd string) (string, *perrorx.ErrorX) {
	ans, err := putils.Sh(pconfig.GConfig.MAdb + " -s " + avd + " forward --remove-all")
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}
