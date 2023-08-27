package ahocorasick

import (
	_ "embed"
	"fmt"
	"github.com/tmikus/ahocorasick_rs/matchkind"
	"testing"

	ac "github.com/BobuSumisu/aho-corasick"
)

func ExampleAhoCorasick_FindAll_basic() {
	automaton := NewAhoCorasickBuilder().SetMatchKind(matchkind.Standard).Build([]string{"append", "appendage", "app"})
	defer automaton.Close()
	haystack := "append the app to the appendage"
	matches := automaton.FindAll(haystack)
	for _, match := range matches {
		fmt.Println(match.PatternIndex, haystack[match.Start:match.End], match.Start, match.End)
	}
	// Output:
	// 2 app 0 3
	// 2 app 11 14
	// 2 app 22 25
}

func ExampleAhoCorasick_FindAll_leftmost_first() {
	automaton := NewAhoCorasickBuilder().SetMatchKind(matchkind.LeftMostFirst).Build([]string{"append", "appendage", "app"})
	defer automaton.Close()
	haystack := "append the app to the appendage"
	matches := automaton.FindAll(haystack)
	for _, match := range matches {
		fmt.Println(match.PatternIndex, haystack[match.Start:match.End], match.Start, match.End)
	}
	// Output:
	// 0 append 0 6
	// 2 app 11 14
	// 0 append 22 28
}

func ExampleAhoCorasick_FindAll_leftmost_longest() {
	automaton := NewAhoCorasickBuilder().SetMatchKind(matchkind.LeftMostLongest).Build([]string{"append", "appendage", "app"})
	defer automaton.Close()
	haystack := "append the app to the appendage"
	matches := automaton.FindAll(haystack)
	for _, match := range matches {
		fmt.Println(match.PatternIndex, haystack[match.Start:match.End], match.Start, match.End)
	}
	// Output:
	// 0 append 0 6
	// 2 app 11 14
	// 1 appendage 22 31
}

func ExampleAhoCorasick_FindFirst_basic() {
	automaton := NewAhoCorasickBuilder().SetMatchKind(matchkind.Standard).Build([]string{"b", "abc", "abcd"})
	defer automaton.Close()
	haystack := "abcd"
	match := automaton.FindFirst(haystack)
	fmt.Println(haystack[match.Start:match.End])
	// Output: b
}

func ExampleAhoCorasick_FindFirst_leftmost_first() {
	automaton := NewAhoCorasickBuilder().SetMatchKind(matchkind.LeftMostFirst).Build([]string{"b", "abc", "abcd"})
	defer automaton.Close()
	haystack := "abcd"
	match := automaton.FindFirst(haystack)
	fmt.Println(haystack[match.Start:match.End])
	// Output: abc
}

func ExampleAhoCorasick_FindFirst_leftmost_longest() {
	automaton := NewAhoCorasickBuilder().SetMatchKind(matchkind.LeftMostLongest).Build([]string{"b", "abc", "abcd"})
	defer automaton.Close()
	haystack := "abcd"
	match := automaton.FindFirst(haystack)
	fmt.Println(haystack[match.Start:match.End])
	// Output: abcd
}

func ExampleAhoCorasick_IsMatch() {
	automaton := NewAhoCorasick([]string{"foo", "bar", "quux", "baz"})
	defer automaton.Close()
	fmt.Println(automaton.IsMatch("xxx bar xxx"))
	fmt.Println(automaton.IsMatch("xxx qux xxx"))
	// Output:
	// true
	// false
}

//go:embed data/sherlock.txt
var SHERLOCK string

var patterns = []string{
	"ADL", "ADl", "AdL", "Adl", "BAK", "BAk", "BAK", "BaK", "Bak",
	"BaK", "HOL", "HOl", "HoL", "Hol", "IRE", "IRe", "IrE", "Ire",
	"JOH", "JOh", "JoH", "Joh", "SHE", "SHe", "ShE", "She", "WAT",
	"WAt", "WaT", "Wat", "aDL", "aDl", "adL", "adl", "bAK", "bAk",
	"bAK", "baK", "bak", "baK", "hOL", "hOl", "hoL", "hol", "iRE",
	"iRe", "irE", "ire", "jOH", "jOh", "joH", "joh", "sHE", "sHe",
	"shE", "she", "wAT", "wAt", "waT", "wat", "ſHE", "ſHe", "ſhE",
	"ſhe",
}

func BenchmarkAhoCorasickGo(b *testing.B) {
	trie := ac.NewTrieBuilder().AddStrings(patterns).Build()

	for n := 0; n < b.N; n++ {
		trie.MatchString(SHERLOCK)
	}
}

func BenchmarkAhoCorasickRs(b *testing.B) {
	automaton := NewAhoCorasick(patterns)
	defer automaton.Close()

	for n := 0; n < b.N; n++ {
		automaton.FindAll(SHERLOCK)
	}
}
