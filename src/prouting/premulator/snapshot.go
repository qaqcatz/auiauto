package premulator

import (
	"auiauto/pconfig"
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"auiauto/putils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
)

// 获取待查询avd下的所有快照, 返回[]string切片
func RSnapshot(c *gin.Context) {
	avd := c.Query("avd")
	avdName, err := padb.AdbEmuAvdName(avd)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "get avd name error: " + err.MType + ":" + err.MDescription)
		return
	}
	snapshotPath := path.Join(pconfig.GConfig.MAvd, avdName, "snapshots")
	if !putils.FileExist(snapshotPath) {
		err := perrorx.NewErrorXFileNotFound(snapshotPath, nil)
		logrus.Errorf(err.Error())
		c.String(http.StatusInternalServerError,  err.MType + ":" + err.MDescription)
		return
	}
	ans, err := putils.GetDirsStartWith("", snapshotPath)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "get snapshots error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, ans)
}

// 加载快照
func RLoadSnapshot(c *gin.Context) {
	avd := c.Query("avd")
	name := c.Query("name")
	ans, err := padb.AdbEmuLoadSnapshot(avd, name)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "load snapshot error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, ans)
}

// 保存快照
func RSaveSnapshot(c *gin.Context) {
	avd := c.Query("avd")
	name := c.Query("name")
	ans, err := padb.AdbEmuSaveSnapshot(avd, name)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "save snapshot error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, ans)
}

// 删除快照
func RDeleteSnapshot(c *gin.Context) {
	avd := c.Query("avd")
	name := c.Query("name")
	ans, err := padb.AdbEmuDeleteSnapshot(avd, name)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError,  "delete snapshot error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, ans)
}