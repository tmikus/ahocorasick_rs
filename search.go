package main

/*
#cgo LDFLAGS: -L./lib -lahocorasick_rs
#include "./lib/ahocorasick_rs.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

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

func (ac *AhoCorasick) Search(text string) []int {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	foundCount := C.long(0)
	cMatches := C.search_automaton(ac.automaton, cText, C.size_t(len(text)), &foundCount)
	goSlice := (*[1 << 30]C.size_t)(unsafe.Pointer(cMatches))[:foundCount:foundCount]

	result := make([]int, int(foundCount))
	for i, val := range goSlice {
		result[i] = int(val)
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
	matches := aho.Search(text)
	defer aho.Close()

	fmt.Printf("Found matches at indexes: %v\n", matches)
}
