package pevent

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/putils"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// 动作序列
type Events struct {
	// 初始快照
	MInitSnapshot string `json:"initSnapshot"`
	// 初始用例
	MInitTestcase string `json:"initTestcase"`
	// 是否重装
	MReInstall bool `json:"reInstall"`
	// 是否重启app
	MReStart bool `json:"reStart"`
	// 动作序列
	MEvents []*Event `json:"events"`
}

// 深拷贝
func (events *Events) Copy() *Events {
	newEvents := Events{events.MInitSnapshot,
		events.MInitTestcase,
		events.MReInstall,
		events.MReStart,
		make([]*Event, len(events.MEvents))}
	for i := 0; i < len(events.MEvents); i++ {
		event := events.MEvents[i]
		newPrefix := make([]int, len(event.MPrefix))
		for j := 0; j < len(event.MPrefix); j++ {
			newPrefix[j] = event.MPrefix[j]
		}
		newEvents.MEvents[i] = &Event{event.MId, event.MType, event.MValue,
			event.MObject, newPrefix, event.MDesc}
	}
	return &newEvents
}

// 从projectId/testcases/{caseName}下读取testcase.json, 存储为Events
func ReadEventsStd(projectId string, caseName string) (*Events, *perrorx.ErrorX) {
	jsonPath := pdba.DBURLProjectIdTestcaseTestcase(projectId, caseName)
	events, err := ReadEvents(jsonPath)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	return events, nil
}

// 从path下读取Events
func ReadEvents(path string) (*Events, *perrorx.ErrorX) {
	// 读取path, 路径不存在报错
	if !putils.FileExist(path) {
		return nil, perrorx.NewErrorXFileNotFound(path, nil)
	}
	// 获取json data
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, perrorx.NewErrorXReadFile(path, err.Error(), nil)
	}
	// 绑定到Events上
	var events Events
	err = json.Unmarshal(jsonData, &events)
	if err != nil {
		return nil, perrorx.NewErrorXUnmarshal(err.Error(), nil)
	}
	return &events, nil
}

// 将events写入projectId/testcases/caseName/testcase.json, 路径不存在时会自动创建
func WriteEventsStd(projectId string, caseName string, events *Events) *perrorx.ErrorX {
	projectPath := pdba.DBURLProjectId(projectId)
	casePath := pdba.DBURLProjectIdTestcase(projectId, caseName)
	// 没有项目新建项目
	if !putils.FileExist(projectPath) {
		_ = os.Mkdir(projectPath, 0777)
		_ = os.Mkdir(path.Join(projectPath, "testcases"), 0777)
	}
	// caseName存在拒绝覆盖, 不存在新建
	if putils.FileExist(casePath) {
		return perrorx.NewErrorXFileExist(casePath, nil)
	} else {
		_ = os.MkdirAll(casePath, 0777)
	}
	err := WriteEvents(path.Join(casePath, "testcase.json"), events)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}

// 将Events写到path下
func WriteEvents(path string, events *Events) *perrorx.ErrorX {
	jsonData, err := json.Marshal(events)
	if err != nil {
		return perrorx.NewErrorXMarshal(err.Error(), nil)
	}
	err = ioutil.WriteFile(path, jsonData, 0777)
	if err != nil {
		return perrorx.NewErrorXWriteFile(path, err.Error(), nil)
	}
	return nil
}
