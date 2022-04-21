package padb

import (
	"auiauto/pconfig"
	"auiauto/perrorx"
	"auiauto/putils"
	"strings"
)

// 执行adb devices, 返回当前可用的avd列表, 存储在切片中(已经进行了格式处理)
func AdbDevices() ([]string, *perrorx.ErrorX) {
	var ans []string
	str, err := putils.Sh(pconfig.GConfig.MAdb + " devices")
	if err != nil {
		return make([]string, 0), perrorx.TransErrorX(err)
	}
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
		avd := sp[0]
		ans = append(ans, avd)
	}
	return ans, nil
}