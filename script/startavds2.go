package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	emulatorCmd := "/home/android/Android/Sdk/emulator/emulator"

	var wg sync.WaitGroup
	for i := 12; i < 24; i++ {
		name := "auiauto"+strconv.Itoa(i+1)
		port := 5554+(i*2)

		wg.Add(1)
		go func (avdName string, portStr string) {
			snapshot := "-snapshot init"
			//snapshot := ""
			noWindow := "-no-window"
			//noWindow := ""
			_, err := Sh("nohup " + emulatorCmd + " " +
				"-avd "+avdName+" "+
				"-port "+portStr+" "+snapshot+" "+noWindow)
			if err != nil {
				fmt.Println(err.Error())
			}
			wg.Done()
		} (name, strconv.Itoa(port))
		fmt.Println("avd " + name + " running...")
		time.Sleep(time.Millisecond*3000) // avoid shock
	}
	wg.Wait()
	fmt.Println("bye~")
}