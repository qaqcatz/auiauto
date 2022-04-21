package pdba

import (
	"auiauto/pconfig"
	"path"
)

// 数据库结构:
//	database
//	--kernel: 核心文件
//	----abuilder, gradle插桩插件, android-gradle版本2~3时使用0.x版本, >3时使用1+版本, 不支持<2, 源码可去qaqcatz的github仓库中的antrance项目获取
//	----antrance.apk, Accessibility Service, 负责协助adb执行部分动作, 以及收集各个app的覆盖率信息, 源码可去qaqcatz的github仓库中的antrance项目获取
//	----soot.jar, *.class, 通过soot将*.class注入到app中, 包括与antrance进行通信的覆盖率日志采集服务, 崩溃捕获等,
//	修改后的soot源码可去qaqcatz的github仓库中的soot项目获取, *.class是antrance项目生成的, 可以查看antrance项目下的kernel.sh文件查看这几个class的来源
//	----jar.sh, 当soot进行jar包分析时需要使用jar命令打包解包, 注意如果没有在usr/bin下配置过jar命令的话需要手动修改这个文件的jar路径
//	--tmp: 临时文件, 比如屏幕截图, 一些统计图表等, 都是不重要的数据, 放心删除
//	--projects: 待分析的项目, 在android studio使用gradle插件abuilder, 执行abuilder任务可自动在projects下生成项目文件
//	(0.x版本还需要手动复制apk, 手动编写apk元信息, 1+版本可全自动实现整个流程)
//	----apk: 保存项目的插桩后apk文件和元信息
//	------app.apk: 项目的插桩后apk文件
//	------config.txt: apk元信息, 格式: applicationId@mainActivity, 主要用于告诉adb怎么启动apk
//	----cfg: soot生成的控制流图
//	----classes: soot分析过程中的临时class目录, 可删除
//	----jars: soot分析过程中的临时jar目录, 存储着解包后的class文件, 可删除
//	----debugjimple: soot分析过程中生成的jimple中间代码, 主要用于开发者调试, 可删除
//	----src: 项目源码
//	----clssrcmap: 记录源文件与class文件的对应关系, 覆盖率展示时会用到
//	----logIdSig.txt: 我们插桩时使用了数组id优化, 需要把数组id还原成实际的语句, logIdSig记录了这个对应关系
//	----testcases: 保存用户录制的测试用例, 尽管没有在程序中限制, 但强烈推荐采用前缀命名法: 自定义前缀+'_'+crash/pass+自定义id
//	分析时我们采用的是前缀统计法, 即用户输入一个前缀, 我们分析所有包含这个前缀的文件.
//	另外有时我们也会判断文件名中是否包含crash/pass来判断错误/正确用例, 因此推荐采用这种方式命名
//	------cover.json: 覆盖率文件, 格式可以去src/pkernel/pcoverage/coverage.go的Coverage类查看
//	------stmtlog.json: 原始覆盖率文件(数组id没有映射到语句), 开发者调试用
//  ------testcase.json: 保存用户录制的用例, 格式可以参考src/pkernel/pevent/events.go下的Events类
//	    分析文件保存为*.txt, 格式可以去src/pkernel/panalyze下的各个子模块查看, 在文件开头有相应注释
//	    聚焦的分析语句保存未*.json, 比如fix.json是聚焦fixing commit中的语句, rootcause.json是聚焦我们分析出的root cause语句, 学术研究用.
//	平时用top100就可以, 不需要提供聚焦语句.
//	----test: 保存测试软件录制的测试用例, 目前auiauto自带的测试工具aua已经停止维护, 现在只有monkey可用.
//	下面是各个tester, 每个tester下是这个tester生成的各个用例, 命名格式采用严格的前缀命名法, 不要做任何修改.
//	另外tester的执行日志可在*.log中查看, 分析结果可以在*.txt中查看.

// 这里主要提供数据库url服务, 开发过程中应该严格用这里的url访问数据库, 不要自己拼接路径

func DBURLKernel() string {
	return path.Join(pconfig.GConfig.MDatabase, "kernel")
}

func DBURLKernelAntranceAPK() string {
	return path.Join(DBURLKernel(), "antrance.apk")
}

func DBURLKernelSoot() string {
	return path.Join(DBURLKernel(), "soot.jar")
}

func DBURLTmp() string {
	return path.Join(pconfig.GConfig.MDatabase, "tmp")
}

func DBURLTmpScreenshot() string {
	return path.Join(DBURLTmp(), "screenshot.png")
}

func DBURLTmpCharts() string {
	return path.Join(DBURLTmp(), "charts")
}

func DBURLProjects() string {
	return path.Join(pconfig.GConfig.MDatabase, "projects")
}

func DBURLProjectId(projectId string) string {
	return path.Join(DBURLProjects(), projectId)
}

func DBURLProjectIdAPK(projectId string) string {
	return path.Join(DBURLProjectId(projectId), "apk")
}

func DBURLProjectIdAPKAPP(projectId string) string {
	return path.Join(DBURLProjectIdAPK(projectId), "app.apk")
}

func DBURLProjectIdAPKConfig(projectId string) string {
	return path.Join(DBURLProjectIdAPK(projectId), "config.txt")
}

func DBURLProjectIdCFG(projectId string) string {
	return path.Join(DBURLProjectId(projectId), "cfg")
}

func DBURLProjectIdSRC(projectId string) string {
	return path.Join(DBURLProjectId(projectId), "src")
}

func DBURLProjectIdTest(projectId string) string {
	return path.Join(DBURLProjectId(projectId), "test")
}

func DBURLProjectIdTester(projectId string, tester string) string {
	return path.Join(DBURLProjectIdTest(projectId), tester)
}

func DBURLProjectIdTesterTestCase(projectId string, tester string, caseName string) string {
	return path.Join(DBURLProjectIdTester(projectId, tester), caseName)
}

func DBURLProjectIdTesterTestCaseCover(projectId string, tester string, caseName string) string {
	return path.Join(DBURLProjectIdTesterTestCase(projectId, tester, caseName), "cover.json")
}

func DBURLProjectIdTesterTestCaseStmtlog(projectId string, tester string, caseName string) string {
	return path.Join(DBURLProjectIdTesterTestCase(projectId, tester, caseName), "stmtlog.json")
}

func DBURLProjectIdTestcases(projectId string) string {
	return path.Join(DBURLProjectId(projectId), "testcases")
}

func DBURLProjectIdTestcase(projectId string, caseName string) string {
	return path.Join(DBURLProjectIdTestcases(projectId), caseName)
}

func DBURLProjectIdTestcaseCover(projectId string, caseName string) string {
	return path.Join(DBURLProjectIdTestcase(projectId, caseName), "cover.json")
}

func DBURLProjectIdTestcaseStmtlog(projectId string, caseName string) string {
	return path.Join(DBURLProjectIdTestcase(projectId, caseName), "stmtlog.json")
}

func DBURLProjectIdTestcaseTestcase(projectId string, caseName string) string {
	return path.Join(DBURLProjectIdTestcase(projectId, caseName), "testcase.json")
}

func DBURLProjectIdClssrcmap(projectId string) string {
	return path.Join(DBURLProjectId(projectId), "clssrcmap")
}

func DBURLProjectIdLogidsig(projectId string) string {
	return path.Join(DBURLProjectId(projectId), "logIdSig.txt")
}