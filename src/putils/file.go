package putils

import (
	"auiauto/perrorx"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 判断文件是否存在
func FileExist(myPath string) bool {
	_, err := os.Stat(myPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// 获取dirPath下所有以prefix为前缀的文件, 返回的是名字不是路径
func GetDirsStartWith(prefix string, dirPath string) ([]string, *perrorx.ErrorX) {
	ans := make([]string, 0)
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, perrorx.NewErrorXReadDir(dirPath, err.Error(), nil)
	}
	for _, fi := range dir {
		if fi.IsDir() && strings.HasPrefix(fi.Name(), prefix) {
			ans = append(ans, fi.Name())
		}
	}
	return ans, nil
}

// 将文件src拷贝到文件dst, dst不存在会新建, 但文件夹要自己创建
func FileCopy(src string, dst string) *perrorx.ErrorX {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return perrorx.NewErrorXFileCopy(err.Error(), nil)
	}
	if !sourceFileStat.Mode().IsRegular() {
		return perrorx.NewErrorXFileCopy(src + " is not a regular file", nil)
	}
	source, err := os.Open(src)
	if err != nil {
		return perrorx.NewErrorXFileCopy(err.Error(), nil)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return perrorx.NewErrorXFileCopy(err.Error(), nil)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return perrorx.NewErrorXFileCopy(err.Error(), nil)
	}
	return nil
}

// cp -r from to
func CopyR(from, to string) *perrorx.ErrorX {
	_, err := Sh("cp -r " + from + " " + to)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}

// cp from to
func Copy(from, to string) *perrorx.ErrorX {
	_, err := Sh("cp" + from + " " + to)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}