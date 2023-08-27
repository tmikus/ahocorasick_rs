package ahocorasick

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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
