package puitree

import (
	"auiauto/perrorx"
	"auiauto/pkernel/pahttp"
	"auiauto/pkernel/pevent"
	"encoding/xml"
)

// todo 解决ui tree有时获取过慢问题
// ui tree, 从antrance中获取xml(或文件中存储下来的xml), 通过ParseUITree解析成UITree
type UITree struct {
	// 这个name是设计遗留问题, 恒为空
	MRoot  xml.Name  `xml:"rt"`
	MNodes []*UINode `xml:"nd"`
}

// dfs遍历所有节点, 对每个节点执行lamda函数
func (uiTree *UITree) Foreach(f func(father *UINode, prefix string, node *UINode)) {
	for _, root := range uiTree.MNodes {
		uiTree.foreachDFS(nil, "", root, f)
	}
}

// 辅助Foreach进行dfs
func (uiTree *UITree)foreachDFS(father *UINode, prefix string, cur *UINode, f func(father *UINode, prefix string, cur *UINode)) {
	f(father, prefix, cur)
	for _, child := range cur.MNodes {
		uiTree.foreachDFS(cur, prefix+"|  ", child, f)
	}
}

// 通过Foreach计算每个节点的父节点, 计算ui tree后会自动执行
func (uiTree *UITree) CalFather() {
	uiTree.Foreach(func(father *UINode, prefix string, cur *UINode) {
		cur.MFather = father
	})
}

// ui tree中是否包含event的作用对象, 也就是说event是否能在当前ui tree上执行
// 根据object匹配, 如果有多个, 找最相似的那个
// 相似度定义:
// prefix长度不一致相似度为-1(这个其实已经涵盖在了object比较中, 因为object中包含深度信息)
// 否则定义为对应位index相等的个数, 比如两个前缀1 2 3 4, 1 0 3 3, 0号位(1)和2号位(3)相等, 相似度为2
func (uiTree *UITree) CanPerform(event *pevent.Event) (bool, *UINode) {
	// keyevent, swipe, wait, rotate这种没有object的全局动作默认可以执行, 或许是个多余的判断, 不过还是留着吧
	if event.IsGlobal() {
		return true, nil
	}
	var ans *UINode = nil
	ansSim := -1
	uiTree.Foreach(func (father *UINode, prefix string, node *UINode) {
		if node.EventObjectCode() == event.MObject {
			nodeSim := event.CalSim(node.EventObjectPrefix())
			if ans == nil || ansSim < nodeSim {
				ans = node
				ansSim = nodeSim
			}
		}
	})
	return ans != nil, ans
}

// 将xml数据解析成UITree
func parseUITree(xmlData []byte) (*UITree, *perrorx.ErrorX) {
	uiTree := UITree{}
	err := xml.Unmarshal(xmlData, &uiTree)
	if err != nil {
		return nil, perrorx.NewErrorXUnmarshal(err.Error(), nil)
	}
	uiTree.CalFather()
	return &uiTree, nil
}

// 从antrance获取将xml并解析成UITree
func GetAndParseUITree(avd string) (*UITree, *perrorx.ErrorX) {
	xmlData, err := pahttp.GetUITree(avd)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	uiTree, err := parseUITree([]byte(xmlData))
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	return uiTree, nil
}