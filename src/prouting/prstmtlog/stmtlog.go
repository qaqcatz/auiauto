package prstmtlog

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/pcoverage"
	"auiauto/putils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// 获取当前的stmtlog
func RStmtlogNow(c *gin.Context) {
	avd := c.Query("avd")
	projectId := c.Query("projectId")
	caseName := c.Query("caseName")

	testcasesPath := pdba.DBURLProjectIdTestcases(projectId)
	casePath := pdba.DBURLProjectIdTestcase(projectId, caseName)
	// 没有testcases新建一个
	if !putils.FileExist(testcasesPath) {
		_ = os.Mkdir(testcasesPath, 0777)
	}
	// caseName不存在新建
	if !putils.FileExist(casePath) {
		_ = os.MkdirAll(casePath, 0777)
	}

	stmtLog, err := pcoverage.SaveStmtLogAndLineCoverageStd(avd, projectId, caseName)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RStmtlogNow: " + err.MType + ":" + err.MDescription)
		return
	}

	ans := "uncrash"
	if stmtLog.MStatus == "false" {
		ans = "crash"
	}

	c.String(http.StatusOK, ans)
}
