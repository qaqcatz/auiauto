package premulator

import (
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 获取屏幕截图的base64编码
func RScreenShot(c *gin.Context) {
	avd := c.Query("avd")
	data, err := padb.GetScreenShot(avd)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "get screenshot error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, data)
}