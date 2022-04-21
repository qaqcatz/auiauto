package prcharts

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pcharts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type base64Data struct {
	MBase64 string `json:"base64"`
}

// 保存表格
func RSaveCharts(c *gin.Context) {
	casePrefix := c.Query("casePrefix")
	analyzeType := c.Query("analyzeType")
	analyzeFile := c.Query("analyzeFile")
	factor := c.Query("factor")
	projectId := c.Query("projectId")
	var b64Data base64Data
	err := c.BindJSON(&b64Data)
	if err != nil {
		err_ := perrorx.NewErrorXBindJson(err.Error(), nil)
		logrus.Errorf(err_.Error())
		c.String(http.StatusInternalServerError, "save error(bingJSON): " + err_.MType + ":" + err_.MDescription)
		return
	}
	err_ := pcharts.SaveCharts(casePrefix, analyzeType, analyzeFile, factor, projectId, b64Data.MBase64)
	if err_ != nil {
		logrus.Errorf(perrorx.TransErrorX(err_).Error())
		c.String(http.StatusInternalServerError, "RCombineAnalyze error: " + err_.MType + ":" + err_.MDescription)
		return
	}
	c.JSON(http.StatusOK, "save png successful")
}
