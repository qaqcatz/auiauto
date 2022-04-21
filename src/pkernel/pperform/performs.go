package pperform

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"auiauto/pkernel/pahttp"
	"auiauto/pkernel/pcoverage"
	"auiauto/pkernel/pevent"
	"auiauto/putils"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

// 若apkPath存在, 使用defer实现在PostPerforms正常/异常执行后关闭app
// 这里忽略了异常, 不用担心app没启动时调用kill报错的情况
func postPerformsClean(apkPath string, avd string, projectId string) {
	if !putils.FileExist(apkPath) {
		return
	}
	applicationId, _, err := padb.GetAppConfig(projectId)
	if err != nil {
		return
	}
	_, _  = padb.AdbKill(avd, applicationId)
}

// init app, 即加载init快照, 重新安装, 重新启动, 执行初始用例
func PerformInit(avd string, projectId string, events *pevent.Events) *perrorx.ErrorX {
	// load snapshot
	if events.MInitSnapshot != "" {
		_, err := padb.AdbEmuLoadSnapshot(avd, events.MInitSnapshot)
		if err != nil {
			return perrorx.TransErrorX(err)
		}
		// 默认等待3s, 给snapshot准备时间
		logrus.Debugf("PostPerforms: wait 3s after loading snapshot")
		time.Sleep(time.Millisecond*3000)
	}
	// reInstall
	if events.MReInstall {
		_, err := padb.AdbOnlyInstall(avd, projectId)
		if err != nil {
			return perrorx.TransErrorX(err)
		}
	}
	// reStart
	if events.MReStart {
		_, err := padb.AdbOnlyStart(avd, projectId)
		if err != nil {
			return perrorx.TransErrorX(err)
		}
		// 默认等待3s, 给app准备时间
		logrus.Debugf("PostPerforms: wait 3s after app start(default)")
		time.Sleep(time.Millisecond*3000)
	}

	// antrance init
	err := pahttp.GetInit(avd)
	if err != nil {
		return perrorx.TransErrorX(err)
	}

	// load init case
	if events.MInitTestcase != "" {
		initCasePath := pdba.DBURLProjectIdTestcaseTestcase(projectId, events.MInitTestcase)
		if !putils.FileExist(initCasePath) {
			return perrorx.NewErrorXFileNotFound(initCasePath, nil)
		}
		initEvents, err := pevent.ReadEventsStd(projectId, events.MInitTestcase)
		if err != nil {
			return perrorx.TransErrorX(err)
		}
		for i := 0; i < len(initEvents.MEvents); i++ {
			_, err := PostPerform(avd, initEvents.MEvents[i], true)
			if err != nil {
				return perrorx.TransErrorX(err)
			}
		}
	}
	return nil
}

// 执行动作序列, 返回执行成功的动作数, 相当于告诉用户[0, return)区间的动作执行成功
// 1. 初始化
//		从projectId/testcases/caseName/testcase.json中读取动作序列
// 		判断projectId下是否有apk目录.
//		有的话说明这是一个特定app的复现脚本, 需要执行init app逻辑, 跳转到2;
// 		否则就是一个普通的复现脚本, 用户需要自己保证当前界面处于动作序列开始位置, 跳转到3.
// 2. init app: 根据testcase.json中的配置依次执行加载初始快照, 重新安装app, 重新启动app操作
// 		之后向antrance发送init请求, 开启动作复现. 复现开始后要先根据testcase.json中的配置执行
//		初始动作序列的加载(如登陆, 授权等操作).
//		注意: 加载初始快照后要等待3s, 防止我们在模拟器中安装的antrance服务来不及恢复.
//		app重新启动后也等待3s, 给app准备时间.
// 3. 执行动作序列, 注意一个动作执行失败时不要立即报错, 执行失败的情况下也要尝试执行后续的日志获取操作(这么设计其实是为变异服务的,
//		变异后的序列可能无法执行完, 但可能也会被接受, 并需要获取日志, 虽然现在来看我们用不到变异了)
//		记录好执行成功了多少个动作, 根据apkPath是否存在跳转到4(存在)或5(不存在)
// 4. 日志获取: 根据logIdSig计算stmtlog.json, 接着结合cfg计算cover.json
//		注意: 如果cover计算失败的话视为动作序列全部执行失败, 没什么特殊意义, 规范一点
//		app执行后等待3s, 防止一些代码没来得及覆盖, 以及崩溃没来得及发生
// 5. 返回执行成功的动作数
// 最后一个返回值表示app是否崩溃, 有error的情况默认为false
func PostPerforms(avd string, projectId string, caseName string) (int, *perrorx.ErrorX, bool) {
	// 1. 初始化
	// 清除cover.json
	_ = os.Remove(pdba.DBURLProjectIdTestcaseCover(projectId, caseName))
	// 先获取一下apkPath, 根据apkPath是否存在判断当前用例是特定app的测试还是普通的复现脚本
	apkPath := pdba.DBURLProjectIdAPKAPP(projectId)
	// 若apkPath存在, 使用defer实现在PostPerforms正常/异常执行后关闭app
	defer postPerformsClean(apkPath, avd, projectId)
	// 读取测试用例
	events, err := pevent.ReadEventsStd(projectId, caseName)
	if err != nil {
		return 0, perrorx.TransErrorX(err), false
	}

	// 2. 根据apkPath是否存在判断是否要执行init app操作
	if putils.FileExist(apkPath) {
		logrus.Debugf("init app >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		err := PerformInit(avd, projectId, events)
		if err != nil {
			return 0, perrorx.TransErrorX(err), false
		}
	}

	// 3. 执行动作序列
	logrus.Debugf("performs >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// 注意一个动作执行失败时不要立即报错, 执行失败的情况下也要尝试执行后续的日志获取操作(这么设计其实是为变异服务的,
	// 变异后的序列可能无法执行完, 但可能也会被接受, 并需要获取日志, 虽然现在来看我们用不到变异了)
	// 这里performErr是为了暂存中间可能产生的错误
	var performErr *perrorx.ErrorX = nil
	// i记录好执行成功了多少个动作
	i := 0
	for ; i < len(events.MEvents); i++ {
		_, err := PostPerform(avd, events.MEvents[i], false)
		if err != nil {
			performErr = perrorx.NewErrorXPerforms("perform event " + strconv.Itoa(i) + " error", err)
			break
		}
	}

	// 4.根据apkPath是否存在判断是否要执行日志获取操作
	// 根据logIdSig计算stmtlog.json, 接着结合cfg计算cover.json
	// 如果cover计算失败的话视为全部失败
	// crashFlag表示执行成功的情况下app是否崩溃
	crashFlag := false
	if putils.FileExist(apkPath) {
		// 等3s, 防止一些代码没来得及覆盖, 以及崩溃没来得及发生
		logrus.Debugf("PostPerforms: wait 3s before getting stmt log(default)")
		time.Sleep(time.Millisecond*3000)

		// 获取, 计算, 保存stmtlog和cover.json
		stmtLog, err := pcoverage.SaveStmtLogAndLineCoverageStd(avd, projectId, caseName)
		if err != nil {
			return 0, perrorx.TransErrorX(err), false
		}
		if stmtLog.MStatus == "false" {
			crashFlag = true
		}
	}

	// 5. 只是perform错误的话返回成功数量+perform失败原因
	if performErr != nil {
		return i, performErr, false
	}

	return i, nil, crashFlag
}
