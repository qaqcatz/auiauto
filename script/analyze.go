package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	ip := "127.0.0.1"
	port := "8082"
	casePrefix := "art"
	analyzeFile := "fix"
	factor := "Op2"
	for i := 0; i < len(Projects); i++ {
		fmt.Println("analyze project " + strconv.Itoa(i) + " " + Projects[i])
		resp, err := http.Get("http://" + ip + ":" + port + "/analyze?" +
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
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			fmt.Println(strconv.Itoa(resp.StatusCode))
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			_ = resp.Body.Close()
		}
	}
}

