package prperform

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pperform"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// 复现events
func RPerformEs(c *gin.Context) {
	avd := c.Query("avd")
	projectId := c.Query("projectId")
	caseName := c.Query("caseName")

	// 复现events
	i, err := pperform.PostPerformEs(avd, projectId, caseName)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "post performs error: "+ err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, strconv.Itoa(i))
}

