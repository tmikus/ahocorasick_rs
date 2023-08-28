package ahocorasick

// AhoCorasickKind is principally used as an input to the [AhoCorasickBuilder.SetStartKind] method.
// Its documentation goes into more detail about each choice.
type AhoCorasickKind int

const (
	AhoCorasickKindNonContinuousNFA AhoCorasickKind = 1 // Use a noncontiguous NFA.
	AhoCorasickKindContinuousNFA    AhoCorasickKind = 2 // Use a contiguous NFA.
	AhoCorasickKindDFA              AhoCorasickKind = 3 // Use a AhoCorasickKindDFA. Warning: DFAs typically use a large amount of memory.
)
