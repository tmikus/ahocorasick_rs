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

	maxMatches := 10000 // Adjust this based on your needs
	matches := make([]C.size_t, maxMatches)

	numMatches := C.search_automaton(ac.automaton, cText, C.size_t(len(text)), (*C.size_t)(&matches[0]))

	result := make([]int, numMatches)
	for i := 0; i < int(numMatches); i++ {
		result[i] = int(matches[i])
	}

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
