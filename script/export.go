package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	os.MkdirAll("../database-cp/projects", 0777)
	for _, project := range Projects {
		os.Mkdir(path.Join("../database-cp/projects", project), 0777)
		testPath := path.Join("../database/projects", project, "test")
		_, err := Sh("cp -r "+testPath+" "+path.Join("../database-cp/projects", project, "test"))
		if err != nil {
			fmt.Println(err.Error())
			break
		} else {
			fmt.Println("cp "+testPath)
		}
		testcasePath := path.Join("../database/projects", project, "testcases")
		_, err = Sh("cp -r "+testcasePath+" "+path.Join("../database-cp/projects", project, "testcases"))
		if err != nil {
			fmt.Println(err.Error())
			break
		} else {
			fmt.Println("cp "+testcasePath)
		}
		srcPath := path.Join("../database/projects", project, "src")
		_, err = Sh("cp -r "+srcPath+" "+path.Join("../database-cp/projects", project, "src"))
		if err != nil {
			fmt.Println(err.Error())
			break
		} else {
			fmt.Println("cp "+srcPath)
		}
		clssrcmapPath := path.Join("../database/projects", project, "clssrcmap")
		_, err = Sh("cp "+clssrcmapPath+" "+path.Join("../database-cp/projects", project, "clssrcmap"))
		if err != nil {
			fmt.Println(err.Error())
			break
		} else {
			fmt.Println("cp "+clssrcmapPath)
		}
	}
}

