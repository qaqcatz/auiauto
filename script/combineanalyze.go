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
	casePrefix := "art"
	analyzeFile := "fix"
	factor := "Ochiai"
	for i := 0; i < len(Projects); i++ {
		fmt.Println("combine analyze project " + strconv.Itoa(i) + " " + Projects[i])
		resp, err := http.Get("http://" + ip + ":" + port + "/combineanalyze?" +
			"&projectId=" + Projects[i] +
			"&casePrefix=" + casePrefix +
			"&analyzeFile=" + analyzeFile +
			"&factor=" + factor)
		if err != nil {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println(err.Error())
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
		if resp != nil {
			data, err_ := ioutil.ReadAll(resp.Body)
			if err_ != nil {
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				fmt.Println(err_.Error())
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			}
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			fmt.Println(strconv.Itoa(resp.StatusCode) + ": " + string(data))
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			_ = resp.Body.Close()
		}
	}
}
