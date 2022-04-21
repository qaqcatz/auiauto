package pcoverage

// 用于给覆盖结果排序, 主要是为了显示时好看一些
type SortPair struct {
	MSid     int
	MEventId string
}

type SortPairs []SortPair

func (s SortPairs) Len() int           { return len(s) }
func (s SortPairs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortPairs) Less(i, j int) bool { return s[i].MSid < s[j].MSid }
