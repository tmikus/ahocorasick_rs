package ahocorasickkind

// AhoCorasickKind is principally used as an input to the [AhoCorasickBuilder.SetStartKind] method.
// Its documentation goes into more detail about each choice.
type AhoCorasickKind int

const (
	NonContinuousNFA AhoCorasickKind = 1 // Use a noncontiguous NFA.
	ContinuousNFA    AhoCorasickKind = 2 // Use a contiguous NFA.
	DFA              AhoCorasickKind = 3 // Use a DFA. Warning: DFAs typically use a large amount of memory.
)
