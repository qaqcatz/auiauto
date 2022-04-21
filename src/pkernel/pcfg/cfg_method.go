package pcfg

import (
	"auiauto/perrorx"
	"strconv"
	"strings"
)

// 控制流图中的每个method
type CFGMethod struct {
	// 函数签名, 唯一标识
	MMethodSig string     `json:"methodSig"`
	// 从json中解析出来的原始语句编码, 还未解码
	MStmtStrs  []string   `json:"stmts"`
	// 由MStmtStrs解码而来
	MStmts     []*CFGStmt `json:"-"`
}

// 将jid@sid@type@targets@fallThrough解析为CFGStmt
func (method *CFGMethod) Parse() *perrorx.ErrorX {
	method.MStmts = make([]*CFGStmt, len(method.MStmtStrs))
	for _, stmtStr := range method.MStmtStrs {
		// 0:jid, 1:sid, 2:type, 3:targets, 4:fallThrough
		stmtStrSP := strings.Split(stmtStr, "@")

		var stmt CFGStmt
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
		// 3.targets
		stmt.MTargets = make([]int, 0)
		if stmtStrSP[3] != "" {
			targetStrs := strings.Split(stmtStrSP[3], "#")
			for _, targetStr := range targetStrs {
				var targetJid int
				targetJid, err = strconv.Atoi(targetStr)
				if err != nil {
					return perrorx.NewErrorXAtoI(targetStr, nil)
				}
				stmt.MTargets = append(stmt.MTargets, targetJid)
			}
		}
		// 4.fallThrough
		stmt.MFallThrough = stmtStrSP[4] == "1"

		method.MStmts[stmt.MJid] = &stmt
	}
	return nil
}
