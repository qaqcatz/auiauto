package padb

import (
	"auiauto/perrorx"
)

// 设置/取消自动旋转
func AdbAccRotation(avd string, i string) *perrorx.ErrorX {
	_, err := AdbShell(avd, "content insert --uri content://settings/system --bind name:s:accelerometer_rotation --bind value:i:"+i)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}

// 旋转, 0竖1横
func AdbUserRotation(avd string, i string) *perrorx.ErrorX {
	_, err := AdbShell(avd, "content insert --uri content://settings/system --bind name:s:user_rotation --bind value:i:"+i)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}
