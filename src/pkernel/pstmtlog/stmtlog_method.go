package pstmtlog

import (
	"auiauto/perrorx"
	"strconv"
	"strings"
)

// stmtlog中的各个函数
type StmtLogMethod struct {
	// 函数签名
	MethodSig     string
	// 由于设计问题, 函数入口插桩比较特殊, 只记录了methodsig, 没有记录stmtsig,
	// 因此不能将其存入Stmts, 要存储在StmtLogMethod中
	// 这里的MethodEventId相当于记录函数入口被哪些动作执行.
	MethodEventId string
	// 未解码的stmt
	StmtStrs      []string
	// 解码后的stmt
	Stmts []*StmtLogStmt
}

// 将StmtStrs解码成Stmts
func (method *StmtLogMethod) Parse() *perrorx.ErrorX {
	method.Stmts = make([]*StmtLogStmt, len(method.StmtStrs))
	for i, stmtStr := range method.StmtStrs {
		// 0:jid, 1:sid, 2:type, 3:value
		stmtStrSP := strings.Split(stmtStr, "@")

		var stmt StmtLogStmt
		var err error
		// 0.jid
		stmt.MJid, err = strconv.Atoi(stmtStrSP[0])
		if err != nil {
			return perrorx.NewErrorXAtoI(stmtStrSP[0], nil)
		}
		// 1.sid
		stmt.MSid, err = strconv.Atoi(stmtStrSP[1])
		if err != nil {
			return perrorx.NewErrorXAtoI(stmtStrSP[1], nil)
		}
		// 2.type
		stmt.MType = stmtStrSP[2]
		// 3.value
		stmt.MValue, err = strconv.Atoi(stmtStrSP[3])
		if err != nil {
			return perrorx.NewErrorXAtoI(stmtStrSP[3], nil)
		}
		// 4.event id
		stmt.MEventId = stmtStrSP[4]

		method.Stmts[i] = &stmt
	}
	return nil
}
