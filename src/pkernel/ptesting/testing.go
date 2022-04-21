package ptesting

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"auiauto/pkernel/pcoverage"
	"auiauto/pkernel/pevent"
	"auiauto/pkernel/pperform"
	"auiauto/pkernel/pstmtlog"
	"auiauto/putils"
	"bufio"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// 测试接口
type AndroidTester interface {
	// 获取本次测试对应的模拟器
	GetAvd() string
	// 获取测试工具名字, 如monkey
	GetTestGenName() string
	// 获取要生成的测试用例总消耗, +表示总数, -表示总时间(ms)
	GetCost() int
	// 获取本次测试的项目id
	GetProjectId() string
	// 获取本次测试的基础crash case
	GetCrashCase() string
	// 获取本次测试的前缀标识, 如random
	GetTestPrefix() string
	// 单次测试
	Start(int) (string, *perrorx.ErrorX)
}

// 在projectId/test下创建测试文件夹testGenName
// 已存在的话会自动覆盖, 这一点要注意
func createTestDir(projectId string, testGenName string) *perrorx.ErrorX {
	testGenPath := pdba.DBURLProjectIdTester(projectId, testGenName)
	// 没有testGenPath新建
	if !putils.FileExist(testGenPath) {
		_ = os.MkdirAll(testGenPath, 0777)
	}
	return nil
}

// 将错误用例复制到projectId/test/testGenName/testPrefix_crash_0
// 返回错误栈信息
func copyCrashCase(projectId string, crashCase string, testGenName string, testPrefix string) (string, *perrorx.ErrorX) {
	// 验证stmtlog是否崩溃
	stmtLog, err := pstmtlog.ReadStmtLogStd(projectId, crashCase)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	if stmtLog.MStatus != "false" {
		return "", perrorx.NewErrorXGenInitSnapBeforeTest("stmtLog.MStatus != false", nil)
	}
	srcPath := pdba.DBURLProjectIdTestcase(projectId, crashCase)
	dstPath := pdba.DBURLProjectIdTesterTestCase(projectId, testGenName, testPrefix+"_crash_0")
	if !putils.FileExist(dstPath) {
		_ = os.MkdirAll(dstPath, 0777)
	}
	err = putils.FileCopy(path.Join(srcPath, "stmtlog.json"), path.Join(dstPath, "stmtlog.json"))
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	err = putils.FileCopy(path.Join(srcPath, "cover.json"), path.Join(dstPath, "cover.json"))
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return stmtLog.MStackTraceOrigin, nil
}

// 根据错误用例的init快照和init testcase创建测试的初始测试快照, 返回初始测试快照的名字
// 一般情况下初始测试快照的名字是projectId_testGenName_crashCase_t
// 如果events有init快照, 没有init用例, 不需要重新安装, 重新启动, 则初始测试快照就是events的init快照
func genInitSnapBeforeTest(avd string, projectId string, testGenName string, crashCase string) (string,
	*perrorx.ErrorX) {
	// 使用ReadEvents获取Events, 创建初始测试快照
	events, err := pevent.ReadEventsStd(projectId, crashCase)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	// 如果events有init快照, 没有init用例, 不需要重新安装, 重新启动, 则初始测试快照就是events的init快照
	if events.MInitSnapshot != "" && events.MInitTestcase == "" &&
		!events.MReInstall && !events.MReStart {
		return events.MInitSnapshot, nil
	}
	// 执行app init操作
	err = pperform.PerformInit(avd, projectId, events)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	// 先等待3s, 让动作执行完毕
	logrus.Debugf(projectId + " genInitSnapBeforeTest: wait 3s before creating snapshot")
	time.Sleep(time.Millisecond*3000)
	// 执行玩app init操作后创建初始测试快照, 名字为projectId_testGenName_crashCase_t
	initTestSnapName := projectId+"_"+testGenName+"_"+crashCase+"_t"
	_, err = padb.AdbEmuSaveSnapshot(avd, initTestSnapName)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return initTestSnapName, nil
}

// 为节省空间, 测试结束后删初始测试快照, 简单封装一下AdbEmuDeleteSnapshot
func deleteInitSanpAfterTest(avd string, snapName string) {
	// 这里一定要注意以_t结尾才能安全删除!!!
	if strings.HasSuffix(snapName, "_t") {
		_, _ = padb.AdbEmuDeleteSnapshot(avd, snapName)
	}
}

