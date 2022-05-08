package pstatistic

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"bufio"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

// 组合分析结果
type StaCombine struct {
	// 组合总数2^(n-1)
	MN int `json:"n"`
	// 正确用例数
	MPN int `json:"pn"`
	// 错误用例数
	MFN int `json:"fn"`
	// 所有组合中最好的map
	MBestMap string `json:"bestMap"`
	// bestMap对应的正确用例数量
	MBestPassNum int `json:"bestPassNum"`
	// bestMap对应的错误用例数量
	MBestCrashNum int `json:"bestCrashNum"`
	// bestMap对应的rank
	MBestRank []int `json:"bestRank"`
	// 用例数(容量为n, 下标+1对应用例数):[每种组合的map]
	MNumMaps [][]string `json:"numMaps"`
}

// 统计组合分析结果
func StatisticCombineStd(projectId string, analyzeFile string, factor string, casePrefix string) (*StaCombine, *perrorx.ErrorX) {
	testCasePath := pdba.DBURLProjectIdTestcases(projectId)
	analyzePath := path.Join(testCasePath, "combineanalyze_"+casePrefix+"_"+analyzeFile+"_"+factor+".txt")
	file, err := os.Open(analyzePath)
	if err != nil {
		return nil, perrorx.NewErrorXOpen(analyzePath, nil)
	}
	scanner := bufio.NewScanner(file)
	// 第一行两个部分, 表示用例总数n(组合总数2^(n-1)), 聚焦语句数量s
	scanner.Scan()
	line := strings.TrimSpace(scanner.Text())
	sp := strings.Split(line, " ")
	if len(sp) != 2 {
		return nil, perrorx.NewErrorXSplitN(len(sp), 2, nil)
	}
	n, err1 := strconv.Atoi(sp[0])
	s, err2 := strconv.Atoi(sp[1])
	if err1 != nil || err2 != nil {
		return nil, perrorx.NewErrorXAtoI(sp[0]+"|"+sp[1], nil)
	}
	// 统计最佳map及其对应的正确/错误用例数以及rank
	bestMap := 0.0
	bestPassNum := -1
	bestCrashNum := -1
	bestRank := make([]int, 0)
	// 用例数(容量为n, 下标+1对应用例数):[每种组合的map]
	numMaps := make([][]string, n)
	for i := 0; i < n; i++ {
		numMaps[i] = make([]string, 0)
	}
	// 第二行空格打印所用的用例
	scanner.Scan()
	line = strings.TrimSpace(scanner.Text())
	cases := strings.Split(line, " ")
	if len(cases) != n {
		return nil, perrorx.NewErrorXSplitN(len(cases), n, nil)
	}
	pn := 0
	fn := 0
	for i := 0; i < n; i++ {
		if strings.Contains(cases[i], "pass") {
			pn += 1
		}
		if strings.Contains(cases[i], "crash") {
			fn += 1
		}
	}
	// 接下来2^(n-1)个组合, 每个组合:
	for i := 0; i < (1<<(n-1)); i++ {
		// 	第一行一个01字符串(长度为n,原始错误用例一定为1), 二进制形式表示所选的用例, le存储
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		// 统计1的个数
		curNum := 0
		// 根据1对应的case是否包含pass|crash统计正确用例和错误用例的个数, 注释我们用的le存储, 可以直接顺序对应
		curPassNum := 0
		curCrashNum := 0
		if len(line) != n {
			return nil, perrorx.NewErrorXSplitN(len(line), n, nil)
		}
		for j := 0; j < len(line); j++ {
			if line[j] == '1' {
				curNum++
				if strings.Contains(cases[j], "pass") {
					curPassNum++
				} else if strings.Contains(cases[j], "crash") {
					curCrashNum++
				}
			}
		}
		// 接下来s个数, 代表聚焦语句文件中每个聚焦语句的排名
		scanner.Scan()
		line = strings.TrimSpace(scanner.Text())
		sp := strings.Split(line, " ")
		if len(sp) != s {
			return nil, perrorx.NewErrorXSplitN(len(sp), s, nil)
		}
		rank := make([]int, s)
		for j := 0; j < s; j++ {
			rank[j], err = strconv.Atoi(sp[j])
			if err != nil {
				return nil, perrorx.NewErrorXAtoI(sp[j], nil)
			}
		}

		rank_ := make([]int, 0)
		for _, r := range rank {
			if r <= 0 {
				continue
			}
			rank_ = append(rank_, r)
		}
		sort.Ints(rank_)
		curMap := 0.0
		for j := 0; j < len(rank_); j++ {
			// 这里注意用k+1
			curMap += float64(j+1) / float64(rank_[j])
		}
		curMap /= float64(s)
		// 更新bestMap
		if curMap > bestMap {
			bestMap = curMap
			bestPassNum = curPassNum
			bestCrashNum = curCrashNum
			bestRank = rank
		}
		// 更新numMaps
		numMaps[curNum-1] = append(numMaps[curNum-1], strconv.FormatFloat(curMap, 'f', 5, 64))
	}
	ans := &StaCombine {
		MN: n,
		MPN: pn,
		MFN: fn,
		MBestMap: strconv.FormatFloat(bestMap, 'f', 5, 64),
		MBestPassNum: bestPassNum,
		MBestCrashNum: bestCrashNum,
		MBestRank: bestRank,
		MNumMaps: numMaps,
	}
	return ans, nil
}
