package ptesting

import (
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"auiauto/pkernel/pevent"
	"strconv"
	"time"
)

// monkey测试
type MonkeyTester struct {
	avd string
	cost int
	projectId string
	crashCase string
	testPrefix string
	// param暂时没用
	param string
}

func NewMonkeyTester(avd string, cost int, projectId string, crashCase string, testPrefix string, param string) *MonkeyTester {
	return &MonkeyTester {
		avd: avd,
		cost: cost,
		projectId: projectId,
		crashCase: crashCase,
		testPrefix: testPrefix,
		param: param,
	}
}

// 获取本次测试对应的模拟器
func (monkey *MonkeyTester) GetAvd() string {
	return monkey.avd
}

// 获取要生成的测试用例的总数, +表示数量, -表示时间
func (monkey *MonkeyTester) GetCost() int {
	return monkey.cost
}

// 获取本次测试的项目id
func (monkey *MonkeyTester) GetProjectId() string {
	return monkey.projectId
}

// 获取本次测试的基础crash case
func (monkey *MonkeyTester) GetCrashCase() string {
	return monkey.crashCase
}

// 获取测试工具名字, 如monkey
func (monkey *MonkeyTester) GetTestGenName() string {
	return "monkey"
}

// 获取本次测试的前缀标识, 如rd
func (monkey *MonkeyTester) GetTestPrefix() string {
	return monkey.testPrefix
}

// 单次测试
// https://blog.csdn.net/weixin_33554506/article/details/117798443
// adb shell monkey -p org.mozilla.rocket.debug.hzy -s 1 -v --throttle 500 100
// 最后一个指定monkey event数量的参数设置为10倍的crashcase evnet数
// 返回执行的命令, 有错误流的话打印错误流
func (monkey *MonkeyTester) Start(id int) (string, *perrorx.ErrorX) {
	applicationId, _, err := padb.GetAppConfig(monkey.projectId)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	events, err := pevent.ReadEventsStd(monkey.projectId, monkey.crashCase)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	monkeyEventNum := 10*len(events.MEvents)
	monkeyCmd := "monkey -p " + applicationId + " -s " + strconv.Itoa(id) + " -v --throttle 500 " +
		strconv.Itoa(monkeyEventNum)
	ans := monkeyCmd + "\n"
	// monkey报错时可能是app崩溃, 而不是monkey本身的原因. 因此我们不抛出异常, 仅作为结果返回.

	cha := make(chan *perrorx.ErrorX, 1)
	go func() {
		_, err = padb.AdbShell(monkey.avd, monkeyCmd)
		cha <- err
	} ()

	select {
	case err := <-cha:
		if err != nil {
			ans += err.Error() + "\n"
		}
	case <-time.After(time.Duration((len(events.MEvents)+2) * 3000) * time.Millisecond):
		ans += perrorx.NewErrorXTimeout(nil).Error()+"\n"
	}

	return ans, nil
}
