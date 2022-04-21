package puitree

import (
	"strconv"
)

// ui tree节点
type UINode struct {
	// 节点在树中的深度
	MDp int `xml:"dp,attr"`
	// 节点在父亲children中的下标
	MIdx int `xml:"idx,attr"`
	// ui控件的边界, left@top@right@bottom, 表示左上, 右下坐标
	MBds string `xml:"bds,attr"`
	// 包
	MPkg string `xml:"pkg,attr"`
	// 类
	MCls string `xml:"cls,attr"`
	// resource id
	MRes string `xml:"res,attr"`
	// 描述
	MDsc string `xml:"dsc,attr"`
	// ui控件的文本信息, 比如TextView的text
	MTxt string `xml:"txt,attr"`
	// 操作码:
	// 支持|运算, 从低位到高位:
	// 1位: clickable
	// 2位: longClickable
	// 3位: editable
	// 4位: scrollable
	// 5位: checkable
	MOp int `xml:"op,attr"`
	// 状态码:
	// 1位: checked
	MSta    int       `xml:"sta,attr"`
	MNodes  []*UINode `xml:"nd"`
	MFather *UINode   `xml:"-"`
}

// 节点是否可点击
func (uiNode *UINode) IsClickable() bool {
	return (uiNode.MOp & 1) != 0
}

// 节点是否可长按
func (uiNode *UINode) IsLongClickable() bool {
	return (uiNode.MOp & 2) != 0
}

// 节点是否可编辑
func (uiNode *UINode) IsEditable() bool {
	return (uiNode.MOp & 4) != 0
}

// 节点是否可滑动
func (uiNode *UINode) IsScrollable() bool {
	return (uiNode.MOp & 8) != 0
}

// 计算uinode的object, 包@类@res id@可操作性@深度, 与event中的MObject做对比
func (uiNode *UINode) EventObjectCode() string {
	return uiNode.MPkg +"@"+uiNode.MCls +"@"+uiNode.MRes +"@"+
		strconv.Itoa(uiNode.MOp)+"@"+strconv.Itoa(uiNode.MDp)
}

// 计算uinode的index前缀
func (uiNode *UINode) EventObjectPrefix() []int {
	ans := make([]int, 0)
	cur := uiNode
	for cur != nil {
		ans = append(ans, cur.MIdx)
		cur = cur.MFather
	}
	// reverse
	for i, j := 0, len(ans)-1; i < j; i, j = i+1, j-1 {
		ans[i], ans[j] = ans[j], ans[i]
	}
	return ans
}
