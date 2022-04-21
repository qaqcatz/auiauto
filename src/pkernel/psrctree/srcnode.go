package psrctree

// 源码树的节点
type SourceNode struct {
	// 包/源文件名
	MName string `json:"name"`
	// 全名
	MFullName string `json:"fullName"`
	// 子树涵盖的源码总行数
	MTotalNum int `json:"totalNum"`
	// 子树涵盖的被覆盖的源码总行数
	MCoverNum int `json:"coverNum"`
	// 子节点
	MChildren []*SourceNode `json:"children"`
	// 子节点map, 方便根据name进行查找
	MChildrenMap map[string]*SourceNode `json:"-"`
	// 若当前节点是源文件, 则用Codes存储其代码行
	MCodes []string `json:"-"`
	// 记录每行代码被哪些动作执行过, 这里就不是二进制表示了, 需要将二进制解析为实际动作
	MEventIds [][]int `json:"-"`
	// MCodes 中的每行代码被多少正确用例执行过
	MPassed []int `json:"-"`
	// MCodes 中的每行代码被多少错误用例执行过
	MFailed []int `json:"-"`
	// MCodes 中每行代码的可疑度排名
	MRanks []int `json:"-"`
}
