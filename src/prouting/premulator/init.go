package premulator

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pahttp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 初始化antrance, 主要是清空antrance的stmtlog
func RInit(c *gin.Context) {
	avd := c.Query("avd")
	err := pahttp.GetInit(avd)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "get init error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, "init")
}
