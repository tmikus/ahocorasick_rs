package ahocorasick

/*
#cgo LDFLAGS: -laho_corasick_ffi
#include "./ahocorasick_rs.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// AhoCorasickBuilder is a builder for configuring an [AhoCorasick] automaton.
type AhoCorasickBuilder struct {
	asciiCaseInsensitive bool
	byteClasses          bool
	denseDepth           *uint
	kind                 *AhoCorasickKind
	matchKind            MatchKind
	prefilter            bool
	startKind            StartKind
}

// NewAhoCorasickBuilder creates a new builder for configuring an [AhoCorasick] automaton.
//
// The builder provides a way to configure a number of things, including ASCII case insensitivity and what kind of match semantics are used.
func NewAhoCorasickBuilder() *AhoCorasickBuilder {
	return &AhoCorasickBuilder{
		asciiCaseInsensitive: false,
		byteClasses:          true,
		denseDepth:           nil,
		kind:                 nil,
		matchKind:            MatchKindStandard,
		prefilter:            true,
		startKind:            StartKindUnanchored,
	}
}

// Build creates an [AhoCorasick] automaton using the configuration set on this builder.
//
// A builder may be reused to create more automatons.
func (b *AhoCorasickBuilder) Build(patterns []string) *AhoCorasick {
	cPatterns := make([]*C.char, len(patterns))
	for i, pattern := range patterns {
		cPatterns[i] = C.CString(pattern)
		defer C.free(unsafe.Pointer(cPatterns[i]))
	}
	options := C.AhoCorasickBuilderOptions{
		ascii_case_insensitive: boolToCInt(b.asciiCaseInsensitive),
		byte_classes:           boolToCInt(b.byteClasses),
		dense_depth:            (*C.size_t)(unsafe.Pointer(b.denseDepth)),
		kind:                   (*C.size_t)(unsafe.Pointer(b.kind)),
		match_kind:             C.size_t(b.matchKind),
		prefilter:              boolToCInt(b.prefilter),
		start_kind:             C.size_t(b.startKind),
	}
	automaton := C.build_automaton(
		(**C.char)(&cPatterns[0]),
		C.size_t(len(patterns)),
		(*C.AhoCorasickBuilderOptions)(unsafe.Pointer(&options)),
	)
	return &AhoCorasick{
		automaton: automaton,
	}
}

// SetAsciiCaseInsensitive enables ASCII-aware case-insensitive matching.
//
// When this option is enabled, searching will be performed without respect to case for ASCII letters (a-z and A-Z) only.
//
// Enabling this option does not change the search algorithm, but it may increase the size of the automaton.
func (b *AhoCorasickBuilder) SetAsciiCaseInsensitive(asciiCaseInsensitive bool) *AhoCorasickBuilder {
	b.asciiCaseInsensitive = asciiCaseInsensitive
	return b
}

// SetByteClasses sets a debug setting for whether to attempt to shrink the size of the automaton’s alphabet or not.
//
// This option is enabled by default and should never be disabled unless one is debugging the underlying automaton.
//
// When enabled, some (but not all) [AhoCorasick] automatons will use a map from all possible bytes to their
// corresponding equivalence class.
// Each equivalence class represents a set of bytes that does not discriminate between a match and a non-match
// in the automaton.
//
// The advantage of this map is that the size of the transition table can be reduced drastically from
// #states * 256 * sizeof(u32) to #states * k * sizeof(u32) where k is the number of equivalence classes
// (rounded up to the nearest power of 2). As a result, total space usage can decrease substantially.
// Moreover, since a smaller alphabet is used, automaton compilation becomes faster as well.
//
// WARNING: This is only useful for debugging automatons. Disabling this does not yield any speed advantages.
// Namely, even when this is disabled, a byte class map is still used while searching. The only difference is that
// every byte will be forced into its own distinct equivalence class. This is useful for debugging the actual generated
// transitions because it lets one see the transitions defined on actual bytes instead of the equivalence classes.
func (b *AhoCorasickBuilder) SetByteClasses(byteClasses bool) *AhoCorasickBuilder {
	b.byteClasses = byteClasses
	return b
}

// SetDenseDepth sets the limit on how many states use a dense representation for their transitions.
// Other states will generally use a sparse representation.
//
// A dense representation uses more memory but is generally faster, since the next transition in a dense representation
// can be computed in a constant number of instructions. A sparse representation uses less memory but is generally
// slower, since the next transition in a sparse representation requires executing a variable number of instructions.
//
// This setting is only used when an [AhoCorasick] implementation is used that supports the dense versus sparse
// representation trade off. Not all do.
//
// This limit is expressed in terms of the depth of a state, i.e., the number of transitions from the starting
// state of the automaton. The idea is that most of the time searching will be spent near the starting state of the
// automaton, so states near the start state should use a dense representation. States further away from the start
// state would then use a sparse representation.
//
// By default, this is set to a low but non-zero number. Setting this to 0 is almost never what you want, since it is
// likely to make searches very slow due to the start state itself being forced to use a sparse representation.
// However, it is unlikely that increasing this number will help things much, since the most active states have
// a small depth. More to the point, the memory usage increases super-linearly as this number increases.
func (b *AhoCorasickBuilder) SetDenseDepth(denseDepth *uint) *AhoCorasickBuilder {
	b.denseDepth = denseDepth
	return b
}

// SetKind sets the type of underlying automaton to use.
//
// Currently, there are four choices:
//
// [ahocorasickkind.AhoCorasickKindNonContinuousNFA] instructs the searcher to use a noncontinuous::NFA. A noncontinuous NFA
// is the fastest to be built, has moderate memory usage and is typically the slowest to execute a search.
//
// [ahocorasickkind.AhoCorasickKindContinuousNFA] instructs the searcher to use a contiguous::NFA. A contiguous NFA is a little slower
// to build than a noncontinuous NFA, has excellent memory usage and is typically a little slower than a AhoCorasickKindDFA for a search.
//
// [ahocorasickkind.AhoCorasickKindDFA] instructs the searcher to use a dfa::AhoCorasickKindDFA. A AhoCorasickKindDFA is very slow to build, uses exorbitant
// amounts of memory, but will typically execute searches the fastest.
//
// nil (the default) instructs the searcher to choose the “best” Aho-Corasick implementation.
// This choice is typically based primarily on the number of patterns.
//
// Setting this configuration does not change the time complexity for constructing the Aho-Corasick automaton
// (which is O(p) where p is the total number of patterns being compiled). Setting this to [ahocorasickkind.AhoCorasickKindDFA] does
// however reduce the time complexity of non-overlapping searches from O(n + p) to O(n), where n is
// the length of the haystack.
//
// In general, you should probably stick to the default unless you have some kind of reason to use a specific
// Aho-Corasick implementation. For example, you might choose [ahocorasickkind.AhoCorasickKindDFA] if you don’t care about memory
// usage and want the fastest possible search times.
//
// Setting this guarantees that the searcher returned uses the chosen implementation. If that implementation could
// not be constructed, then an error will be returned. In contrast, when None is used, it is possible for it to attempt
// to construct, for example, a contiguous NFA and have it fail. In which case, it will fall back to using a noncontinuous NFA.
//
// If nil is given, then one may use [AhoCorasick.GetKind] to determine which [AhoCorasick] implementation was chosen.
//
// Note that the heuristics used for choosing which [ahocorasickkind.AhoCorasickKind] may be changed in a semver compatible release.
func (b *AhoCorasickBuilder) SetKind(kind *AhoCorasickKind) *AhoCorasickBuilder {
	b.kind = kind
	return b
}

// SetMatchKind sets the desired match semantics.
//
// The default is [matchkind.MatchKindStandard], which corresponds to the match semantics supported by the standard textbook
// description of the Aho-Corasick algorithm. Namely, matches are reported as soon as they are found.
// Moreover, this is the only way to get overlapping matches or do stream searching.
//
// The other kinds of match semantics that are supported are [matchkind.MatchKindLeftMostFirst] and [matchkind.MatchKindLeftMostLongest].
// The former corresponds to the match you would get if you were to try to match each pattern at each position
// in the haystack in the same order that you give to the automaton. That is, it returns the leftmost match corresponding
// to the earliest pattern given to the automaton. The latter corresponds to finding the longest possible match
// among all leftmost matches.
//
// For more details on match semantics, see the documentation for MatchKind.
//
// Note that setting this to [matchkind.MatchKindLeftMostFirst] or [matchkind.MatchKindLeftMostLongest] will cause some search routines
// on [AhoCorasick] to return an error (or panic if you’re using the infallible API). Notably, this includes
// stream and overlapping searches.
func (b *AhoCorasickBuilder) SetMatchKind(matchKind MatchKind) *AhoCorasickBuilder {
	b.matchKind = matchKind
	return b
}

// SetPrefilter enables heuristic prefilter optimizations.
//
// When enabled, searching will attempt to quickly skip to match candidates using specialized literal search routines.
// A prefilter cannot always be used, and is generally treated as a heuristic. It can be useful to disable this
// if the prefilter is observed to be suoptimal for a particular workload.
//
// Currently, prefilters are typically only active when building searchers with a small (less than 100) number of patterns.
//
// This is enabled by default.
func (b *AhoCorasickBuilder) SetPrefilter(prefilter bool) *AhoCorasickBuilder {
	b.prefilter = prefilter
	return b
}

// SetStartKind sets the starting state configuration for the automaton.
//
// Every Aho-Corasick automaton is capable of having two start states: one that is used for unanchored searches
// and one that is used for anchored searches. Some automatons, like the NFAs, support this with almost zero additional cost.
// Other automatons, like the AhoCorasickKindDFA, require two copies of the underlying transition table to support both simultaneously.
//
// Because there may be an added non-trivial cost to supporting both, it is possible to configure which starting
// state configuration is needed.
//
// Indeed, since anchored searches tend to be somewhat more rare, only unanchored searches are supported by default.
// Thus, [startkind.StartKindUnanchored] is the default.
//
// Note that when this is set to [startkind.StartKindUnanchored], then running an anchored search will result in an error
// (or a panic if using the infallible APIs). Similarly, when this is set to [startkind.StartKindAnchored], then running
// an unanchored search will result in an error (or a panic if using the infallible APIs).
// When [startkind.StartKindBoth] is used, then both unanchored and anchored searches are always supported.
//
// Also note that even if an [AhoCorasick] searcher is using an NFA internally (which always supports both unanchored
// and anchored searches), an error will still be reported for a search that isn’t supported by the configuration
// set via this method. This means, for example, that an error is never dependent on which internal
// implementation of [AhoCorasick] is used.
func (b *AhoCorasickBuilder) SetStartKind(startKind StartKind) *AhoCorasickBuilder {
	b.startKind = startKind
	return b
}

func boolToCInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
