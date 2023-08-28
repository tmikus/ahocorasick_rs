package ahocorasick

/*
#cgo darwin,arm64 LDFLAGS: -L./lib/darwin -lahocorasick_rs_arm64
#cgo darwin,amd64 LDFLAGS: -L./lib/darwin -lahocorasick_rs_amd64
#cgo linux,arm64 LDFLAGS: -L./lib/linux -lahocorasick_rs_arm64
#cgo linux,amd64 LDFLAGS: -L./lib/linux -lahocorasick_rs_amd64
#cgo windows,amd64 LDFLAGS: -L./lib/windows -lahocorasick_rs_amd64
#include "./lib/ahocorasick_rs.h"
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Match represents a match found by an [AhoCorasick] automaton.
type Match struct {
	// The ending position of the match.
	End uint
	// Returns the ID of the pattern that matched.
	//
	// The ID of a pattern is derived from the position in which it was originally inserted into the corresponding searcher. The first pattern has identifier 0, and each subsequent pattern is 1, 2 and so on.
	PatternIndex uint
	// The starting position of the match.
	Start uint
}

// AhoCorasick is an automaton for searching multiple strings in linear time.
//
// The [AhoCorasick] type supports a few basic ways of constructing an automaton, with the default being [NewAhoCorasick].
// However, there are a fair number of configurable options that can be set by using [AhoCorasickBuilder] instead.
// Such options include, but are not limited to, how matches are determined, simple case insensitivity,
// whether to use a AhoCorasickKindDFA or not and various knobs for controlling the space-vs-time trade-offs taken when building the automaton.
type AhoCorasick struct {
	automaton *C.AhoCorasick
}

// NewAhoCorasick creates a new Aho-Corasick automaton using the default configuration.
//
// The default configuration optimizes for less space usage, but at the expense of longer search times.
// To change the configuration, use [AhoCorasickBuilder].
//
// This uses the default [matchkind.MatchKindStandard] match semantics, which reports a match as soon as it is found.
// This corresponds to the standard match semantics supported by textbook descriptions of the Aho-Corasick algorithm.
func NewAhoCorasick(patterns []string) *AhoCorasick {
	cPatterns := make([]*C.char, len(patterns))
	cLengths := make([]C.size_t, len(patterns))
	for i, pattern := range patterns {
		cPatterns[i] = (*C.char)(unsafe.Pointer(unsafe.StringData(pattern)))
		cLengths[i] = C.size_t(len(pattern))
	}
	automaton := C.create_automaton(
		(**C.char)(&cPatterns[0]),
		(*C.size_t)(&cLengths[0]),
		C.size_t(len(patterns)),
	)
	result := &AhoCorasick{
		automaton: automaton,
	}
	runtime.SetFinalizer(result, func(c *AhoCorasick) {
		C.free_automaton(c.automaton)
	})
	return result
}

// FindAll returns an iterator of non-overlapping matches, using the match semantics that this automaton was constructed with.
//
// input may be any type that is cheaply convertible to an Input. This includes, but is not limited to, &str and &[u8].
//
// This is the infallible version of [AhoCorasick.TryFindIter].
func (ac *AhoCorasick) FindAll(input string) []Match {
	cText := (*C.char)(unsafe.Pointer(unsafe.StringData(input)))
	foundCount := C.long(0)
	cMatches := C.find_iter(ac.automaton, cText, C.size_t(len(input)), &foundCount)
	result := make([]Match, int(foundCount))
	if foundCount > 0 {
		goSlice := (*[1 << 30]C.AhoCorasickMatch)(unsafe.Pointer(cMatches))[:foundCount:foundCount]
		for i, val := range goSlice {
			result[i] = Match{
				End:          uint(val.end),
				PatternIndex: uint(val.pattern_index),
				Start:        uint(val.start),
			}
		}
		C.free(unsafe.Pointer(cMatches))
	}
	return result
}

// FindFirst returns the location of the first match according to the match semantics that this automaton was constructed with.
//
// input may be any type that is cheaply convertible to an Input. This includes, but is not limited to, &str and &[u8].
//
// This is the infallible version of [AhoCorasick.TryFind].
func (ac *AhoCorasick) FindFirst(input string) *Match {
	cText := (*C.char)(unsafe.Pointer(unsafe.StringData(input)))
	match := C.find(ac.automaton, cText, C.size_t(len(input)))
	if match == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(match))
	return &Match{
		End:          uint(match.end),
		PatternIndex: uint(match.pattern_index),
		Start:        uint(match.start),
	}
}

// GetKind returns the kind of the [AhoCorasick] automaton used by this searcher.
//
// Knowing the Aho-Corasick kind is principally useful for diagnostic purposes. In particular, if no specific kind
// was given to [AhoCorasickBuilder.SetKind], then one is automatically chosen and this routine will report which one.
//
// Note that the heuristics used for choosing which [ahocorasickkind.AhoCorasickKind] may be changed in a semver compatible release.
func (ac *AhoCorasick) GetKind() AhoCorasickKind {
	kind := C.get_kind(ac.automaton)
	return AhoCorasickKind(kind)
}

// IsMatch returns true if and only if this automaton matches the haystack at any position.
//
// Input may be any type that is cheaply convertible to an Input. This includes, but is not limited to, &str and &[u8].
//
// Aside from convenience, when [AhoCorasick] was built with leftmost-first or leftmost-longest semantics,
// this might result in a search that visits less of the haystack than [AhoCorasick.FindFirst] would otherwise.
// (For standard semantics, matches are always immediately returned once they are seen, so there is no way for this to do less work in that case.)
//
// Note that there is no corresponding fallible routine for this method. If you need a fallible version of this,
// then [AhoCorasick.TryFind] can be used with Input::earliest enabled.
func (ac *AhoCorasick) IsMatch(input string) bool {
	cText := (*C.char)(unsafe.Pointer(unsafe.StringData(input)))
	isMatch := C.is_match(ac.automaton, cText, C.size_t(len(input)))
	return int(isMatch) != 0
}
