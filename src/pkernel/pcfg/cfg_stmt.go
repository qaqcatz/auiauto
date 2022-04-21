package pcfg

// 控制流图中的每条语句
type CFGStmt struct {
	// 字节码id(函数中的行号)
	MJid         int
	// 源码id(行号)
	MSid         int
	// 语句类型
	// type分为branch(br),goto(gt),normal(n):
	// 1.br: targets是若干用'#'分隔的数字, 表示目标语句jid, 编号0对应false/default分支, >=1对应true/各个case分支)
	// 2.gt: targets表示跳转语句jid
	// 3.n: targets为空
	MType        string
	// 语句指向的目标, 可以参考type中的说明
	MTargets     []int
	// fallThrough(0/1)对应soot中的stmt.fallsThrough(), 表示当前语句和jid+1的语句执行上是否相连)
	MFallThrough bool
}