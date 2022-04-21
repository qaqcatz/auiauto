package prstatistic

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pstatistic"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RStatisticRdTesting(c *gin.Context) {
	projectId := c.Query("projectId")
	tester := c.Query("tester")
	casePrefix := c.Query("casePrefix")
	ans, err := pstatistic.RDTestingAnalyzeStd(projectId, tester, casePrefix)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RStatisticRdTesting: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, ans)
}