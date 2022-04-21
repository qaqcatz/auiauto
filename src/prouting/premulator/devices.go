package premulator

import (
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 获取当前可用的模拟器列表, 相当于执行adb devices
// 返回一个[]string, 数组的每一项为一个可用的模拟器名
func RDevices(c *gin.Context) {
	data, err := padb.AdbDevices()
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "adb devices error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, data)
}