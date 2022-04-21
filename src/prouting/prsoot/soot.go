package prsoot

import (
	"auiauto/pdba"
	"auiauto/putils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 调用soot进行插桩
func RSoot(c *gin.Context) {
	database := c.Query("database")
	projectId := c.Query("projectId")
	inputPath := c.Query("inputPath")
	outputPath := c.Query("outputPath")
	sootId := c.Query("sootId")
	ans, err := putils.Sh("java -cp " + pdba.DBURLKernelSoot() + " xmu.wrxlab.abuilder.ABuilderMain " +
		database + " " +
		projectId + " " +
		inputPath + " " +
		outputPath + " " +
		sootId)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	logrus.Errorf("[soot ans]\n" + ans + "\n[soot err]\n" + errStr)
	c.String(http.StatusOK, "ok")
}
