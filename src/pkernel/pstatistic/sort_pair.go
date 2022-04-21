package pstatistic

type sortPair struct {
	MValue float64
	MRank []int
}

type sortPairs []sortPair

func (s sortPairs) Len() int           { return len(s) }
func (s sortPairs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortPairs) Less(i, j int) bool { return s[i].MValue < s[j].MValue }

