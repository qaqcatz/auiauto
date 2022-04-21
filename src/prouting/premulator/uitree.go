package premulator

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pahttp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 获取uitree的xml编码
func RUITree(c *gin.Context) {
	avd := c.Query("avd")
	data, err := pahttp.GetUITree(avd)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "get ui tree error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, data)
}