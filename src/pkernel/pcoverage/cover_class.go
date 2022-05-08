package pcoverage

// 被覆盖的每个类
type CoverClass struct {
	// 类名, stmtlog里解析出来是带$的, 而这里处理过, 是不带$的
	MClassName string `json:"className"`
	// 被覆盖的代码行列表
	MLines []int `json:"lines"`
	// 与lines对应, 表示每行被哪些动作覆盖(二进制表示)
	MEventIds []string `json:"eventids"`
}

