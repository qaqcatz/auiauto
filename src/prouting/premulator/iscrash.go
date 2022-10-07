package premulator

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pahttp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RIsCrash(c *gin.Context) {
	avd := c.Query("avd")
	ans, err := pahttp.GetIsCrash(avd)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "get is crash error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, ans)
}
