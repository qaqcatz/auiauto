package pstatistic

import (
	"auiauto/perrorx"
	"strconv"
)

// 汇聚分析结果(rd, art, two, combine)
type StaAll struct {
	MStaEachs []*StaEach `json:"staEachs"`
	// avg map
	MArtAvgMap string       `json:"artAvgMap"`
	MTwoAvgMap string       `json:"twoAvgMap"`
	MCombineBestMapAvgMap string `json:"combineBestMapAvgMap"`
	MRdBestAvgMapAvgMap string `json:"rdBestAvgMapAvgMap"`
	// mrr
	MArtMrr string       `json:"artMrr"`
	MTwoMrr string       `json:"twoMrr"`
	MCombineBestMapMrr string `json:"combineBestMapMrr"`
	MRdBestAvgMidMapMrr string `json:"rdBestAvgMidMapMrr"`
	MRdBestAvgMaxMapMrr string `json:"rdBestAvgMaxMapMrr"`
	MRdBestAvg34MapMrr string `json:"rdBestAvg34MapMrr"`
	// top: 1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50
	MArtTop []int `json:"artTop"`
	MTwoTop []int `json:"twoTop"`
	MCombineBestMapTop []int `json:"combineTop"`
	MRdBestAvgMidMapTop []int `json:"rdBestAvgMidMapTop"`
	MRdBestAvgMaxMapTop []int `json:"rdBestAvgMaxMapTop"`
	MRdBestAvg34MapTop []int `json:"rdBestAvg34MapTop"`
}

// 汇聚分析中的单个项目
type StaEach struct {
	MProjectId string `json:"projectId"`
	// rd
	// 各个正确/错误用例组合中的最优平均map
	MRdBestAvgMap string `json:"rdBestAvgMap"`
	// rdBestAvgMap对应的正确用例数
	MRdBestAvgPassNum int `json:"rdBestAvgPassNum"`
	// rdBestAvgMap对应的错误用例数
	MRdBestAvgCrashNum int `json:"rdBestAvgCrashNum"`
	// rdBestAvgMap对应中的中位数
	MRdBestAvgMidMap string `json:"rdBestAvgMidMap"`
	// rdBestAvgMidMap对应的rank
	MRdBestAvgMidRank []int `json:"rdBestAvgMidRank"`
	// rdBestAvgMidMap对应的最优值与rank
	MRdBestAvgMaxMap  string `json:"rdBestAvgMaxMap"`
	MRdBestAvgMaxRank []int  `json:"rdBestAvgMaxRank"`
	// rdBestAvgMidMap对应的最差值与rank
	MRdBestAvgMinMap  string `json:"rdBestAvgMinMap"`
	MRdBestAvgMinRank []int  `json:"rdBestAvgMinRank"`
	// rdBestAvgMidMap对应的1/4线与rank
	MRdBestAvg14Map  string `json:"rdBestAvg14Map"`
	MRdBestAvg14Rank []int  `json:"rdBestAvg14Rank"`
	// rdBestAvgMidMap对应的3/4线与rank
	MRdBestAvg34Map  string `json:"rdBestAvg34Map"`
	MRdBestAvg34Rank []int  `json:"rdBestAvg34Rank"`
	// art
	MArtMap string `json:"artMap"`
	MArtRank []int `json:"artRank"`
	// two
	MTwoMap string `json:"twoMap"`
	MTwoRank []int `json:"twoRank"`
	// combine
	// 所有组合中最好的map
	MCombineBestMap string `json:"combineBestMap"`
	// combineBestMap对应的正确用例数量
	MCombineBestPassNum int `json:"combineBestPassNum"`
	// combineBestMap对应的错误用例数量
	MCombineBestCrashNum int `json:"combineBestCrashNum"`
	// combineBestMap对应的rank
	MCombineBestRank []int `json:"combineBestRank"`
}

