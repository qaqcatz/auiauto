package pcfg

import (
	"auiauto/perrorx"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// 控制流图, 通过解析soot产出的cfg json而来
// soot传来的cfg json格式:
// { "methods"(类文件下的函数):[
//     {
//       "methodSig"(函数签名):"de.rampro.activitydiary.ui.generic.DetailViewHolders@void onClick(android.view.View)",
//       "stmts"(按soot结构信息记录的语句):[
//         "jid@sid@type@targets@fallThrough"(jid:语句在函数中的字节码id, sid:语句在文件中的源码行,
//           type分为branch(br),goto(gt),normal(n):
//           1.br: targets是若干用'#'分隔的数字, 表示目标语句jid, 编号0对应false/default分支, >=1对应true/各个case分支)
//           2.gt: targets表示跳转语句jid
//           3.n: targets为空
//           fallThrough(0/1)对应soot中的stmt.fallsThrough(), 表示当前语句和jid+1的语句执行上是否相连)
//       ]
//     }
//   ]
// }
type CFG struct {
	// 控制流图中的method集合
	MCFGMethods   []*CFGMethod `json:"methods"`
	// 将Methods中的每一项转换为methodSig:Method保存在map中
	MCFGMethodMap map[string]*CFGMethod `json:"-"`
}

// 将Methods中的每一项转换为methodSig:Method保存在map中
func (cfg *CFG) Parse() *perrorx.ErrorX {
	cfg.MCFGMethodMap = make(map[string]*CFGMethod)
	for _, method := range cfg.MCFGMethods {
		cfg.MCFGMethodMap[method.MMethodSig] = method
		err := method.Parse()
		if err != nil {
			return perrorx.TransErrorX(err)
		}
	}
	return nil
}

func (cfg *CFG) Debug() {
	for methodSig, method := range cfg.MCFGMethodMap {
		fmt.Println("==================================================")
		fmt.Println(methodSig)
		fmt.Println("==================================================")
		for _, stmt := range method.MStmts {
			fmt.Println(stmt.MJid, stmt.MSid, stmt.MType, stmt.MTargets, stmt.MFallThrough)
		}
	}
}

// 在MCFGMethodMap中查询给定method是否存在
func (cfg *CFG) Search(methodSig string) *CFGMethod {
	if method, ok := cfg.MCFGMethodMap[methodSig]; ok {
		return method
	} else {
		return nil
	}
}

// 从cfgPath中加载CFG, 转化为CFGMethods的形式返回
func ReadCFG(cfgPath string) (*CFG, *perrorx.ErrorX) {
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, perrorx.NewErrorXReadFile(cfgPath, err.Error(), nil)
	}
	var cfg CFG
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, perrorx.NewErrorXUnmarshal(err.Error(), nil)
	}
	err_ := cfg.Parse()
	if err_ != nil {
		return nil, perrorx.TransErrorX(err_)
	}
	return &cfg, nil
}