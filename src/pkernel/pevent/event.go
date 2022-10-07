package pevent

// 单个动作
type Event struct {
	// 唯一id, 其实是在events中的次序
	MId int `json:"id"`
	// 类型, 如click, edit, scroll, keyevent ...
	MType string `json:"type"`
	// 一些type需要value, 如scroll, value 0, 1, 2, 3分别表示不同的滑动方向
	// 特别地, 如果value以@prewait{...}开头的话需要解析{}中的等待时间(ms), 表示做这个动作前需要等待多少毫秒,
	// 没有这个标识的话默认等待750ms
	MValue string `json:"value"`
	// 非全局event是有作用对象的, 比如click, 需要作用到明确的ui控件
	// 格式: 包@类@res id@操作码@深度
	// 操作码:
	// 支持|运算, 从低位到高位:
	// 0位: clickable
	// 1位: longClickable
	// 2位: editable
	// 3位: scrollable
	// 4位: checkable
	MObject string `json:"object"`
	// object的index前缀
	// ui元素匹配规则:
	// object匹配, 如果有多个, 找最相似的那个
	// 相似度定义:
	// prefix长度不一致相似度为-1(这个其实已经涵盖在了object比较中, 因为object中包含深度信息)
	// 否则定义为对应位index相等的个数, 比如两个前缀1 2 3 4, 1 0 3 3, 0号位(1)和2号位(3)相等, 相似度为2
	MPrefix []int `json:"prefix"`
	// 描述信息, 方便人阅读
	MDesc string `json:"desc"`
}

// 是否为全局动作, 全局动作是指不需要作用对象的动作, 像返回, home等keyevent是不需要作用对象的, 而click, edit这种需要作用在一个控件上
// 注意后续添加新动作时一定要考虑是否需要更新IsGlobal函数
func (event *Event) IsGlobal() bool {
	return event.MType == "keyevent" || event.MType == "swipe" || event.MType == "dswipe" ||
		event.MType == "wait" || event.MType == "rotate"
}

// 计算两个prefix相似度, prefix长度不一致相似度为-1(这个其实已经涵盖在了object比较中, 因为object中包含深度信息)
// 否则定义为对应位index相等的个数, 比如两个前缀1 2 3 4, 1 0 3 3, 0号位(1)和2号位(3)相等, 相似度为2
func (event *Event) CalSim(prefix []int) int {
	if len(event.MPrefix) != len(prefix) {
		return -1
	} else {
		sum := 0
		for i := 0; i < len(prefix); i++ {
			if event.MPrefix[i] == prefix[i] {
				sum += 1
			}
		}
		return sum
	}
}