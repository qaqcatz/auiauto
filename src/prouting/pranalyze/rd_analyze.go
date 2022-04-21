package pranalyze

import (
	"auiauto/perrorx"
	"auiauto/pkernel/panalyze"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// 随机采样分析
func RRDAnalyze(c *gin.Context) {
	projectId := c.Query("projectId")
	analyzeFile := c.Query("analyzeFile")
	factor := c.Query("factor")
	tester := c.Query("tester")
	casePrefix := c.Query("casePrefix")
	randNumStr := c.Query("randNum")
	randNum, err := strconv.Atoi(randNumStr)
	if err != nil{
		err_ := perrorx.NewErrorXAtoI(randNumStr, nil)
		logrus.Errorf(err_.Error())
		c.String(http.StatusInternalServerError, "RTesting: " + err_.MType + ":" + err_.MDescription)
		return
	}
	err_ := panalyze.RDAnalyzeStd(projectId, analyzeFile, factor, tester, casePrefix, randNum)
	if err_ != nil{
		logrus.Errorf(perrorx.TransErrorX(err_).Error())
		c.String(http.StatusInternalServerError, "RTesting: " + err_.MType + ":" + err_.MDescription)
		return
	}
	c.JSON(http.StatusOK, "success")
}