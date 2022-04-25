package prstatistic

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pevent"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// 获取一个用例的事件长度, 论文统计用
func REventsLength(c *gin.Context) {
	projectId := c.Query("projectId")
	caseName := c.Query("caseName")
	events, err := pevent.ReadEventsStd(projectId, caseName)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "load error(read events): " + err.MType + ":" + err.MDescription)
		return
	}
	c.String(http.StatusOK, strconv.Itoa(len(events.MEvents)))
}