// 汇聚分析结果统计(rd, art, two, combine)
func StatisticAllStd(analyzeFile string, factor string, tester string) (*StaAll, *perrorx.ErrorX) {
	projects, err := getAllProjects()
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	staEachs := make([]*StaEach, 0)
	// 统计map
	artAvgMap := 0.0
	twoAvgMap := 0.0
	combineBestMapAvgMap := 0.0
	rdBestAvgMapAvgMap := 0.0
	// 统计mrr
	artMrr := 0.0
	twoMrr := 0.0
	combineBestMapMrr := 0.0
	rdBestAvgMidMapMrr := 0.0
	rdBestAvgMaxMapMrr := 0.0
	rdBestAvg34MapMrr := 0.0
	// 统计top
	artTop := make([]int, 11)
	twoTop := make([]int, 11)
	combineBestMapTop := make([]int, 11)
	rdBestAvgMidMapTop := make([]int, 11)
	rdBestAvgMaxMapTop := make([]int, 11)
	rdBestAvg34MapTop := make([]int, 11)
	for _, projectId := range projects {
		rdAns, err := StatisticRdStd(projectId, analyzeFile, factor, tester, "rd")
		if err != nil {
			return nil, perrorx.TransErrorX(err)
		}
		artAns, err := StatisticNormalStd(projectId, analyzeFile, factor, "art")
		if err != nil {
			return nil, perrorx.TransErrorX(err)
		}
		twoAns, err := StatisticNormalStd(projectId, analyzeFile, factor, "two")
		if err != nil {
			return nil, perrorx.TransErrorX(err)
		}
		combineAns, err := StatisticCombineStd(projectId, analyzeFile, factor, "art")
		if err != nil {
			return nil, perrorx.TransErrorX(err)
		}
		staEachs = append(staEachs, &StaEach {
			MProjectId: projectId,
			MRdBestAvgMap: rdAns.MBestAvgMap,
			MRdBestAvgPassNum: rdAns.MBestAvgPassNum,
			MRdBestAvgCrashNum: rdAns.MBestAvgCrashNum,
			MRdBestAvgMidMap: rdAns.MBestAvgMidMap,
			MRdBestAvgMidRank: rdAns.MBestAvgMidRank,
			MRdBestAvgMaxMap: rdAns.MBestAvgMap,
			MRdBestAvgMaxRank: rdAns.MBestAvgMaxRank,
			MRdBestAvgMinMap: rdAns.MBestAvgMinMap,
			MRdBestAvgMinRank: rdAns.MBestAvgMinRank,
			MRdBestAvg14Map: rdAns.MBestAvg14Map,
			MRdBestAvg14Rank: rdAns.MBestAvg14Rank,
			MRdBestAvg34Map: rdAns.MBestAvg34Map,
			MRdBestAvg34Rank: rdAns.MBestAvg34Rank,
			MArtMap: artAns.MMap,
			MArtRank: artAns.MRank,
			MTwoMap: twoAns.MMap,
			MTwoRank: twoAns.MRank,
			MCombineBestMap: combineAns.MBestMap,
			MCombineBestPassNum: combineAns.MBestPassNum,
			MCombineBestCrashNum: combineAns.MBestCrashNum,
			MCombineBestRank: combineAns.MBestRank,
		})
		// 计算art avg map
		curArtMap, err_ := strconv.ParseFloat(artAns.MMap, 64)
		artAvgMap += curArtMap
		if err_ != nil {
			return nil, perrorx.NewErrorXParseFloat(artAns.MMap, nil)
		}
		// 计算art mrr
		minArtRank := -1
		for i := 0; i < len(artAns.MRank); i++ {
			if artAns.MRank[i] == 0 {
				continue
			}
			if minArtRank == -1 || artAns.MRank[i] < minArtRank {
				minArtRank = artAns.MRank[i]
			}
		}
		if minArtRank != -1 {
			artMrr += 1.0/float64(minArtRank)
		}
		// 计算art的top
		if minArtRank >= 1 {
			if minArtRank == 1 {
				artTop[0]++
			} else {
				idx := (minArtRank-1)/5+1
				if idx < 11 {
					artTop[idx]++
				}
			}
		}
		// 计算two avg map
		curTwoMap, err_ := strconv.ParseFloat(twoAns.MMap, 64)
		twoAvgMap += curTwoMap
		if err_ != nil {
			return nil, perrorx.NewErrorXParseFloat(twoAns.MMap, nil)
		}
		// 计算two mrr
		minTwoRank := -1
		for i := 0; i < len(twoAns.MRank); i++ {
			if twoAns.MRank[i] == 0 {
				continue
			}
			if minTwoRank == -1 || twoAns.MRank[i] < minTwoRank {
				minTwoRank = twoAns.MRank[i]
			}
		}
		if minTwoRank != -1 {
			twoMrr += 1.0/float64(minTwoRank)
		}
		// 计算two的top
		if minTwoRank >= 1 {
			if minTwoRank == 1 {
				twoTop[0]++
			} else {
				idx := (minTwoRank-1)/5+1
				if idx < 11 {
					twoTop[idx]++
				}
			}
		}
		// 计算combine best map的avg map
		curCombineBestMapMap, err_ := strconv.ParseFloat(combineAns.MBestMap, 64)
		combineBestMapAvgMap += curCombineBestMapMap
		if err_ != nil {
			return nil, perrorx.NewErrorXParseFloat(combineAns.MBestMap, nil)
		}
		// 计算combine best map的mrr
		minCombineBestMapRank := -1
		for i := 0; i < len(combineAns.MBestRank); i++ {
			if combineAns.MBestRank[i] == 0 {
				continue
			}
			if minCombineBestMapRank == -1 || combineAns.MBestRank[i] < minCombineBestMapRank {
				minCombineBestMapRank = combineAns.MBestRank[i]
			}
		}
		if minCombineBestMapRank != -1 {
			combineBestMapMrr += 1.0/float64(minCombineBestMapRank)
		}
		// 计算combine best map的top
		if minCombineBestMapRank >= 1 {
			if minCombineBestMapRank == 1 {
				combineBestMapTop[0]++
			} else {
				idx := (minCombineBestMapRank-1)/5+1
				if idx < 11 {
					combineBestMapTop[idx]++
				}
			}
		}
		// 计算rd best avg map的avg map
		curRdBestAvgMapMap, err_ := strconv.ParseFloat(rdAns.MBestAvgMap, 64)
		rdBestAvgMapAvgMap += curRdBestAvgMapMap
		if err_ != nil {
			return nil, perrorx.NewErrorXParseFloat(rdAns.MBestAvgMap, nil)
		}
		// 计算rd best avg mid map的mrr
		minRdBestAvgMidMapRank := -1
		for i := 0; i < len(rdAns.MBestAvgMidRank); i++ {
			if rdAns.MBestAvgMidRank[i] == 0 {
				continue
			}
			if minRdBestAvgMidMapRank == -1 || rdAns.MBestAvgMidRank[i] < minRdBestAvgMidMapRank {
				minRdBestAvgMidMapRank = rdAns.MBestAvgMidRank[i]
			}
		}
		if minRdBestAvgMidMapRank != -1 {
			rdBestAvgMidMapMrr += 1.0/float64(minRdBestAvgMidMapRank)
		}
		// 计算rd best avg mid map的top
		if minRdBestAvgMidMapRank >= 1 {
			if minRdBestAvgMidMapRank == 1 {
				rdBestAvgMidMapTop[0]++
			} else {
				idx := (minRdBestAvgMidMapRank-1)/5+1
				if idx < 11 {
					rdBestAvgMidMapTop[idx]++
				}
			}
		}
		// 计算rd best avg max map的mrr
		minRdBestAvgMaxMapRank := -1
		for i := 0; i < len(rdAns.MBestAvgMaxRank); i++ {
			if rdAns.MBestAvgMaxRank[i] == 0 {
				continue
			}
			if minRdBestAvgMaxMapRank == -1 || rdAns.MBestAvgMaxRank[i] < minRdBestAvgMaxMapRank {
				minRdBestAvgMaxMapRank = rdAns.MBestAvgMaxRank[i]
			}
		}
		if minRdBestAvgMaxMapRank != -1 {
			rdBestAvgMaxMapMrr += 1.0/float64(minRdBestAvgMaxMapRank)
		}
		// 计算rd best avg max map的top
		if minRdBestAvgMaxMapRank >= 1 {
			if minRdBestAvgMaxMapRank == 1 {
				rdBestAvgMaxMapTop[0]++
			} else {
				idx := (minRdBestAvgMaxMapRank-1)/5+1
				if idx < 11 {
					rdBestAvgMaxMapTop[idx]++
				}
			}
		}
		// 计算rd best avg 3/4 map的mrr
		minRdBestAvg34MapRank := -1
		for i := 0; i < len(rdAns.MBestAvg34Rank); i++ {
			if rdAns.MBestAvg34Rank[i] == 0 {
				continue
			}
			if minRdBestAvg34MapRank == -1 || rdAns.MBestAvg34Rank[i] < minRdBestAvg34MapRank {
				minRdBestAvg34MapRank = rdAns.MBestAvg34Rank[i]
			}
		}
		if minRdBestAvg34MapRank != -1 {
			rdBestAvg34MapMrr += 1.0/float64(minRdBestAvg34MapRank)
		}
		// 计算rd best avg 3/4 map的top
		if minRdBestAvg34MapRank >= 1 {
			if minRdBestAvg34MapRank == 1 {
				rdBestAvg34MapTop[0]++
			} else {
				idx := (minRdBestAvg34MapRank-1)/5+1
				if idx < 11 {
					rdBestAvg34MapTop[idx]++
				}
			}
		}
	}
	// 计算art avg map
	artAvgMap /= float64(len(projects))
	// 计算two avg map
	twoAvgMap /= float64(len(projects))
	// 计算combine best map的avg map
	combineBestMapAvgMap /= float64(len(projects))
	// 计算rd best avg map的avg map
	rdBestAvgMapAvgMap /= float64(len(projects))
	// 计算art mrr
	artMrr /= float64(len(projects))
	// 计算two mrr
	twoMrr /= float64(len(projects))
	// 计算combine best map的mrr
	combineBestMapMrr /= float64(len(projects))
	// 计算rd best avg mid map的mrr
	rdBestAvgMidMapMrr /= float64(len(projects))
	// 计算rd best avg max map的mrr
	rdBestAvgMaxMapMrr /= float64(len(projects))
	// 计算rd best avg 3/4 map的mrr
	rdBestAvg34MapMrr /= float64(len(projects))
	// 计算各个top, 前缀和
	for i := 1; i < 11; i++ {
		artTop[i] = artTop[i] + artTop[i-1]
		twoTop[i] = twoTop[i] + twoTop[i-1]
		combineBestMapTop[i] = combineBestMapTop[i] + combineBestMapTop[i-1]
		rdBestAvgMidMapTop[i] = rdBestAvgMidMapTop[i] + rdBestAvgMidMapTop[i-1]
		rdBestAvgMaxMapTop[i] = rdBestAvgMaxMapTop[i] + rdBestAvgMaxMapTop[i-1]
		rdBestAvg34MapTop[i] = rdBestAvg34MapTop[i] + rdBestAvg34MapTop[i-1]
	}
	return &StaAll{
		MStaEachs: staEachs,
		MArtAvgMap: strconv.FormatFloat(artAvgMap, 'f', 5, 64),
		MTwoAvgMap: strconv.FormatFloat(twoAvgMap, 'f', 5, 64),
		MCombineBestMapAvgMap: strconv.FormatFloat(combineBestMapAvgMap, 'f', 5, 64),
		MRdBestAvgMapAvgMap: strconv.FormatFloat(rdBestAvgMapAvgMap, 'f', 5, 64),
		MArtMrr: strconv.FormatFloat(artMrr, 'f', 5, 64),
		MTwoMrr: strconv.FormatFloat(twoMrr, 'f', 5, 64),
		MCombineBestMapMrr: strconv.FormatFloat(combineBestMapMrr, 'f', 5, 64),
		MRdBestAvgMidMapMrr: strconv.FormatFloat(rdBestAvgMidMapMrr, 'f', 5, 64),
		MRdBestAvgMaxMapMrr: strconv.FormatFloat(rdBestAvgMaxMapMrr, 'f', 5, 64),
		MRdBestAvg34MapMrr: strconv.FormatFloat(rdBestAvg34MapMrr, 'f', 5, 64),
		MArtTop: artTop,
		MTwoTop: twoTop,
		MCombineBestMapTop: combineBestMapTop,
		MRdBestAvgMidMapTop: rdBestAvgMidMapTop,
		MRdBestAvgMaxMapTop: rdBestAvgMaxMapTop,
		MRdBestAvg34MapTop: rdBestAvg34MapTop,
	}, nil
}