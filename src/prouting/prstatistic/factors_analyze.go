package prstatistic

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pstatistic"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 统计因子分析结果
func RStatisticFactors(c *gin.Context) {
	analyzeFile := c.Query("analyzeFile")
	casePrefix := c.Query("casePrefix")
	ans, err := pstatistic.StatisticFactorsStd(analyzeFile, casePrefix)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RAnalyze error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, ans)
}