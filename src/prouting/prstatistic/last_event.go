package prstatistic

import (
	"auiauto/perrorx"
	"auiauto/pkernel/panalyze"
	"auiauto/pkernel/pevent"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 统计每条聚焦语句被哪些动作覆盖过, 论文统计用
//
func RLastEvent(c *gin.Context) {
	projectId := c.Query("projectId")
	analyzeFile := c.Query("analyzeFile")
	_, susCodes, err := panalyze.ReadSourceTreeAndNormalAnalyzeStd(projectId, "origin", analyzeFile, "Ochiai")
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RAnalyze error: " + err.MType + ":" + err.MDescription)
		return
	}
	events, err := pevent.ReadEventsStd(projectId, "origin_crash")
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "RAnalyze error: " + err.MType + ":" + err.MDescription)
		return
	}
	//ans := "[" + strconv.Itoa(len(events.MEvents)) + "]"
	//for i := 0; i < susCodes.Len(); i++ {
	//	ans += "["+susCodes[i].MClassName+":"+strconv.Itoa(susCodes[i].MLine)+"]"
	//	id := susCodes[i].MIdx
	//	eventIds := susCodes[i].MOriginNode.MEventIds[id]
	//	for j := 0; j < len(eventIds); j++ {
	//		ans += strconv.Itoa(eventIds[j]) + " "
	//	}
	//	ans += "\n"
	//}
	ans := "0"
	for i := 0; i < susCodes.Len(); i++ {
		id := susCodes[i].MIdx
		if susCodes[i].MOriginNode == nil {
			continue
		}
		eventIds := susCodes[i].MOriginNode.MEventIds[id]
		for j := 0; j < len(eventIds); j++ {
			if eventIds[j] >= len(events.MEvents) {
				ans = "1"
			}
		}
	}
	c.String(http.StatusOK, ans)
}