package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	adbCmd := "/home/android/Android/Sdk/platform-tools/adb"
	for i := 36; i < 46; i++ {
		port := 5554+(i*2)
		avdId := "emulator-"+strconv.Itoa(port)
		ans, err := Sh(adbCmd + " -s " + avdId + " emu kill")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(ans)
		}
		time.Sleep(time.Millisecond*500) // avoid shock
	}
}