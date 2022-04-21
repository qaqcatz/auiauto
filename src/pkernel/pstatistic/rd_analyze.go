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

// 统计随机采样结果
type StaRd struct {
	MPassNum  int `json:"passNum"`
	MCrashNum int `json:"crashNum"`
	MRandNum  int `json:"randNum"`
	// 各个正确/错误用例组合中的最优平均map
	MBestAvgMap string `json:"bestAvgMap"`
	// bestAvgMap对应的正确用例数
	MBestAvgPassNum int `json:"bestAvgPassNum"`
	// bestAvgMap对应的错误用例数
	MBestAvgCrashNum int `json:"bestAvgCrashNum"`
	// bestAvgMap对应中的中位数
	MBestAvgMidMap string `json:"bestAvgMidMap"`
	// bestAvgMidMap对应的rank
	MBestAvgMidRank []int `json:"bestAvgMidRank"`
	// bestAvgMap对应的最优值与rank
	MBestAvgMaxMap  string `json:"bestAvgMaxMap"`
	MBestAvgMaxRank []int  `json:"bestAvgMaxRank"`
	// bestAvgMap对应的最差值与rank
	MBestAvgMinMap  string `json:"bestAvgMinMap"`
	MBestAvgMinRank []int  `json:"bestAvgMinRank"`
	// bestAvgMap对应的1/4线与rank
	MBestAvg14Map  string `json:"bestAvg14Map"`
	MBestAvg14Rank []int  `json:"bestAvg14Rank"`
	// bestAvgMap对应的3/4线与rank
	MBestAvg34Map  string `json:"bestAvg34Map"`
	MBestAvg34Rank []int  `json:"bestAvg34Rank"`

	// x:正确用例数 y:错误用例数 z:平均map
	MPassCrashAvgMapX []int    `json:"passCrashAvgMapX"`
	MPassCrashAvgMapY []int    `json:"passCrashAvgMapY"`
	MPassCrashAvgMapZ []string `json:"passCrashAvgMapZ"`
}

