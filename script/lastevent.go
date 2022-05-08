package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	ip := "127.0.0.1"
	port := "8082"
	analyzeFile := "fix"
	for i := 0; i < len(Projects); i++ {
		resp, err := http.Get("http://" + ip + ":" + port + "/lastevent?" +
			"&projectId=" + Projects[i] +
			"&analyzeFile="+analyzeFile)
		if err != nil {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println(Projects[i] + " " + err.Error())
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			return
		}
		if resp != nil {
			data, err_ := ioutil.ReadAll(resp.Body)
			if err_ != nil {
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				fmt.Println(Projects[i] + " " +  err_.Error())
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				return
			}
			fmt.Print(Projects[i])
			for j := 0; j < 80-len(Projects[i]); j++ {
				fmt.Print(" ")
			}
			fmt.Println(string(data))
			_ = resp.Body.Close()
		}
	}
}
















































