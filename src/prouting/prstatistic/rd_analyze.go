package prstatistic

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pstatistic"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RStatisticRd(c *gin.Context) {
	projectId := c.Query("projectId")
	analyzeFile := c.Query("analyzeFile")
	factor := c.Query("factor")
	tester := c.Query("tester")
	casePrefix := c.Query("casePrefix")
	ans, err := pstatistic.StatisticRdStd(projectId, analyzeFile, factor, tester, casePrefix)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RAnalyze error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, ans)
}
