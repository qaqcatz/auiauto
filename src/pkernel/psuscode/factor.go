package psuscode

import (
	"math"
	"sort"
)

// ochiai
func DoOchiai(susCodes SusCodes) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		a01 := float64(susCodes[i].MA01)
		//a00 := float64(susCodes[i].MA00)
		susCodes[i].MValue = int(1000000 * a11 / math.Sqrt((a11+a01)*(a11+a10)))
	}
	sortRank(susCodes)
}

// tarantula
func DoTarantula(susCodes SusCodes) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		a01 := float64(susCodes[i].MA01)
		a00 := float64(susCodes[i].MA00)

		x := a11 / (a11 + a01)
		y := a10 / (a10 + a00)
		susCodes[i].MValue = int(1000000 * x / (x+y))
	}
	sortRank(susCodes)
}

// barinel
func DoBarinel(susCodes SusCodes) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		//a01 := float64(susCodes[i].MA01)
		//a00 := float64(susCodes[i].MA00)

		susCodes[i].MValue = int(1000000 * (1.0 - a10 / (a10+a11)))
	}
	sortRank(susCodes)
}

// dstar
func DoDStar(susCodes SusCodes) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		a01 := float64(susCodes[i].MA01)
		//a00 := float64(susCodes[i].MA00)

		susCodes[i].MValue = int(1000000 * a11*a11 / (a10+a01))
	}
	sortRank(susCodes)
}

// op2
func DoOp2(susCodes SusCodes) {
	for i := 0; i < len(susCodes); i++ {
		a11 := float64(susCodes[i].MA11)
		a10 := float64(susCodes[i].MA10)
		//a01 := float64(susCodes[i].MA01)
		a00 := float64(susCodes[i].MA00)

		susCodes[i].MValue = int(1000000 * (a11 - a10/(a10+a00+1)))
	}
	sortRank(susCodes)
}

// 实验性功能开关, 表示是否将value相同语句的排名标记为他们的中位数
// 这个变量是只读的, 无需担心并发问题
var gMidFlag = true

func sortRank(susCodes SusCodes) {
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
