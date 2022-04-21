package prperform

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pperform"
	"auiauto/pkernel/pevent"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// 执行click/longclick/edit/scroll/keyevent/swipe操作
func RPerform(c *gin.Context) {
	avd := c.Query("avd")
	// 解析event
	var event pevent.Event
	err := c.BindJSON(&event)
	if err != nil {
		err_ := perrorx.NewErrorXBindJson(err.Error(), nil)
		logrus.Errorf(err_.Error())
		c.String(http.StatusInternalServerError, "perform error(bind json): " + err_.MType + ":" + err_.MDescription)
		return
	}
	// 发送post请求(PostPerform会使用http+adb执行操作)
	_, err_ := pperform.PostPerform(avd, &event, false)
	if err_ != nil {
		logrus.Errorf(perrorx.TransErrorX(err_).Error())
		c.String(http.StatusInternalServerError, "perform error(post): " + err_.MType + ":" + err_.MDescription)
		return
	}
	// 睡一会, 防止前端立即dump获取更新前的界面
	logrus.Debugf("RPerform: waitAfter: 750ms")
	time.Sleep(time.Millisecond*750)
	c.String(http.StatusOK, "perform")
}
