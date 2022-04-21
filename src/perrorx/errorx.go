package perrorx

type ErrorX struct {
	MType        string      // 异常类型
	MDescription string      // 异常描述
	MStackTrace  *StackTrace // 栈信息
	MCauseBy     *ErrorX     // 相当于java异常中的cause by
}

func (errorX *ErrorX) Error() string {
	ans := ""
	currentError := errorX
	for {
		ans += "[" + errorX.MType + "] " + errorX.MDescription + "\n"
		ans += currentError.MStackTrace.ToString()
		if currentError.MCauseBy == nil {
			break
		}
		currentError = currentError.MCauseBy
	}
	return ans
}

// 新建一个ErrorX, 若其依赖于另一个ErrorX(设为X), 则将X传给参数cause, 否则将cause置nil即可
func NewErrorX(mType string, mDescription string, cause *ErrorX) *ErrorX {
	ans := &ErrorX{
		MType:        mType,
		MDescription: mDescription,
		MStackTrace:  NewStackTrace(),
		MCauseBy:     cause,
	}
	// 注意将当前语句加到栈信息中, 由于NewErrorX还会被各种类型的错误创建函数封装, 这里skip设为3
	ans.MStackTrace.AddStackTrace(3)
	return ans
}

// 调用TransErrorX可以将当前语句记录在errorX的栈中, 并将errorX重新返回, 用于异常向上传递
func TransErrorX(errorX *ErrorX) *ErrorX {
	errorX.MStackTrace.AddStackTrace(2)
	return errorX
}