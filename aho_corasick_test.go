package ahocorasick

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func ExampleAhoCorasick_FindAll_basic() {
	automaton := NewAhoCorasickBuilder().SetMatchKind(MatchKindStandard).Build([]string{"append", "appendage", "app"})
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
	automaton := NewAhoCorasickBuilder().SetMatchKind(MatchKindLeftMostFirst).Build([]string{"append", "appendage", "app"})
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
	automaton := NewAhoCorasickBuilder().SetMatchKind(MatchKindLeftMostLongest).Build([]string{"append", "appendage", "app"})
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
	automaton := NewAhoCorasickBuilder().SetMatchKind(MatchKindStandard).Build([]string{"b", "abc", "abcd"})
	haystack := "abcd"
	match := automaton.FindFirst(haystack)
	fmt.Println(haystack[match.Start:match.End])
	// Output: b
}

func ExampleAhoCorasick_FindFirst_leftmost_first() {
	automaton := NewAhoCorasickBuilder().SetMatchKind(MatchKindLeftMostFirst).Build([]string{"b", "abc", "abcd"})
	haystack := "abcd"
	match := automaton.FindFirst(haystack)
	fmt.Println(haystack[match.Start:match.End])
	// Output: abc
}

func ExampleAhoCorasick_FindFirst_leftmost_longest() {
	automaton := NewAhoCorasickBuilder().SetMatchKind(MatchKindLeftMostLongest).Build([]string{"b", "abc", "abcd"})
	haystack := "abcd"
	match := automaton.FindFirst(haystack)
	fmt.Println(haystack[match.Start:match.End])
	// Output: abcd
}

func ExampleAhoCorasick_GetKind() {
	automaton := NewAhoCorasick([]string{"foo", "bar", "quux", "baz"})
	fmt.Println(automaton.GetKind() == AhoCorasickKindDFA)
	// Output: true
}

func ExampleAhoCorasick_IsMatch() {
	automaton := NewAhoCorasick([]string{"foo", "bar", "quux", "baz"})
	fmt.Println(automaton.IsMatch("xxx bar xxx"))
	fmt.Println(automaton.IsMatch("xxx qux xxx"))
	// Output:
	// true
	// false
}

func TestAhoCorasick(t *testing.T) {
	Convey("GIVEN a list of 1000 patterns", t, func() {
		patterns := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			patterns[i] = fmt.Sprintf("pattern_%d", i)
		}

		Convey("WHEN a new AhoCorasick is created", func() {
			automaton := NewAhoCorasick(patterns)

			Convey("THEN the automaton should be able to match all patterns", func() {
				for i := 0; i < 1000; i++ {
					So(automaton.IsMatch(fmt.Sprintf("pattern_%d", i)), ShouldBeTrue)
				}
			})
		})
	})
}
