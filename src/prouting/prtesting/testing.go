package prtesting

import (
	"auiauto/perrorx"
	"auiauto/pkernel/ptesting"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// 根据参数启动相应的测试器进行测试
func RTesting(c *gin.Context) {
	tester := c.Query("tester")
	avd := c.Query("avd")
	numStr := c.Query("testNum")
	num, err_ := strconv.Atoi(numStr)
	if err_ != nil{
		err := perrorx.NewErrorXAtoI(numStr, nil)
		logrus.Errorf(err.Error())
		c.String(http.StatusInternalServerError, "RTesting: " + err.MType + ":" + err.MDescription)
		return
	}
	projectId := c.Query("projectId")
	crashCase := c.Query("crashCase")
	testPrefix := c.Query("testPrefix")
	param := c.Query("testParam")

	var err *perrorx.ErrorX = nil
	switch tester {
	case "monkey":
		err = ptesting.Testing(ptesting.NewMonkeyTester(avd, num, projectId, crashCase, testPrefix, param))
	default:
		err = perrorx.NewErrorXRTesting("unknown tester " + tester, nil)
	}

	if err != nil {
		logrus.Errorf(err.Error())
		c.String(http.StatusInternalServerError, "RTesting: " + err.MType + ":" + err.MDescription)
		return
	}

	c.String(http.StatusOK, "success!")
}