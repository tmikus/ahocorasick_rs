package ahocorasick

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func ExampleAhoCorasickBuilder_SetAsciiCaseInsensitive() {
	automaton := NewAhoCorasickBuilder().SetAsciiCaseInsensitive(true).Build([]string{"FOO", "bAr", "BaZ"})
	fmt.Println(len(automaton.Search("foo bar baz")))
	// Output: 3
}

func TestNewAhoCorasickBuilder(t *testing.T) {
	Convey("Given a new AhoCorasickBuilder", t, func() {
		builder := NewAhoCorasickBuilder()

		Convey("Then the builder is not nil", func() {
			So(builder, ShouldNotBeNil)
		})

		Convey("When the builder is used to build an AhoCorasick", func() {
			automaton := builder.Build([]string{"foo", "bar"})

			Convey("Then the automaton is not nil", func() {
				So(automaton, ShouldNotBeNil)
			})
		})
	})
}
