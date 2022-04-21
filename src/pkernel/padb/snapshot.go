package padb

import (
	"auiauto/pconfig"
	"auiauto/perrorx"
	"auiauto/putils"
)

// 加载快照
func AdbEmuLoadSnapshot(avd string, name string) (string, *perrorx.ErrorX) {
	_, err := AdbReconnectOffline()
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	ans, err := putils.Sh(pconfig.GConfig.MAdb + " -s " + avd + " emu avd snapshot load " + name)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	_, err = AdbReconnectOffline()
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}

// 保存快照
func AdbEmuSaveSnapshot(avd string, name string) (string, *perrorx.ErrorX) {
	ans, err := putils.Sh(pconfig.GConfig.MAdb + " -s " + avd + " emu avd snapshot save " + name)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}

// 删除快照
func AdbEmuDeleteSnapshot(avd string, name string) (string, *perrorx.ErrorX) {
	ans, err := putils.Sh(pconfig.GConfig.MAdb + " -s " + avd + " emu avd snapshot delete " + name)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}