// 每轮测试前加载初始快照, 等待3s
// throw TimeoutError if cannot load snapshot in 30s
func loadInitSnapBeforeEachTest(avd string, initSnapName string, projectId string) *perrorx.ErrorX {

	cha := make(chan *perrorx.ErrorX, 1)
	go func() {
		_, err := padb.AdbEmuLoadSnapshot(avd, initSnapName)
		cha <- err
	} ()

	select {
	case err := <-cha:
		if err != nil {
			return perrorx.TransErrorX(err)
		}
	case <-time.After(time.Duration(30000 * time.Millisecond)):
		return perrorx.NewErrorXTimeout(nil)
	}

	// 默认等待3s, 给snapshot准备时间
	logrus.Debugf(projectId + " loadInitSnapBeforeEachTest: wait 3s after loading snapshot")
	time.Sleep(time.Millisecond*3000)
	return nil
}

// used in channel
type stmtLogErrPair struct {
	MStmtLog *pstmtlog.StmtLog
	MErr *perrorx.ErrorX
}

// 每轮测试结束后获取覆盖语句日志, 保存在projectId/test/testGenName/testPrefix_crash|pass_crash|pass id
// 特别地, 对于随机生成的错误用例, 若崩溃栈信息和原始错误用例不相等, 我们将会在用例名最后加X
// throw TimeoutError if cannot get stmt log in 10s
func getStmtLogAfterEachTest(originalStackTrace string, crashPass []int, avd string,
	projectId string, testGenName string, testPrefix string) *perrorx.ErrorX {
	// 先放在临时目录下
	dirTmp := pdba.DBURLProjectIdTesterTestCase(projectId, testGenName, testPrefix)
	// dir不存在新建
	if !putils.FileExist(dirTmp) {
		_ = os.MkdirAll(dirTmp, 0777)
	}
	// 每轮测试结束后默认等待3s, 从而有充足的时间拉日志
	logrus.Debugf(projectId + " getStmtLogAfterEachTest: wait 3s before getting cover")
	time.Sleep(time.Millisecond*3000)

	cha := make(chan stmtLogErrPair, 1)
	go func() {
		// 保存stmtlog和coverage
		stmtLog, err := pcoverage.SaveStmtLogAndLineCoverage(avd, dirTmp, projectId)
		cha <- stmtLogErrPair{stmtLog, err}
	} ()

	var stmtLog *pstmtlog.StmtLog = nil
	select {
	case stmtLogErr := <-cha:
		if stmtLogErr.MErr != nil {
			return perrorx.TransErrorX(stmtLogErr.MErr)
		}
		stmtLog = stmtLogErr.MStmtLog
	case <-time.After(time.Duration(10000 * time.Millisecond)):
		return perrorx.NewErrorXTimeout(nil)
	}

	// 有了日志状态后再对dirTmp重命名
	status := stmtLog.MStatus
	statusStr := "crash"
	if status == "true" {
		statusStr = "pass"
	}
	suffix := "_" + statusStr + "_"
	if status == "false" {
		crashPass[0] += 1
		suffix += strconv.Itoa(crashPass[0])
		// 特别地, 对于随机生成的错误用例, 若崩溃栈信息和原始错误用例不相等, 我们将会在用例名最后加_X
		if stmtLog.MStackTraceOrigin != originalStackTrace {
			suffix += "_X"
		}
	} else {
		crashPass[1] += 1
		suffix += strconv.Itoa(crashPass[1])
	}
	// 上次测试失败的话可能有冗余文件夹, 需要删除
	dir := dirTmp+suffix
	if putils.FileExist(dir) {
		_ = os.RemoveAll(dir)
	}
	err_ := os.Rename(dirTmp, dirTmp+suffix)
	if err_ != nil {
		return perrorx.NewErrorXRename(dirTmp+suffix, nil)
	}
	return nil
}

// 通过遍历projectId/test/testGenName/目录下的测试用例
// 返回crashNum和passNum(由Testing()的crashPass切片接收), 并从相应id继续生成测试用例
// crashNum要排除掉原始错误用例, 因此-1, 注意不要小于0
// 我们通过目录名是否包含crash/pass来统计数量(contains cover.json, stmtlog.json), 用户要注意不要写一些额外的"crash"和"pass"
func continueTesting(projectId string, testGenName string) (int, int) {
	dirPath := pdba.DBURLProjectIdTester(projectId, testGenName)
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return 0, 0
	}
	crashNum := 0
	passNum := 0
	for _, fi := range dir {
		if fi.IsDir() && putils.FileExist(path.Join(dirPath, fi.Name(), "cover.json")) &&
			putils.FileExist(path.Join(dirPath, fi.Name(), "stmtlog.json")) {
			if strings.Contains(fi.Name(), "crash") {
				crashNum++
			}
			if strings.Contains(fi.Name(), "pass") {
				passNum++
			}
		}
	}
	if crashNum > 0 {
		crashNum--
	}
	return crashNum, passNum
}

