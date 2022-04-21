package prappmanager

import (
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 安装/启动apk, 安装前会先清理环境
func RInstallAndStart(c *gin.Context) {
	avd := c.Query("avd")
	projectId := c.Query("projectId")
	ans, err := padb.AdbInstallAndStart(avd, projectId)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "install start error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, ans)
}

// 重新安装app
func RReInstall(c *gin.Context) {
	avd := c.Query("avd")
	projectId := c.Query("projectId")
	ans, err := padb.AdbOnlyInstall(avd, projectId)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "install error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, ans)
}

// 重新启动app
func RReStart(c *gin.Context) {
	avd := c.Query("avd")
	projectId := c.Query("projectId")
	ans, err := padb.AdbOnlyStart(avd, projectId)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "start error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, ans)
}
