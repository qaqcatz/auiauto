package main

import (
	"auiauto/pconfig"
	"auiauto/prouting/pranalyze"
	"auiauto/prouting/prappmanager"
	"auiauto/prouting/prcharts"
	"auiauto/prouting/prcoverage"
	"auiauto/prouting/premulator"
	"auiauto/prouting/prperform"
	"auiauto/prouting/prproject"
	"auiauto/prouting/prsoot"
	"auiauto/prouting/prstatistic"
	"auiauto/prouting/prstmtlog"
	"auiauto/prouting/prtestcase"
	"auiauto/prouting/prtesting"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	// 读取配置文件
	pconfig.InitConfig()
	// 配置 logrus
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	// 配置gin
	r := gin.Default()
	r.LoadHTMLGlob("www/html/*")
	r.Static("/css", "www/css")
	r.Static("/js", "www/js")
	r.Static("/favicon.ico", "www/favicon.ico")

	// 主页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auiauto.html", "hzy")
	})
	// 数据统计页面
	r.GET("/sta", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auiautosta.html", "hzysta")
	})

	// emulator
	// 获取当前可用的avd列表
	r.GET("/devices", premulator.RDevices)
	// antrance中存在日志时会禁止用户后续的操作, 需要init才能正常使用
	r.GET("/init", premulator.RInit)
	// 获取屏幕截图的base64编码
	r.GET("/screenshot", premulator.RScreenShot)
	// snapshot
	r.GET("/snapshot", premulator.RSnapshot)
	r.GET("/loadsnapshot", premulator.RLoadSnapshot)
	r.GET("/savesnapshot", premulator.RSaveSnapshot)
	r.GET("/deletesnapshot", premulator.RDeleteSnapshot)
	// 获取uitree
	r.GET("/uitree", premulator.RUITree)

	// appmanager
	// reinstall and start
	r.GET("/installstart", prappmanager.RInstallAndStart)
	// only reinstall
	r.GET("/reinstall", prappmanager.RReInstall)
	// only restart
	r.GET("/restart", prappmanager.RReStart)

	// perform
	// 单步复现
	r.POST("/perform", prperform.RPerform)
	// 复现用例
	r.POST("/performs", prperform.RPerforms)

	// testcase
	// 加载project下的全部用例
	r.GET("/testcases", prtestcase.RLoadTestCases)
	// save case
	r.POST("/savecase", prtestcase.RSaveCase)
	// load case
	r.GET("/loadcase", prtestcase.RLoadCase)

	// project
	// 加载全部projects
	r.GET("/projects", prproject.RLoadProjects)

	// coverage
	// 根据caseName计算覆盖率, 通过eventId指定想要获取哪个动作的覆盖率(为空表示获取全部动作的覆盖率)
	r.GET("/cover", prcoverage.RCover)
	// 获取总覆盖率以及每个动作的覆盖率, 以切片的形式返回
	r.GET("coverAll", prcoverage.RCoverAll)
	// 获取指定源文件的语句及覆盖, ccl:codes and cover lines
	r.GET("/ccl", prcoverage.RCCL)

	// 差异分析
	// analyze
	r.GET("/analyze", pranalyze.RAnalyze)
	// combine analyze
	r.GET("/combineanalyze", pranalyze.RCombineAnalyze)
	// rd analyze
	r.GET("/rdanalyze", pranalyze.RRDAnalyze)

	// soot
	r.GET("/soot", prsoot.RSoot)

	// get stmt log now
	r.GET("/stmtlognow", prstmtlog.RStmtlogNow)

	// 测试
	r.GET("/testing", prtesting.RTesting)

	// 统计
	r.GET("/stardanalyze", prstatistic.RStatisticRd)
	r.GET("/stacombineanalyze", prstatistic.RStatisticCombine)
	r.GET("/stafactorsanalyze", prstatistic.RStatisticFactors)
	r.GET("/staanalyzeall", prstatistic.RStatisticAll)
	r.GET("/stardtesting", prstatistic.RStatisticRdTesting)
	r.GET("/projectlines", prstatistic.RProjectLines)
	r.GET("/eventslength", prstatistic.REventsLength)
	r.GET("/lastevent", prstatistic.RLastEvent)
	// 保存表格
	r.POST("/savecharts", prcharts.RSaveCharts)

	logrus.Infof("auiauto start on " + pconfig.GConfig.MIp + ":" + pconfig.GConfig.MPort)
	_ = r.Run(pconfig.GConfig.MIp +":"+ pconfig.GConfig.MPort)
}