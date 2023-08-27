package ahocorasick

import (
	_ "embed"
	"testing"

	ac "github.com/BobuSumisu/aho-corasick"
)

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
