package panalyze

import (
	"auiauto/perrorx"
	"io/ioutil"
	"path"
	"strings"
)

// 获取dir下的所有以casePrefix为前缀的pass用例和crash(返回路径), 用名字是否包含pass|crash来判断
// ignore _X
func readPassCrash(dirPath string, casePrefix string) ([]string, []string, *perrorx.ErrorX) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, nil, perrorx.NewErrorXReadDir(dirPath, err.Error(), nil)
	}
	passPaths := make([]string, 0)
	crashPaths := make([]string, 0)
	for _, fi := range dir {
		if fi.IsDir() && strings.HasPrefix(fi.Name(), casePrefix) {
			if strings.HasSuffix(fi.Name(), "_X") {
				continue
			}
			if strings.Contains(fi.Name(), "pass") {
				passPaths = append(passPaths, path.Join(dirPath, fi.Name()))
			} else if strings.Contains(fi.Name(), "crash") {
				crashPaths = append(crashPaths, path.Join(dirPath, fi.Name()))
			}
		}
	}
	return passPaths, crashPaths, nil
}

// 获取dir下的所有以casePrefix为前缀的pass用例和crash_0(返回路径), 用名字是否包含pass|crash来判断
// ignore _X
func readPassCrash0(dirPath string, casePrefix string) ([]string, []string, *perrorx.ErrorX) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, nil, perrorx.NewErrorXReadDir(dirPath, err.Error(), nil)
	}
	passPaths := make([]string, 0)
	crashPaths := make([]string, 0)
	for _, fi := range dir {
		if fi.IsDir() && strings.HasPrefix(fi.Name(), casePrefix) {
			if strings.HasSuffix(fi.Name(), "_X") {
				continue
			}
			if strings.Contains(fi.Name(), "pass") {
				passPaths = append(passPaths, path.Join(dirPath, fi.Name()))
			} else if strings.HasSuffix(fi.Name(), "crash_0") {
				crashPaths = append(crashPaths, path.Join(dirPath, fi.Name()))
			}
		}
	}
	return passPaths, crashPaths, nil
}