// 统计随机采样结果
func StatisticRdStd(projectId string, analyzeFile string, factor string, tester string, casePrefix string) (*StaRd, *perrorx.ErrorX) {
	testPath := pdba.DBURLProjectIdTester(projectId, tester)
	analyzePath := path.Join(testPath, "rdanalyze_"+casePrefix+"_"+analyzeFile+"_"+factor+".txt")
	file, err := os.Open(analyzePath)
	if err != nil {
		return nil, perrorx.NewErrorXOpen(analyzePath, nil)
	}
	scanner := bufio.NewScanner(file)
	// 第一行四个部分, 组合总数n(拆分成正确用例数和错误用例数), 每个组合下的随机次数m, 聚焦的语句文件(不能为空), 聚焦语句数量s
	scanner.Scan()
	line := strings.TrimSpace(scanner.Text())
	sp := strings.Split(line, " ")
	if len(sp) != 5 {
		return nil, perrorx.NewErrorXSplitN(len(sp), 5, nil)
	}
	passNum, err1 := strconv.Atoi(sp[0])
	crashNum, err2 := strconv.Atoi(sp[1])
	randNum, err3 := strconv.Atoi(sp[2])
	statementNum, err4 := strconv.Atoi(sp[4])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return nil, perrorx.NewErrorXAtoI(sp[0]+"|"+sp[1]+"|"+
			sp[2]+"|"+sp[4], nil)
	}
	n := passNum * crashNum
	m := randNum
	s := statementNum
	// 统计最佳平均map及其对应的正确/错误用例数
	bestAvgMap := 0.0
	bestAvgPassNum := -1
	bestAvgCrashNum := -1
	// 记录最佳平均map的中位数, max, min, 1/4, 3/4及其排名
	bestAvgMidMap := 0.0
	bestAvgMidRank := make([]int, 0)
	bestAvgMaxMap := 0.0
	bestAvgMaxRank := make([]int, 0)
	bestAvgMinMap := 0.0
	bestAvgMinRank := make([]int, 0)
	bestAvg14Map := 0.0
	bestAvg14Rank := make([]int, 0)
	bestAvg34Map := 0.0
	bestAvg34Rank := make([]int, 0)
	// 统计x:正确用例数 y:错误用例数 z:平均map
	passCrashAvgMapX := make([]int, n)
	passCrashAvgMapY := make([]int, n)
	passCrashAvgMapZ := make([]string, n)
	// 接下来n个组合, 每个组合:
	for i := 0; i < n; i++ {
		// 	第一行两个数, 正确用例数, 错误用例数
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		sp := strings.Split(line, " ")
		if len(sp) != 2 {
			return nil, perrorx.NewErrorXSplitN(len(sp), 2, nil)
		}
		curPassNum, err1 := strconv.Atoi(sp[0])
		curCrashNum, err2 := strconv.Atoi(sp[1])
		if err1 != nil || err2 != nil {
			return nil, perrorx.NewErrorXAtoI(sp[0]+"|"+sp[1], nil)
		}
		// 计算每个组合的平均map
		avgMap := 0.0
		// 保存每次随机的map和rank
		mapRanks := make(sortPairs, 0)
		// 	接下来m次随机, 每次随机:
		for j := 0; j < m; j++ {
			// s个数, 代表聚焦语句文件中每个聚焦语句的排名
			scanner.Scan()
			line := strings.TrimSpace(scanner.Text())
			sp := strings.Split(line, " ")
			if len(sp) != s {
				return nil, perrorx.NewErrorXSplitN(len(sp), s, nil)
			}
			rank := make([]int, s)
			for k := 0; k < s; k++ {
				rank[k], err = strconv.Atoi(sp[k])
				if err != nil {
					return nil, perrorx.NewErrorXAtoI(sp[k], nil)
				}
			}
			// 计算每轮随机的map
			curMap := 0.0
			for k := 0; k < s; k++ {
				// 这里注意用k+1, 判断rank[k]是否为0
				if rank[k] != 0 {
					curMap += float64(k+1) / float64(rank[k])
				}
			}
			curMap /= float64(s)
			// 计算平均map
			avgMap += curMap
			// 记录map和rank
			mapRanks = append(mapRanks, sortPair{curMap, rank})
		}
		// 计算平均map
		avgMap /= float64(m)
		// 更新最优平均map
		if avgMap > bestAvgMap {
			bestAvgMap = avgMap
			bestAvgPassNum = curPassNum
			bestAvgCrashNum = curCrashNum
			// 根据mapRanks获取中位, max, min, 1/4, 3/4对应的rank, 更新bestAvgMidMap
			sort.Sort(mapRanks)
			bestAvgMidMap = mapRanks[m/2].MValue
			bestAvgMidRank = mapRanks[m/2].MRank
			bestAvgMaxMap = mapRanks[m-1].MValue
			bestAvgMaxRank = mapRanks[m-1].MRank
			bestAvgMinMap = mapRanks[0].MValue
			bestAvgMinRank = mapRanks[0].MRank
			bestAvg14Map = mapRanks[m/4].MValue
			bestAvg14Rank = mapRanks[m/4].MRank
			bestAvg34Map = mapRanks[m/4*3].MValue
			bestAvg34Rank = mapRanks[m/4*3].MRank
		}
		// 更新passCrashAvgMap
		passCrashAvgMapX[i] = curPassNum
		passCrashAvgMapY[i] = curCrashNum
		passCrashAvgMapZ[i] = strconv.FormatFloat(avgMap, 'f', 5, 64)
	}
	ans := &StaRd{
		MPassNum:          passNum,
		MCrashNum:         crashNum,
		MRandNum:          randNum,
		MBestAvgMap:       strconv.FormatFloat(bestAvgMap, 'f', 5, 64),
		MBestAvgPassNum:   bestAvgPassNum,
		MBestAvgCrashNum:  bestAvgCrashNum,
		MBestAvgMidMap:    strconv.FormatFloat(bestAvgMidMap, 'f', 5, 64),
		MBestAvgMidRank:   bestAvgMidRank,
		MBestAvgMaxMap:    strconv.FormatFloat(bestAvgMaxMap, 'f', 5, 64),
		MBestAvgMaxRank:   bestAvgMaxRank,
		MBestAvgMinMap:    strconv.FormatFloat(bestAvgMinMap, 'f', 5, 64),
		MBestAvgMinRank:   bestAvgMinRank,
		MBestAvg14Map:     strconv.FormatFloat(bestAvg14Map, 'f', 5, 64),
		MBestAvg14Rank:    bestAvg14Rank,
		MBestAvg34Map:     strconv.FormatFloat(bestAvg34Map, 'f', 5, 64),
		MBestAvg34Rank:    bestAvg34Rank,
		MPassCrashAvgMapX: passCrashAvgMapX,
		MPassCrashAvgMapY: passCrashAvgMapY,
		MPassCrashAvgMapZ: passCrashAvgMapZ,
	}
	return ans, nil
}
