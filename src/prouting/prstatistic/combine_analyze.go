package prstatistic

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pstatistic"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 统计组合分析结果
func RStatisticCombine(c *gin.Context) {
	projectId := c.Query("projectId")
	analyzeFile := c.Query("analyzeFile")
	factor := c.Query("factor")
	casePrefix := c.Query("casePrefix")
	ans, err := pstatistic.StatisticCombineStd(projectId, analyzeFile, factor, casePrefix)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RAnalyze error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, ans)
}
