package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

// copy avd
// be careful of the name of src
func cpAvd(avdPath string, src string, dst string) error {
	fmt.Println("cp " + src + " to " + dst)
	// cp -r src.avd dst.avd
	srcAvdPath := path.Join(avdPath, src+".avd")
	dstAvdPath := path.Join(avdPath, dst+".avd")
	if FileExist(dstAvdPath) {
		return errors.New(dstAvdPath + " exists")
	}
	err := CopyR(srcAvdPath, dstAvdPath)
	if err != nil {
		return err
	}
	// cp src.ini dst.ini
	srcIniPath := path.Join(avdPath, src+".ini")
	dstIniPath := path.Join(avdPath, dst+".ini")
	err = FileCopy(srcIniPath, dstIniPath)
	if err != nil {
		return err
	}
	// strings.Replace(xxx.ini, src, dst), be careful of the name of src
	inis := make([]string, 0)
	inis = append(inis, dstIniPath)
	inis = append(inis, path.Join(dstAvdPath, "config.ini"))
	inis = append(inis, path.Join(dstAvdPath, "hardware-qemu.ini"))
	dstSnapshotPath := path.Join(dstAvdPath, "snapshots")
	if FileExist(dstSnapshotPath) {
		dir, err_ := ioutil.ReadDir(dstSnapshotPath)
		if err_ != nil {
			return err
		}
		for _, fi := range dir {
			if fi.IsDir() {
				inis = append(inis, path.Join(dstSnapshotPath, fi.Name(), "hardware.ini"))
			}
		}
	}
	for _, ini := range inis {
		if FileExist(ini) {
			data, err_ := ioutil.ReadFile(ini)
			if err_ != nil {
				return err
			}
			newData := []byte(strings.ReplaceAll(string(data), src, dst))
			err_ = ioutil.WriteFile(ini, newData, 0777)
			if err_ != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	avdPath := "/home/android/.android/avd"
	for i := 43; i <= 43; i++ {
		err := cpAvd(avdPath, "auiauto0", "auiauto" + strconv.Itoa(i))
		if err != nil {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println(err.Error())
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
	}
}

