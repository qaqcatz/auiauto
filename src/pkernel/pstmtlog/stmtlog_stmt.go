package pstmtlog

// stmtlog中解析出的语句
type StmtLogStmt struct {
	// 字节码id(函数中的行号)
	MJid     int
	// 源码id
	MSid     int
	// type目前只考虑了branch(br)
	MType    string
	// 对于br, value 0表示false/default分支, >=1表示true/各个case分支)
	MValue   int
	// 被哪些动作覆盖, 二进制表示
	MEventId string
}
