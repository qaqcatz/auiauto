package pranalyze

import (
	"auiauto/perrorx"
	"auiauto/pkernel/panalyze"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 组合casePrefix下的用例进行差异分析, 生成分析报告
func RCombineAnalyze(c *gin.Context) {
	projectId := c.Query("projectId")
	casePrefix := c.Query("casePrefix")
	analyzeFile := c.Query("analyzeFile")
	factor := c.Query("factor")
	err := panalyze.ReadSourceTreeAndCombineAnalyzeStd(projectId, casePrefix, analyzeFile, factor)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RCombineAnalyze error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, "please check "+projectId+"/testcases/"+casePrefix+"_"+analyzeFile+"combine_analyze.txt")
}
