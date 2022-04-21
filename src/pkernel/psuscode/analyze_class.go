package psuscode

// 需要分析的文件中的每个类
type AnalyzeClass struct {
	MClassName string `json:"className"`
	MLines     []int  `json:"lines"`
}
