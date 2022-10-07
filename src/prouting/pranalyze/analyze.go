package pranalyze

import (
	"auiauto/perrorx"
	"auiauto/pkernel/panalyze"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 根据casePrefix以及analyzeFile进行差异分析.
func RAnalyze(c *gin.Context) {
	projectId := c.Query("projectId")
	casePrefix := c.Query("casePrefix")
	analyzeFile := c.Query("analyzeFile")
	factor := c.Query("factor")
	sourceTree, susCodes, err := panalyze.ReadSourceTreeAndNormalAnalyzeStd(projectId, casePrefix, analyzeFile, factor)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RAnalyze error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, SourceTreeAndSusCodes{sourceTree, susCodes})
}