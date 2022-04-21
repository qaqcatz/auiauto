package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 13+1 5580 -> 3138
// 25+1 5604 -> 4942
// 42+1 5638 -> 128
// 4+1  5562 -> 261
// 5+1  5564 -> 375
// 6+1  5566 -> 480
// 41+1 5636 -> 32
func main() {
	ip := "127.0.0.1"
	port := "8082"

	var wg sync.WaitGroup
	var mutex sync.Mutex
	curFinished := 0
	for i := 0; i < len(Projects); i++ {
		wg.Add(1)
		go func (projectId string, avdId string) {
			// 每次测试前加载一下初始快照
			_, _ = http.Get("http://" + ip + ":" + port + "/loadsnapshot?avd="+avdId+"&name=init")
			time.Sleep(time.Millisecond*3000)
			fmt.Println(projectId + " start in " + avdId)
			resp, err := http.Get("http://" + ip + ":" + port + "/testing?avd=" + avdId +
				"&projectId=" + projectId +
				"&crashCase=origin_crash" +
				"&tester=monkey" +
				"&testNum=-3600000" +
				"&testPrefix=rd"+
				"&testParam=none")
			if err != nil {
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				fmt.Println(projectId + " " + err.Error())
				fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			}
			if resp != nil {
				data, err_ := ioutil.ReadAll(resp.Body)
				if err_ != nil {
					fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
					fmt.Println(projectId + " " + err_.Error())
					fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				}
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
				fmt.Println(projectId + " " + strconv.Itoa(resp.StatusCode) + ": " + string(data))
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
				_ = resp.Body.Close()
			}
			wg.Done()
			mutex.Lock()
			curFinished++
			fmt.Println(projectId + " " + strconv.Itoa(curFinished) + "/" + strconv.Itoa(len(Projects)))
			mutex.Unlock()
		} (Projects[i], "emulator-"+strconv.Itoa(5554+i*2))
		time.Sleep(time.Millisecond*1000) // avoid shock
	}
	wg.Wait()
}
