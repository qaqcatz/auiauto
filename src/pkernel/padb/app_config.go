package padb

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"io/ioutil"
	"strings"
)

// 读取projectId/apk/config.txt, 获取applicationId和mainActivity
func GetAppConfig(projectId string) (string, string, *perrorx.ErrorX) {
	appConfigPath := pdba.DBURLProjectIdAPKConfig(projectId)
	data, err := ioutil.ReadFile(appConfigPath)
	if err != nil {
		return "", "", perrorx.NewErrorXReadFile(appConfigPath, err.Error(), nil)
	}
	sp := strings.Split(string(data), "@")

	if len(sp) != 2 {
		return "", "", perrorx.NewErrorXGetAppConfig(projectId+"/apk/config.txt", nil)
	}
	return sp[0], sp[1], nil
}
