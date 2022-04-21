package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	ip := "127.0.0.1"
	port := "8082"
	analyzeFile := "rootcause"
	factor := "Ochiai"
	tester := "monkey"
	casePrefix := "rd"
	randNum := 100

	curFinished := 0
	for i := 0; i < len(Projects); i++ {
		fmt.Println("[project] " + strconv.Itoa(i) + " " + Projects[i])
		resp, err := http.Get("http://" + ip + ":" + port + "/rdanalyze?" +
			"&projectId=" + Projects[i] +
			"&analyzeFile=" + analyzeFile +
			"&factor=" + factor +
			"&tester=" + tester +
			"&casePrefix=" + casePrefix +
			"&randNum=" + strconv.Itoa(randNum))
		if err != nil {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println(Projects[i] + " " + err.Error())
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
		if resp != nil {
			data, err_ := ioutil.ReadAll(resp.Body)
			if err_ != nil {
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				fmt.Println(Projects[i] + " " +  err_.Error())
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			}
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			fmt.Println(Projects[i] + " " + strconv.Itoa(resp.StatusCode) + ": " + string(data))
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			_ = resp.Body.Close()
		}
		curFinished++
		fmt.Println(strconv.Itoa(curFinished) + "/" + strconv.Itoa(len(Projects)))
	}

	fmt.Println("fixiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiig")
	analyzeFile = "fix"
	curFinished = 0
	for i := 0; i < len(Projects2); i++ {
		fmt.Println("[project] " + strconv.Itoa(i) + " " + Projects2[i])
		resp, err := http.Get("http://" + ip + ":" + port + "/rdanalyze?" +
			"&projectId=" + Projects2[i] +
			"&analyzeFile=" + analyzeFile +
			"&factor=" + factor +
			"&tester=" + tester +
			"&casePrefix=" + casePrefix +
			"&randNum=" + strconv.Itoa(randNum))
		if err != nil {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println(Projects2[i] + " " + err.Error())
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
		if resp != nil {
			data, err_ := ioutil.ReadAll(resp.Body)
			if err_ != nil {
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				fmt.Println(Projects2[i] + " " +  err_.Error())
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			}
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			fmt.Println(Projects2[i] + " " + strconv.Itoa(resp.StatusCode) + ": " + string(data))
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			_ = resp.Body.Close()
		}
		curFinished++
		fmt.Println(strconv.Itoa(curFinished) + "/" + strconv.Itoa(len(Projects2)))
	}
}
