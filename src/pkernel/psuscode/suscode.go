package psuscode

import "auiauto/pkernel/psrctree"

// 可疑语句
type SusCode struct {
	// 可疑语句对应的srcnode
	MOriginNode *psrctree.SourceNode `json:"-"`
	// 可疑语句在srcnode中的下标
	MIdx int `json:"-"`
	// short class name, 对应src node的full name
	MClassShortName string `json:"classShortName"`
	// dot class name
	MClassName string `json:"className"`
	// 行号
	MLine int `json:"line"`
	// i: testcase j:block
	// counter xij ei
	// a11     1   1   -> failed
	// a10     1   0   -> passed
	// a01     0   1   -> totalFailed-failed
	// a00     0   0   -> totalPassed-passed
	// 简单来说:
	// a11 当前语句被多少错误用例覆盖过, failed(p)
	// a01 错误用例总数-a11, totalfailed = a11+a01
	// a10 当前语句被多少正确用例覆盖过, passed(p)
	// a00 正确用例总数-a10, totalpassed = a10+a00
	MA11  int `json:"a11"`
	MA10  int `json:"a10"`
	MA01  int `json:"a01"`
	MA00  int `json:"a00"`
	// 可疑度*1000000, 转换为整数
	MValue int `json:"value"`
	// 可疑语句排名
	MRank int `json:"rank"`
}
