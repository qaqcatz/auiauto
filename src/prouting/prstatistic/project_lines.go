package prstatistic

import (
	"auiauto/perrorx"
	"auiauto/pkernel/psrctree"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// 返回项目总行数, 论文结果统计时用
func RProjectLines(c *gin.Context) {
	projectId := c.Query("projectId")
	caseName := c.Query("caseName")

	sourceTree, err := psrctree.ReadSourceTreeAndCover(projectId, caseName, "")
	if err != nil {
		logrus.Errorf(perrorx.TransErrorX(err).Error())
		c.String(http.StatusInternalServerError, "ReadSourceTreeAndCover error: " + err.MType + ":" + err.MDescription)
		return
	}

	c.String(http.StatusOK, strconv.Itoa(sourceTree.MRoot.MTotalNum))
}
