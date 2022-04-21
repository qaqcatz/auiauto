package pstatistic

import (
	"auiauto/perrorx"
	"strconv"
)

var FACTORS = []string{
	"Ochiai",
	"Tarantula",
	"Barinel",
	"DStar",
	"Op2",
}

// 影响因子分析统计
type StaFactors struct {
	MFactors         []string          `json:"factors"`
	MStaFactorsEachs []*StaFactorsEach `json:"staFactorsEachs"`
	// mrr
	MMrrs []string `json:"mrrs"`
	// top: 1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50
	MTops [][]int `json:"tops"`
}

// 每个项目的影响因子分析统计
type StaFactorsEach struct {
	MProjectId string `json:"projectId"`
	MMaps  []string `json:"maps"`
	MRanks [][]int  `json:"ranks"`
}

// 因子分析结果统计
func StatisticFactorsStd(analyzeFile string, casePrefix string) (*StaFactors, *perrorx.ErrorX) {
	projects, err := getAllProjects()
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	staFactorsEachs := make([]*StaFactorsEach, 0)
	// 统计mrr
	mrrs := make([]float64, len(FACTORS))
	// 统计top
	tops := make([][]int, len(FACTORS))
	for i := 0; i < len(tops); i++ {
		tops[i] = make([]int, 11)
	}
	for _, projectId := range projects {
		staFactorsEach := &StaFactorsEach{
			MProjectId: projectId,
			MMaps:  make([]string, len(FACTORS)),
			MRanks: make([][]int, len(FACTORS)),
		}
		staFactorsEachs = append(staFactorsEachs, staFactorsEach)
		for i := 0; i < len(FACTORS); i++ {
			ans, err := StatisticNormalStd(projectId, analyzeFile, FACTORS[i], casePrefix)
			if err != nil {
				return nil, perrorx.TransErrorX(err)
			}
			staFactorsEach.MMaps[i] = ans.MMap
			staFactorsEach.MRanks[i] = ans.MRank
			// 计算mrr
			minRank := -1
			for j := 0; j < len(ans.MRank); j++ {
				if ans.MRank[j] == 0 {
					continue
				}
				if minRank == -1 || ans.MRank[j] < minRank {
					minRank = ans.MRank[j]
				}
			}
			if minRank != -1 {
				mrrs[i] += 1.0 / float64(minRank)
			}
			// 计算top
			if minRank >= 1 {
				if minRank == 1 {
					tops[i][0]++
				} else {
					idx := (minRank-1)/5 + 1
					if idx < 11 {
						tops[i][idx]++
					}
				}
			}
		}
	}
	// 计算mrr
	for i := 0; i < len(FACTORS); i++ {
		mrrs[i] /= float64(len(projects))
	}
	// 转成字符数组
	mrrsStr := make([]string, 0)
	for i := 0; i < len(FACTORS); i++ {
		fltStr := strconv.FormatFloat(mrrs[i], 'f', 5, 64)
		mrrsStr = append(mrrsStr, fltStr)
	}
	// 计算各个top, 前缀和
	for i := 0; i < len(FACTORS); i++ {
		for j := 1; j < 11; j++ {
			tops[i][j] = tops[i][j] + tops[i][j-1]
		}
	}
	return &StaFactors{
		MFactors: FACTORS,
		MStaFactorsEachs: staFactorsEachs,
		MMrrs: mrrsStr,
		MTops: tops,
	}, nil
}
