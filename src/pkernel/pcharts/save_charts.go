package pcharts

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// 将图表数据存储在tmp/charts/casePrefix_rd|combine_analyzeFile_factor_(3D)projectId.png
func SaveCharts(casePrefix string, analyzeType string, analyzeFile string, factor string, projectId string,
	imgBase64 string) *perrorx.ErrorX {
	if strings.Contains(imgBase64, "data:image/png;base64,") {
		imgBase64 = strings.Replace(imgBase64, "data:image/png;base64,", "", 1)
	} else {
		return perrorx.NewErrorXSaveCharts("imgBase64.prefix != data:image/png;base64,", nil)
	}
	img, _ := base64.StdEncoding.DecodeString(imgBase64)
	imgDir := path.Join(pdba.DBURLTmpCharts(), analyzeType+"_"+casePrefix+"_"+analyzeFile+"_"+factor)
	os.MkdirAll(imgDir, 0777)
	imgPath := path.Join(imgDir, projectId+".png")
	err := ioutil.WriteFile(imgPath, img, 0777)
	if err != nil {
		return perrorx.NewErrorXWriteFile(imgPath, err.Error(), nil)
	}
	return nil
}