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
	"unsafe"
)

// Match represents a match found by an AhoCorasick automaton.
type Match struct {
	End          uint
	PatternIndex uint
	Start        uint
}

// AhoCorasick is an automaton for searching multiple strings in linear time.
//
// The AhoCorasick type supports a few basic ways of constructing an automaton, with the default being NewAhoCorasick.
// However, there are a fair number of configurable options that can be set by using AhoCorasickBuilder instead.
// Such options include, but are not limited to, how matches are determined, simple case insensitivity,
// whether to use a DFA or not and various knobs for controlling the space-vs-time trade-offs taken when building the automaton.
//
// Make sure to call AhoCorasick.Close() when you are done with the automaton.
type AhoCorasick struct {
	automaton *C.AhoCorasick
}

// NewAhoCorasick creates a new Aho-Corasick automaton using the default configuration.
//
// The default configuration optimizes for less space usage, but at the expense of longer search times.
// To change the configuration, use AhoCorasickBuilder.
//
// This uses the default matchkind.Standard match semantics, which reports a match as soon as it is found.
// This corresponds to the standard match semantics supported by textbook descriptions of the Aho-Corasick algorithm.
//
// Make sure to call AhoCorasick.Close() when you are done with the automaton.
func NewAhoCorasick(patterns []string) *AhoCorasick {
	cPatterns := make([]*C.char, len(patterns))
	for i, pattern := range patterns {
		cPatterns[i] = C.CString(pattern)
		defer C.free(unsafe.Pointer(cPatterns[i]))
	}

	automaton := C.create_automaton((**C.char)(&cPatterns[0]), C.size_t(len(patterns)))

	return &AhoCorasick{
		automaton: automaton,
	}
}

// FindAll returns an iterator of non-overlapping matches, using the match semantics that this automaton was constructed with.
//
// input may be any type that is cheaply convertible to an Input. This includes, but is not limited to, &str and &[u8].
//
// This is the infallible version of AhoCorasick.TryFindIter.
func (ac *AhoCorasick) FindAll(input string) []Match {
	cText := C.CString(input)
	defer C.free(unsafe.Pointer(cText))
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
// This is the infallible version of AhoCorasick.TryFind.
func (ac *AhoCorasick) FindFirst(input string) *Match {
	cText := C.CString(input)
	defer C.free(unsafe.Pointer(cText))
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

// IsMatch returns true if and only if this automaton matches the haystack at any position.
//
// Input may be any type that is cheaply convertible to an Input. This includes, but is not limited to, &str and &[u8].
//
// Aside from convenience, when AhoCorasick was built with leftmost-first or leftmost-longest semantics,
// this might result in a search that visits less of the haystack than AhoCorasick.FindFirst would otherwise.
// (For standard semantics, matches are always immediately returned once they are seen, so there is no way for this to do less work in that case.)
//
// Note that there is no corresponding fallible routine for this method. If you need a fallible version of this,
// then AhoCorasick.TryFind can be used with Input::earliest enabled.
func (ac *AhoCorasick) IsMatch(text string) bool {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	isMatch := C.is_match(ac.automaton, cText, C.size_t(len(text)))
	return int(isMatch) != 0
}

// Close frees the memory associated with this automaton.
func (ac *AhoCorasick) Close() {
	C.free_automaton(ac.automaton)
}
