package prstatistic

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pstatistic"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 统计组合分析结果
func RStatisticAll(c *gin.Context) {
	analyzeFile := c.Query("analyzeFile")
	factor := c.Query("factor")
	tester := c.Query("tester")
	ans, err := pstatistic.StatisticAllStd(analyzeFile, factor, tester)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RAnalyze error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, ans)
}
