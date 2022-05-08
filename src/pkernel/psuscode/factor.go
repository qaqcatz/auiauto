package psuscode

import (
	"auiauto/perrorx"
	"math"
	"sort"
	"strconv"
	"strings"
)

func DoSusFormula(factor string, susCodes SusCodes, st map[string]int) *perrorx.ErrorX {
	switch factor {
	case "Ochiai":
		doOchiai(susCodes, st)
	case "Tarantula":
		doTarantula(susCodes, st)
	case "Barinel":
		doBarinel(susCodes, st)
	case "DStar":
		doDStar(susCodes, st)
	case "Op2":
		doOp2(susCodes, st)
	default:
		return perrorx.NewErrorXReadSourceTreeAndAnalyze("unknown factor", nil)
	}
	return nil
}

// ochiai
func doOchiai(susCodes SusCodes, st map[string]int) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		a01 := float64(susCodes[i].MA01)
		//a00 := float64(susCodes[i].MA00)
		susCodes[i].MValue = int(1000000 * a11 / math.Sqrt((a11+a01)*(a11+a10)))
	}
	sortRank(susCodes, st)
}

// tarantula
func doTarantula(susCodes SusCodes, st map[string]int) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		a01 := float64(susCodes[i].MA01)
		a00 := float64(susCodes[i].MA00)

		x := a11 / (a11 + a01)
		y := a10 / (a10 + a00)
		susCodes[i].MValue = int(1000000 * x / (x+y))
	}
	sortRank(susCodes, st)
}

// barinel
func doBarinel(susCodes SusCodes, st map[string]int) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		//a01 := float64(susCodes[i].MA01)
		//a00 := float64(susCodes[i].MA00)

		susCodes[i].MValue = int(1000000 * (1.0 - a10 / (a10+a11)))
	}
	sortRank(susCodes, st)
}

// dstar
func doDStar(susCodes SusCodes, st map[string]int) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		a01 := float64(susCodes[i].MA01)
		//a00 := float64(susCodes[i].MA00)

		susCodes[i].MValue = int(1000000 * a11*a11 / (a10+a01))
	}
	sortRank(susCodes, st)
}

// op2
func doOp2(susCodes SusCodes, st map[string]int) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		//a01 := float64(susCodes[i].MA01)
		a00 := float64(susCodes[i].MA00)

		susCodes[i].MValue = int(1000000 * (a11 - a10/(a10+a00+1)))
	}
	sortRank(susCodes, st)
}

// 实验性功能开关, 表示是否将value相同语句的排名标记为他们的中位数
// 这个变量是只读的, 无需担心并发问题
var gMidFlag = false

func sortRank(susCodes SusCodes, st map[string]int) {
	// 插桩那边可能会手动捕获异常, 多出来一些MyCrashxxx的类, 这里过滤掉
	for i := 0; i < len(susCodes); i++ {
		sig := susCodes[i].MClassName+"@"+strconv.Itoa(susCodes[i].MLine)
		if strings.Contains(susCodes[i].MClassName, "MyCrash") {
			susCodes[i].MValue = -1
		}
		if rk, ok := st[sig]; ok {
			if rk <= 0 {
				continue
			}
			if rk <= 10 {
				susCodes[i].MValue += 1000000 / rk
			} else {
				susCodes[i].MValue += 100000
			}
		}
	}
	//// 测试用, 统计纯栈信息的效果
	//for i := 0; i < len(susCodes); i++ {
	//	sig := susCodes[i].MClassName+"@"+strconv.Itoa(susCodes[i].MLine)
	//	susCodes[i].MValue = -1
	//	if rk, ok := st[sig]; ok {
	//		if rk <= 0 {
	//			continue
	//		}
	//		susCodes[i].MValue = 1000000 / rk
	//	}
	//}
	sort.Sort(susCodes)
	if gMidFlag {
		for i := 0; i < len(susCodes); i++ {
			susCodes[i].MRank = i+1
		}
		for i := 1; i < len(susCodes); i++ {
			if susCodes[i].MValue != susCodes[i-1].MValue {
				continue
			}
			start := i-1
			end := i
			for susCodes[i].MValue == susCodes[i-1].MValue {
				end = i
				i++
				if i >= len(susCodes) {
					break
				}
			}
			for j := start; j <= end; j++ {
				susCodes[j].MRank = (start+1+end+1)/2
			}
		}
	} else {
		for i := 0; i < len(susCodes); i++ {
			susCodes[i].MRank = i+1
		}
	}
}
