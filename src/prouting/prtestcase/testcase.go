package prtestcase

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/pevent"
	"auiauto/putils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 在projectId下保存caseName用例, 没有项目新建项目, caseName存在拒绝覆盖
func RSaveCase(c *gin.Context) {
	projectId := c.Query("projectId")
	caseName := c.Query("caseName")
	var events pevent.Events
	err := c.BindJSON(&events)
	if err != nil {
		err_ := perrorx.NewErrorXBindJson(err.Error(), nil)
		logrus.Errorf(err_.Error())
		c.String(http.StatusInternalServerError, "save error(bingJSON): " + err_.MType + ":" + err_.MDescription)
		return
	}
	err_ := pevent.WriteEventsStd(projectId, caseName, &events)
	if err_ != nil {
		logrus.Errorf(perrorx.TransErrorX(err_).Error())
		c.String(http.StatusInternalServerError, "save error(write file): " + err_.MType + ":" + err_.MDescription)
		return
	}
	c.String(http.StatusOK, "save")
}

// load case
func RLoadCase(c *gin.Context) {
	projectId := c.Query("projectId")
	caseName := c.Query("caseName")
	events, err := pevent.ReadEventsStd(projectId, caseName)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "load error(read events): " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, events)
}

// 加载project下的testcases
func RLoadTestCases(c *gin.Context) {
	project := c.Query("project")
	ans, err := putils.GetDirsStartWith("", pdba.DBURLProjectIdTestcases(project))
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "read project's testcase error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, ans)
}
