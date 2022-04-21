package pahttp

import (
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"bytes"
	"io/ioutil"
	"net/http"
)

// 通用的http request方法, 自动开启端口转发
func AntranceRequest(method string, avd string, paramUrl string, jsonData []byte) (string, *perrorx.ErrorX) {
	_, err := padb.AdbReconnectOffline()
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	// 获取ip
	avdIp, err := padb.AdbWlanIp(avd)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	// 判断是否需要forward, 只有连接模拟器时才需要forward
	port := "8624"
	if avdIp == "127.0.0.1" {
		port, err = connectAvd(avd)
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
	}
	var resp *http.Response = nil
	var err_ error
	if method == "GET" { // GET
		resp, err_ = http.Get("http://" + avdIp + ":" + port + "/" + paramUrl)
	} else { // POST
		resp, err_ = http.Post("http://"+avdIp+":"+port+"/"+paramUrl, "application/json", bytes.NewBuffer(jsonData))
	}
	if resp == nil {
		return "", perrorx.NewErrorXAntranceRquest(method+" http://"+avdIp+":"+port+"/"+paramUrl,
			"resp == nil", nil)
	}
	if err_ != nil {
		return "", perrorx.NewErrorXAntranceRquest(method+" http://"+avdIp+":"+port+"/"+paramUrl,
			method+" error: "+err_.Error(), nil)
	}
	data, err_ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err_ != nil {
		return "", perrorx.NewErrorXReadAll(err_.Error(), nil)
	}
	if resp.StatusCode != http.StatusOK {
		return "", perrorx.NewErrorXAntranceRquest(method+" http://"+avdIp+":"+port+"/"+paramUrl,
			"http.StatusOK", nil)
	}
	return string(data), nil
}