// 开启测试, 测试过程中会将日志写在projectId/test/testGenName/下保存日志testPrefix.log
// 主要包含各轮迭代的时间戳, 以及执行记录
func Testing(androidTester AndroidTester) *perrorx.ErrorX {
	// 创建测试目录
	err := createTestDir(androidTester.GetProjectId(), androidTester.GetTestGenName())
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	// 复制初始错误用例到测试目录
	originStackTrace, err := copyCrashCase(androidTester.GetProjectId(), androidTester.GetCrashCase(),
		androidTester.GetTestGenName(), androidTester.GetTestPrefix())
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	// 创建初始测试快照
	initSnapName, err := genInitSnapBeforeTest(androidTester.GetAvd(), androidTester.GetProjectId(),
		androidTester.GetTestGenName(), androidTester.GetCrashCase())
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	// 正常/异常结束后删除初始测试快照
	defer deleteInitSanpAfterTest(androidTester.GetAvd(), initSnapName)

	// 生成测试用例

	// projectId/test/testGenName/下保存日志testPrefix.log
	// 主要包含各轮迭代的时间戳, 以及执行记录, 追加写入
	logPath := path.Join(pdba.DBURLProjectIdTester(androidTester.GetProjectId(), androidTester.GetTestGenName()), androidTester.GetTestPrefix()+".log")
	if !putils.FileExist(logPath) {
		_, _ = os.Create(logPath)
	}
	logFile, err_ := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0777)
	if err_ != nil {
		return perrorx.NewErrorXOpen(logPath, nil)
	}
	//及时关闭file句柄
	defer logFile.Close()
	//写入文件时，使用带缓存的 *Writer
	writer := bufio.NewWriter(logFile)

	// 终止条件: 数量或时间, 根据cost的正负来判断
	// 分别记录测试过程中crash, pass的数量, 从而标记其id
	crashPass := make([]int, 2)
	crashPass[0] ,crashPass[1] = continueTesting(androidTester.GetProjectId(), androidTester.GetTestGenName())
	startTime := time.Now()

	// 开始
	logrus.Debugf("[testing] projectId=" + androidTester.GetProjectId() + " start at " + startTime.String())
	_, _ = writer.WriteString("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
	_, _ = writer.WriteString("[testing] startTime=" + startTime.String() + "\n")
	_, _ = writer.WriteString("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
	_ = writer.Flush()
	for {
		i := crashPass[0] + crashPass[1] + 1
		duration := time.Since(startTime)
		curTimStr := time.Now().String()

		if androidTester.GetCost() >= 0 {
			if i >= androidTester.GetCost() {
				break
			}
		} else {
			if duration.Milliseconds() >= int64(-androidTester.GetCost()) {
				break
			}
		}

		logrus.Debugf("[testing] projectId=" + androidTester.GetProjectId() +
			" iteration=" + strconv.Itoa(i) + " (before the iteration:) curTime=" + curTimStr)
		_, _ = writer.WriteString("[testing] iteration=" + strconv.Itoa(i) + " (before the iteration:) crashNum=" + strconv.Itoa(crashPass[0]) +
			" passNum=" + strconv.Itoa(crashPass[1]) + " curTime=" + curTimStr + "\n")
		_ = writer.Flush()

		// crashPass []int, avd string, projectId string, testGenName string, testPrefix string
		// 每轮测试: 加载初始测试快照, start, 获取覆盖语句
		// 加载初始测试快照
		err := loadInitSnapBeforeEachTest(androidTester.GetAvd(), initSnapName, androidTester.GetProjectId())
		if err != nil {
			_, _ = writer.WriteString("load snapshot error!\n")
			_ = writer.Flush()
			return perrorx.NewErrorXTesting("load snapshot error", err)
		}
		// start
		ans, err := androidTester.Start(i)
		if err != nil {
			_, _ = writer.WriteString("start error!\n")
			_ = writer.Flush()
			return perrorx.NewErrorXTesting("start error", err)
		}
		_, _ = writer.WriteString(ans + "\n")
		_ = writer.Flush()
		// 获取覆盖语句
		err = getStmtLogAfterEachTest(originStackTrace, crashPass, androidTester.GetAvd(), androidTester.GetProjectId(),
			androidTester.GetTestGenName(), androidTester.GetTestPrefix())
		if err != nil {
			_, _ = writer.WriteString("get stmtlog error!\n")
			_ = writer.Flush()
			return perrorx.NewErrorXTesting("get stmtlog error", err)
		}
	}
	finishedTime := time.Now()
	totalTime := time.Since(startTime)
	// 结束
	logrus.Debugf("[testing] projectId=" + androidTester.GetProjectId() + " finished at " + finishedTime.String() +
		", totalTime=" + totalTime.String())
	return nil
}