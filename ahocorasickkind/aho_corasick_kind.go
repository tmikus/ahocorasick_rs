package ahocorasickkind

type AhoCorasickKind int

const (
	NonContinuousNFA AhoCorasickKind = 1
	ContinuousNFA    AhoCorasickKind = 2
	DFA              AhoCorasickKind = 3
)
