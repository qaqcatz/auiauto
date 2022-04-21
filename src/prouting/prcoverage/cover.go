package prcoverage

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pevent"
	"auiauto/pkernel/psrctree"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// 根据caseName计算覆盖率, 通过eventId指定想要获取哪个动作的覆盖率(为空表示获取全部动作的覆盖率)
func RCover(c *gin.Context) {
	projectId := c.Query("projectId")
	caseName := c.Query("caseName")
	eventId := c.Query("eventId")

	sourceTree, err := psrctree.ReadSourceTreeAndCover(projectId, caseName, eventId)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "ReadSourceTreeAndCover error: " + err.MType + ":" + err.MDescription)
		return
	}

	c.JSON(http.StatusOK, sourceTree)
}

// 获取总覆盖率以及每个动作的覆盖率, 以切片的形式返回
func RCoverAll(c *gin.Context) {
	projectId := c.Query("projectId")
	caseName := c.Query("caseName")
	// 读取测试用例
	events, err := pevent.ReadEventsStd(projectId, caseName)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "ReadEventsStd error: " + err.MType + ":" + err.MDescription)
		return
	}
	ans := make([]string, 0)
	// 统计总覆盖率
	sourceTree, err := psrctree.ReadSourceTreeAndCover(projectId, caseName, "")
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "ReadSourceTreeAndCover error: " + err.MType + ":" + err.MDescription)
		return
	}
	ans = append(ans, "[all]"+strconv.Itoa(sourceTree.MRoot.MCoverNum)+"/"+
		strconv.Itoa(sourceTree.MRoot.MTotalNum))
	// 遍历每个event, 通过ReadSourceTreeAndCover统计i!/i的覆盖情况
	for i := 0; i < len(events.MEvents); i++ {
		x_y := "["+strconv.Itoa(i+1)+"]"
		sourceTree, err := psrctree.ReadSourceTreeAndCover(projectId, caseName, strconv.Itoa(i+1)+"!")
		if err != nil {
			logrus.Errorf(perrorx.TransErrorX(err).Error())
			c.String(http.StatusInternalServerError, "ReadSourceTreeAndCover error: " + err.MType + ":" + err.MDescription)
			return
		}
		x_y += strconv.Itoa(sourceTree.MRoot.MCoverNum)+"/"
		sourceTree, err = psrctree.ReadSourceTreeAndCover(projectId, caseName, strconv.Itoa(i+1))
		if err != nil {
			logrus.Errorf(perrorx.TransErrorX(err).Error())
			c.String(http.StatusInternalServerError, "ReadSourceTreeAndCover error: " + err.MType + ":" + err.MDescription)
			return
		}
		x_y += strconv.Itoa(sourceTree.MRoot.MCoverNum)
		ans = append(ans, x_y)
	}

	c.JSON(http.StatusOK, ans)
}

// 获取指定源文件的语句及覆盖, ccl:codes and cover lines
func RCCL(c *gin.Context) {
	projectId := c.Query("projectId")
	dotClassPath := c.Query("dotClassPath")
	if psrctree.GetGSourceTreeJS() == nil || psrctree.GetGSourceTreeJS().MRoot.MName != projectId {
		err := perrorx.NewErrorXCCL("GSourceTree empty or unmatch", nil)
		logrus.Errorf(err.Error())
		c.String(http.StatusInternalServerError, "CCL error: " + err.MType + ":" + err.MDescription)
		return
	}
	codesAndCoverLines, err := psrctree.GetGSourceTreeJS().GetCodesAndCoverLines(dotClassPath)
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "utils.GSourceTree.GetCodesAndCoverLines error: " + err.MType + ":" + err.MDescription)
		return
	}
	c.JSON(http.StatusOK, codesAndCoverLines)
}