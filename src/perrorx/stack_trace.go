package perrorx

import (
	"runtime"
)

// 栈信息
type StackTrace struct {
	MStackElements []*StackElement // 栈元素列表, 越靠前越接近异常的抛出位置
}

func NewStackTrace() *StackTrace {
	return &StackTrace{
		MStackElements: make([]*StackElement, 0),
	}
}

// 将当前执行的语句加入栈中, 注意runtime.Caller()返回为true时才会添加
// 参数skip表示打印哪一层的栈, 0表示打印当前语句, 1表示打印上一层栈的语句, 由MErrorX决定(New时为3, Trans时为2)
func (stackTrace *StackTrace) AddStackTrace(skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		methodName := runtime.FuncForPC(pc).Name()
		stackTrace.MStackElements = append(stackTrace.MStackElements, &StackElement{
			MFilePath:   file,
			MMethodName: methodName,
			MLine:       line,
		})
	}
}

func (stackTrace *StackTrace) ToString() string {
	ans := ""
	for _, stackElement := range stackTrace.MStackElements {
		ans += stackElement.ToString() + "\n"
	}
	return ans
}
