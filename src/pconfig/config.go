package pconfig

import (
	"auiauto/perrorx"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// 启动配置
type Config struct {
	// adb可执行文件所在的路径, 可以用whereis adb查看
	MAdb string `json:"adb"`
	// avd所在的文件夹, 一般在~/.android/下
	MAvd string `json:"avd"`
	// ./database
	MDatabase string `json:"database"`
	// 后端服务的ip
	MIp string `json:"ip"`
	// 后端服务的port
	MPort string `json:"port"`
}

// 全局配置
var GConfig Config

// 从config.json下加载配置
func InitConfig() {
	configPath := "config.json"
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(perrorx.NewErrorXReadFile(configPath, err.Error(), nil))
		os.Exit(1)
	}
	err = json.Unmarshal(data, &GConfig)
	if err != nil {
		fmt.Println(perrorx.NewErrorXUnmarshal(err.Error(), nil))
		os.Exit(1)
	}
}
