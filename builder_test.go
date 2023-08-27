package ahocorasick

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tmikus/ahocorasick_rs/matchkind"
	"testing"
)

func ExampleAhoCorasickBuilder_SetAsciiCaseInsensitive() {
	automaton := NewAhoCorasickBuilder().SetAsciiCaseInsensitive(true).Build([]string{"FOO", "bAr", "BaZ"})
	defer automaton.Close()
	fmt.Println(len(automaton.FindAll("foo bar baz")))
	// Output: 3
}

func ExampleAhoCorasickBuilder_SetMatchKind_standard_semantics() {
	haystack := "abcd"
	automaton := NewAhoCorasickBuilder().SetMatchKind(matchkind.Standard).Build([]string{"b", "abc", "abcd"})
	defer automaton.Close()
	match := automaton.FindFirst(haystack)
	fmt.Println(haystack[match.Start:match.End])
	// Output: b
}

func ExampleAhoCorasickBuilder_SetMatchKind_leftmost_first() {
	haystack := "abcd"
	automaton := NewAhoCorasickBuilder().SetMatchKind(matchkind.LeftMostFirst).Build([]string{"b", "abc", "abcd"})
	defer automaton.Close()
	match := automaton.FindFirst(haystack)
	fmt.Println(haystack[match.Start:match.End])
	// Output: abc
}

func ExampleAhoCorasickBuilder_SetMatchKind_leftmost_longest() {
	haystack := "abcd"
	automaton := NewAhoCorasickBuilder().SetMatchKind(matchkind.LeftMostLongest).Build([]string{"b", "abc", "abcd"})
	defer automaton.Close()
	match := automaton.FindFirst(haystack)
	fmt.Println(haystack[match.Start:match.End])
	// Output: abcd
}

func TestNewAhoCorasickBuilder(t *testing.T) {
	Convey("Given a new AhoCorasickBuilder", t, func() {
		builder := NewAhoCorasickBuilder()

		Convey("Then the builder is not nil", func() {
			So(builder, ShouldNotBeNil)
		})

		Convey("When the builder is used to build an AhoCorasick", func() {
			automaton := builder.Build([]string{"foo", "bar"})
			defer automaton.Close()

			Convey("Then the automaton is not nil", func() {
				So(automaton, ShouldNotBeNil)
			})
		})
	})
}
