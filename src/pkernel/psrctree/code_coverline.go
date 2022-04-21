package psrctree

// 用于可视化显示
type CodesAndCoverLines struct {
	// 每行代码
	Codes []string `json:"codes"`
	// 每行代码的类别
	// -1:未覆盖
	// 0: 覆盖,  lightgreen
	// 1: 可疑度rank1, DarkRed
	// 2: 可疑度rank10, Red
	// 3: 可疑度rank100, LightCoral
	// 4: 可疑度rank other, Yellow
	CodesType []int `json:"codesType"`
}