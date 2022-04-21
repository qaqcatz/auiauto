package padb

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/putils"
	"encoding/base64"
	"io/ioutil"
	"os"
)

// 获取屏幕截图, 执行adb -s avd shell screencap -p /storage/emulated/0/Download/screenshot.png.
// 先删除screenshot.png, 然后将图片存储在database下, 名字为screenshot.png,
// 最后读取图片, 转换成base64编码返回
func GetScreenShot(avd string) (string, *perrorx.ErrorX) {
	screenShotPath := pdba.DBURLTmpScreenshot()
	// 1. 删除原本的screenshot.png
	if putils.FileExist(screenShotPath) {
		_ = os.Remove(screenShotPath)
	}
	// 2. 拉取screenshot.png
	_, err := AdbShell(avd, "screencap -p /sdcard/Download/screenshot.png")
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	_, err = AdbPull(avd, "/sdcard/Download/screenshot.png", screenShotPath)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	// 3. 读取screenshot.png, 转换成base64编码返回
	data, err_ := ioutil.ReadFile(screenShotPath)
	if err_ != nil {
		return "", perrorx.NewErrorXReadFile(screenShotPath, err_.Error(), nil)
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	encoded = "data:image/png;base64,"+encoded
	return encoded, nil
}