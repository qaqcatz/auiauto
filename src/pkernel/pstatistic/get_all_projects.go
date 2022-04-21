package pstatistic

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/putils"
	"io/ioutil"
	"path"
)

// 获取projects下的所有项目, 通过是否包含testcases目录来判断
func getAllProjects() ([]string, *perrorx.ErrorX) {
	ans := make([]string, 0)
	projectsPath := pdba.DBURLProjects()
	dir, err := ioutil.ReadDir(projectsPath)
	if err != nil {
		return nil, perrorx.NewErrorXReadDir(projectsPath, err.Error(), nil)
	}
	for _, fi := range dir {
		if fi.IsDir() {
			if putils.FileExist(path.Join(projectsPath, fi.Name(), "testcases")) {
				ans = append(ans, fi.Name())
			}
		}
	}
	return ans, nil
}
