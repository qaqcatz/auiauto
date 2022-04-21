package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	emulatorCmd := "/home/android/Android/Sdk/emulator/emulator"

	var wg sync.WaitGroup
	//P := "Images-to-PDF_771"
	//P := "FirefoxLite-4942"
	P := "Simple-Music-Player_128"
	for i, project := range Projects {
		if project != P {
			continue
		}
		name := "auiauto"+strconv.Itoa(i+1)
		port := 5554+(i*2)
		wg.Add(1)
		go func (avdName string, portStr string) {
			snapshot := "-snapshot init"
			//snapshot := ""
			//noWindow := "-no-window"
			noWindow := ""
			_, err := Sh("nohup " + emulatorCmd + " " +
				"-avd "+avdName+" "+
				"-port "+portStr+" "+snapshot+" "+noWindow)
			if err != nil {
				fmt.Println(err.Error())
			}
			wg.Done()
		} (name, strconv.Itoa(port))
		fmt.Println("avd " + name + " running...")
		break
	}
	wg.Wait()
	fmt.Println("bye~")
}