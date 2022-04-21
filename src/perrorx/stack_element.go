package perrorx

import "strconv"

// 栈信息中的各个条目
type StackElement struct {
	MFilePath   string // 语句所处的文件
	MMethodName string // 语句所处的函数
	MLine       int    // 语句所处的行
}

func (stackElement *StackElement) ToString() string {
	ans := ""
	ans += stackElement.MFilePath
	ans += "(" + stackElement.MMethodName + ":" + strconv.Itoa(stackElement.MLine) + ")"
	return ans
}
