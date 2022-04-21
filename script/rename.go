package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	for _, project := range Projects {
		dirPath := path.Join("../database/projects", project, "testcases")
		oldFilePath := path.Join(dirPath, "analyze_art_root_cause_Ochiai.txt")
		newFilePath := path.Join(dirPath, "analyze_art_rootcause_Ochiai.txt")
		err := os.Rename(oldFilePath, newFilePath)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
