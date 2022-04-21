package pperform

import (
	"auiauto/perrorx"
	"auiauto/pkernel/padb"
	"auiauto/pkernel/pahttp"
	"auiauto/pkernel/pevent"
	"auiauto/pkernel/puitree"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// 将left@top@right@bottom解析成left top right bottom
func parseBounds(bds string) (int, int, int, int, *perrorx.ErrorX) {
	sp := strings.Split(bds, "@")
	if len(sp) != 4 {
		return -1, -1, -1, -1, perrorx.NewErrorXSplitN(len(sp), 4, nil)
	}
	ans := make([]int, 4)
	for i := 0; i < 4; i++ {
		temp, err :=  strconv.Atoi(sp[i])
		if err != nil {
			return -1, -1, -1, -1, perrorx.NewErrorXAtoI(sp[i], nil)
		}
		ans[i] = temp
	}
	return ans[0], ans[1], ans[2], ans[3], nil
}

// 循环等待包含动作作用对象的ui tree(请求间隔1s), 最长等待10s, 超时报错
func waitAvaUITree(avd string, event *pevent.Event) (*puitree.UINode, *puitree.UITree, *perrorx.ErrorX) {
	t1 := time.Now()
	timeout := 10.0
	for {
		uiTree, err := puitree.GetAndParseUITree(avd)
		if err != nil {
			return nil, nil, perrorx.TransErrorX(err)
		}
		canPerform, uiNode := uiTree.CanPerform(event)
		if canPerform {
			return uiNode, uiTree, nil
		}
		t2 := time.Now()
		d := t2.Sub(t1)
		if d.Seconds() > timeout {
			break
		}
		time.Sleep(time.Millisecond*1000)
		logrus.Debugf("waitAvaUITree: wait 1s")
	}
	return nil, nil, perrorx.NewErrorXWaitAvaUITree("timeout, can not perform this event", nil)
}

// 执行单个动作, 前端点击操作以及PostPerforms的基础, 操作成功时返回动作执行前的ui tree(可能为null, 目前这个返回值用处不大了)
// 1. 动作执行前等待一段时间
// 2. 给antrance设置event id, 这是为了计算覆盖率时能显示每行代码被哪些动作执行过.
//		注意: isInit参数标识当前执行的是初始动作序列中的动作, 此时发送的event id为0
// 3. 根据event是否为global类型, 判断是否需要获取ui tree.
//		若是global类型, 则不许要获取ui tree, 直接执行即可;
//		之后根据event type执行相应的动作.
func PostPerform(avd string, eventOrigin *pevent.Event, isInit bool) (*puitree.UITree, *perrorx.ErrorX) {
	logrus.Debugf("PostPerform: ==================================================")
	logrus.Debugf("PostPerform: " + avd + " perform " + strconv.Itoa(eventOrigin.MId) + " " +
		eventOrigin.MType + " " + eventOrigin.MValue)
	defer func() {
		logrus.Debugf("PostPerform: ==================================================")
	} ()

	// 1. 动作执行前等待一段时间
	// 如果event的value以@prewait{...}开头的话需要解析{}中的等待时间(ms), 表示做这个动作前需要等待多少毫秒,
	// 没有这个标识的话默认等待750ms
	// 因为后面很多地方要用到去除@prewait的value, 我们将eventOrigin拷贝一份, 处理一下value, 还有这里prefix共用即可
	waitBefore := 750
	event := &pevent.Event{eventOrigin.MId, eventOrigin.MType, eventOrigin.MValue,
		eventOrigin.MObject, eventOrigin.MPrefix, eventOrigin.MDesc}
	if strings.HasPrefix(event.MValue, "@prewait") {
		sp := strings.SplitN(event.MValue, "}", 2)
		if len(sp) != 2 {
			return nil, perrorx.NewErrorXSplitN(len(sp), 2, nil)
		}
		event.MValue = sp[1]
		sp_ := strings.Split(sp[0], "{")
		if len(sp_) != 2 {
			return nil, perrorx.NewErrorXSplitN(len(sp_), 2, nil)
		}
		var err error = nil
		waitBefore, err = strconv.Atoi(sp_[1])
		if err != nil {
			return nil, perrorx.NewErrorXAtoI(sp_[1], nil)
		}
	}
	logrus.Debugf("PostPerform: waitBefore: " + strconv.Itoa(waitBefore) + " ms")
	time.Sleep(time.Millisecond*time.Duration(waitBefore))

	// 2. 给antrance设置event id
	eventId := event.MId
	if isInit {
		eventId = 0
	}
	_, err := pahttp.AntranceRequest("GET", avd, "seteventid?id="+strconv.Itoa(eventId), nil)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}

	// 3. 根据event是否为global类型, 判断是否需要获取ui tree. 若是global类型, 则不许要获取ui tree, 直接执行即可;
	// 之后根据event type执行相应的动作.
	// 这个err只是为了方便做一些统一的错误判断
	err = nil
	// 动作作用的ui tree, 目前没什么用了
	var uiTree *puitree.UITree = nil
	if !event.IsGlobal() {
		var uiNode *puitree.UINode = nil
		// 判断ui tree中是否包含当前动作的作用对象, 超时失败, 成功返回作用对象(uiNode)以及当前uiTree
		uiNode, uiTree, err = waitAvaUITree(avd, event)
		if err != nil {
			return nil, perrorx.TransErrorX(err)
		}
		// 解析出uiNode的边界
		left, top, right, bottom := -1, -1, -1, -1
		if uiNode != nil {
			left, top, right, bottom, err = parseBounds(uiNode.MBds)
			if err != nil {
				return nil, perrorx.TransErrorX(err)
			}
		}
		// 计算中心坐标
		x, y := (left+right)/2, (top+bottom)/2

		switch event.MType {
		case "click":
			err = padb.AdbInputTap(avd, x, y)
		case "dclick":
			err = padb.AdbInputTap(avd, x, y)
			if err == nil {
				// 双击间隔
				ms, err := strconv.Atoi(event.MValue)
				if err != nil {
					return nil, perrorx.NewErrorXAtoI(event.MValue, nil)
				}
				time.Sleep(time.Duration(ms) * time.Millisecond)
				err = padb.AdbInputTap(avd, x, y)
			}
		case "longclick":
			// 默认长按1.5s
			err = padb.AdbInputSwipe(avd, x, y, x, y, 1500)
		case "edit":
			// edit比较特殊, adb命令无法实现一些特殊字符的输入, 因此这里借助antrance内的accessibility service实现
			jsonData, err := json.Marshal(event)
			if err == nil {
				_, err = pahttp.AntranceRequest("POST", avd, "perform", jsonData)
			}
		case "editx":
			// editx是为了告诉auiauto这个输入是固定的, 无需变异
			jsonData, err := json.Marshal(event)
			if err == nil {
				_, err = pahttp.AntranceRequest("POST", avd, "perform", jsonData)
			}
		case "scroll":
			// 默认滑动150ms
			switch event.MValue {
			case "0":
				// 左
				err = padb.AdbInputSwipe(avd, x, y, x/2, y, 150)
			case "1":
				// 上
				err = padb.AdbInputSwipe(avd, x, y, x, y/2, 150)
			case "2":
				// 右
				err = padb.AdbInputSwipe(avd, x, y, x+x/2, y, 150)
			case "3":
				// 下
				err = padb.AdbInputSwipe(avd, x, y, x, y+y/2, 150)
			default:
				err = perrorx.NewErrorXPerform("unknown direction", nil)
			}
		case "check":
			switch event.MValue {
			// 其实有快照还原状态的话完全可以用click取代check, 不推荐使用check
			case "0":
				// 取消
				if (uiNode.MSta & 1) != 0 {
					err = padb.AdbInputTap(avd, x, y)
				}
			case "1":
				// 选中
				if (uiNode.MSta & 1) == 0 {
					err = padb.AdbInputTap(avd, x, y)
				}
			default:
				err = perrorx.NewErrorXPerform("unknown check value " + event.MValue, nil)
			}
		default:
			err = perrorx.NewErrorXPerform("unknown type " + event.MType, nil)
		}
	} else {
		switch event.MType {
		case "keyevent":
			err = padb.AdbInputKeyEvent(avd, event.MValue)
		case "swipe":
			// swipe用于任意坐标滑动, 时间为ms, 也可以将起始坐标设为同一个实现任意位置点击/长按
			sp := strings.Split(event.MValue, " ")
			// left top right bottom ms
			if len(sp) != 5 {
				return nil, perrorx.NewErrorXSplitN(len(sp), 5, nil)
			}
			spInt := make([]int, 5)
			for i := 0; i < 5; i++ {
				temp, err :=  strconv.Atoi(sp[i])
				if err != nil {
					return nil, perrorx.NewErrorXAtoI(sp[i], nil)
				}
				spInt[i] = temp
			}
			err = padb.AdbInputSwipe(avd, spInt[0], spInt[1], spInt[2], spInt[3], spInt[4])
		case "dswipe":
			// 连续swipe两次, 间隔wait ms
			sp := strings.Split(event.MValue, " ")
			// left top right bottom ms waitms
			if len(sp) != 6 {
				return nil, perrorx.NewErrorXSplitN(len(sp), 6, nil)
			}
			spInt := make([]int, 6)
			for i := 0; i < 6; i++ {
				temp, err :=  strconv.Atoi(sp[i])
				if err != nil {
					return nil, perrorx.NewErrorXAtoI(sp[i], nil)
				}
				spInt[i] = temp
			}
			err = padb.AdbInputSwipe(avd, spInt[0], spInt[1], spInt[2], spInt[3], spInt[4])
			if err == nil {
				time.Sleep(time.Duration(spInt[5]) * time.Millisecond)
				err = padb.AdbInputSwipe(avd, spInt[0], spInt[1], spInt[2], spInt[3], spInt[4])
			}
		case "wait":
			// 如果最后一步动作结束后还需要等待一段时间的话可以用wait操作, 否则建议使用每个动作的前等待.
			ms, err := strconv.Atoi(event.MValue)
			if err != nil {
				return nil, perrorx.NewErrorXAtoI(event.MValue, nil)
			}
			time.Sleep(time.Millisecond*time.Duration(ms))
		case "rotate":
			// 取消自动旋转, 相当于我们经常设置的固定屏幕不旋转
			err := padb.AdbAccRotation(avd, "0")
			if err != nil {
				return nil, perrorx.TransErrorX(err)
			}
			// 设置旋转方向, 0表示竖着, 1表示横着
			err = padb.AdbUserRotation(avd, event.MValue)
			if err != nil {
				return nil, perrorx.TransErrorX(err)
			}
		default:
			err = perrorx.NewErrorXPerform("unknown type " + event.MType, nil)
		}
	}

	// 统一的错误判断
	if err != nil {
		return nil, err
	}

	return uiTree, nil
}


