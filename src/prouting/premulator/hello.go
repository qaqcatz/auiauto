package premulator

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pahttp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RHello(c *gin.Context) {
	avd := c.Query("avd")
	ans, err := pahttp.GetHello(avd)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "get hello error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, ans)
}