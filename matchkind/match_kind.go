package matchkind

type MatchKind int

const (
	Standard        MatchKind = 1
	LeftMostLongest MatchKind = 2
	LeftMostFirst   MatchKind = 3
)
