package padb

import (
	"auiauto/perrorx"
	"strconv"
	"time"
)

// 模拟点击 adb -s avd shell input tap x y
func AdbInputTap(avd string, x int, y int) *perrorx.ErrorX {
	_, err := AdbShell(avd, "input tap "+strconv.Itoa(x)+" "+strconv.Itoa(y))
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}

// 模拟滑动, adb -s avd shell input swipe x1 y1 x2 y2 time(ms), 默认等待滑动时间
// 若两个坐标相同, 且ms < 100则转为tap操作
func AdbInputSwipe(avd string, x1 int, y1 int, x2 int, y2 int, ms int) *perrorx.ErrorX {
	if x1 == x2 && y1 == y2 && ms < 100 {
		err := AdbInputTap(avd, x1, y1)
		if err != nil {
			return perrorx.TransErrorX(err)
		}
		return nil
	}
	_, err := AdbShell(avd, "input swipe "+strconv.Itoa(x1)+" "+strconv.Itoa(y1)+" "+
		strconv.Itoa(x2)+" "+strconv.Itoa(y2)+" "+strconv.Itoa(ms))
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	time.Sleep(time.Millisecond * time.Duration(ms))
	return nil
}

// 设置全局事件, adb -s avd shell input keyevent KEYCODE_BACK
// https://developer.android.com/reference/android/view/KeyEvent
func AdbInputKeyEvent(avd string, event string) *perrorx.ErrorX {
	_, err := AdbShell(avd, "input keyevent "+event)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}
