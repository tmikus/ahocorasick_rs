package main

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
	"fmt"
	"unsafe"
)

type Match struct {
	End          uint
	PatternIndex uint
	Start        uint
}

type AhoCorasick struct {
	automaton *C.AhoCorasick
}

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

func (ac *AhoCorasick) FindFirst(text string) *Match {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	match := C.find(ac.automaton, cText, C.size_t(len(text)))
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

func (ac *AhoCorasick) IsMatch(text string) bool {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	isMatch := C.is_match(ac.automaton, cText, C.size_t(len(text)))
	return int(isMatch) != 0
}

func (ac *AhoCorasick) Search(text string) []Match {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	foundCount := C.long(0)
	cMatches := C.find_iter(ac.automaton, cText, C.size_t(len(text)), &foundCount)
	goSlice := (*[1 << 30]C.AhoCorasickMatch)(unsafe.Pointer(cMatches))[:foundCount:foundCount]

	result := make([]Match, int(foundCount))
	for i, val := range goSlice {
		result[i] = Match{
			End:          uint(val.end),
			PatternIndex: uint(val.pattern_index),
			Start:        uint(val.start),
		}
	}

	C.free(unsafe.Pointer(cMatches))

	return result
}

func (ac *AhoCorasick) Close() {
	C.free_automaton(ac.automaton)
}

func main() {
	patterns := []string{"foo", "bar", "baz"}
	text := "foobarbaz"

	aho := NewAhoCorasick(patterns)
	match := aho.FindFirst(text)
	if match != nil {
		fmt.Printf("Found match: %v\n", match)
	} else {
		fmt.Println("No match found")
	}

	matches := aho.Search(text)
	defer aho.Close()

	fmt.Printf("Found matches: %v\n", matches)
}
