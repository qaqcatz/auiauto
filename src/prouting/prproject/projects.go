package prproject

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/putils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 加载全部projects
func RLoadProjects(c *gin.Context) {
	temp, err := putils.GetDirsStartWith("", pdba.DBURLProjects())
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "read projects error: " + err.MType + ":" + err.MDescription)
		return
	}
	ans := make([]string, 0)
	for _, t := range temp {
		if t != "kernel" {
			ans = append(ans, t)
		}
	}
	c.JSON(http.StatusOK, ans)
}
