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

// 普通分析结果
type StaNormal struct {
	MMap string `json:"map"`
	MRank []int `json:"rank"`
}

// 统计普通分析结果
func StatisticNormalStd(projectId string, analyzeFile string, factor string, casePrefix string) (*StaNormal, *perrorx.ErrorX) {
	testCasePath := pdba.DBURLProjectIdTestcases(projectId)
	analyzePath := path.Join(testCasePath, "analyze_"+casePrefix+"_"+analyzeFile+"_"+factor+".txt")
	file, err := os.Open(analyzePath)
	if err != nil {
		return nil, perrorx.NewErrorXOpen(analyzePath, nil)
	}
	scanner := bufio.NewScanner(file)
	// 第一行一个数, 聚焦语句数量s
	scanner.Scan()
	line := strings.TrimSpace(scanner.Text())
	s, err := strconv.Atoi(line)
	if err != nil {
		return nil, perrorx.NewErrorXAtoI(line, nil)
	}
	// s个数, 代表聚焦语句文件中每个聚焦语句的排名
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
	// 计算map
	rank_ := make([]int, 0)
	for _, r := range rank {
		if r <= 0 {
			continue
		}
		rank_ = append(rank_, r)
	}
	sort.Ints(rank_)
	myMap := 0.0
	for j := 0; j < len(rank_); j++ {
		// 这里注意用k+1
		myMap += float64(j+1) / float64(rank_[j])
	}
	myMap /= float64(s)
	ans := &StaNormal {
		MMap: strconv.FormatFloat(myMap, 'f', 5, 64),
		MRank: rank,
	}
	return ans, nil
}

