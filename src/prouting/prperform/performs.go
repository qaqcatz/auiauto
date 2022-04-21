package prperform

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pperform"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 复现events
func RPerforms(c *gin.Context) {
	avd := c.Query("avd")
	projectId := c.Query("projectId")
	caseName := c.Query("caseName")

	// 复现events
	_, err, flag := pperform.PostPerforms(avd, projectId, caseName)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "post performs error: "+ err.MType + ":" + err.MDescription)
		return
	}
	crashStr := "uncrash"
	if flag {
		crashStr = "crash"
	}
	c.String(http.StatusOK, crashStr)
}
