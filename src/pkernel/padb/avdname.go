package padb

import (
	"auiauto/pconfig"
	"auiauto/perrorx"
	"auiauto/putils"
	"strings"
)

// 根据avd id获取avd name, 比如avd id为emulator5554, 而avd name是我们自己设定的
func AdbEmuAvdName(avd string) (string, *perrorx.ErrorX) {
	ans, err := putils.Sh(pconfig.GConfig.MAdb + " -s " + avd + " emu avd name")
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	sp := strings.Split(ans, "\n")
	return sp[0] + ".avd", nil
}
