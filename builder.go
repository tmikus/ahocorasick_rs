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
	"github.com/tmikus/ahocorasick_rs/ahocorasickkind"
	"github.com/tmikus/ahocorasick_rs/matchkind"
	"github.com/tmikus/ahocorasick_rs/startkind"
	"unsafe"
)

type AhoCorasickBuilder struct {
	asciiCaseInsensitive bool
	byteClasses          bool
	denseDepth           *uint
	kind                 *ahocorasickkind.AhoCorasickKind
	matchKind            matchkind.MatchKind
	prefilter            bool
	startKind            startkind.StartKind
}

func NewAhoCorasickBuilder() *AhoCorasickBuilder {
	return &AhoCorasickBuilder{
		asciiCaseInsensitive: false,
		byteClasses:          true,
		denseDepth:           nil,
		kind:                 nil,
		matchKind:            matchkind.Standard,
		prefilter:            true,
		startKind:            startkind.Unanchored,
	}
}

func (b *AhoCorasickBuilder) Build(patterns []string) *AhoCorasick {
	cPatterns := make([]*C.char, len(patterns))
	for i, pattern := range patterns {
		cPatterns[i] = C.CString(pattern)
		defer C.free(unsafe.Pointer(cPatterns[i]))
	}
	options := C.AhoCorasickBuilderOptions{
		ascii_case_insensitive: bool_to_c_int(b.asciiCaseInsensitive),
		byte_classes:           bool_to_c_int(b.byteClasses),
		dense_depth:            (*C.size_t)(unsafe.Pointer(b.denseDepth)),
		kind:                   (*C.size_t)(unsafe.Pointer(b.kind)),
		match_kind:             C.size_t(b.matchKind),
		prefilter:              bool_to_c_int(b.prefilter),
		start_kind:             C.size_t(b.startKind),
	}
	automaton := C.build_automaton(
		(**C.char)(&cPatterns[0]),
		C.size_t(len(patterns)),
		(*C.AhoCorasickBuilderOptions)(unsafe.Pointer(&options)),
	)
	return &AhoCorasick{
		automaton: automaton,
	}
}

func (b *AhoCorasickBuilder) SetAsciiCaseInsensitive(asciiCaseInsensitive bool) *AhoCorasickBuilder {
	b.asciiCaseInsensitive = asciiCaseInsensitive
	return b
}

func (b *AhoCorasickBuilder) SetByteClasses(byteClasses bool) *AhoCorasickBuilder {
	b.byteClasses = byteClasses
	return b
}

func (b *AhoCorasickBuilder) SetDenseDepth(denseDepth *uint) *AhoCorasickBuilder {
	b.denseDepth = denseDepth
	return b
}

func (b *AhoCorasickBuilder) SetKind(kind *ahocorasickkind.AhoCorasickKind) *AhoCorasickBuilder {
	b.kind = kind
	return b
}

func (b *AhoCorasickBuilder) SetMatchKind(matchKind matchkind.MatchKind) *AhoCorasickBuilder {
	b.matchKind = matchKind
	return b
}

func (b *AhoCorasickBuilder) SetPrefilter(prefilter bool) *AhoCorasickBuilder {
	b.prefilter = prefilter
	return b
}

func (b *AhoCorasickBuilder) SetStartKind(startKind startkind.StartKind) *AhoCorasickBuilder {
	b.startKind = startKind
	return b
}

func bool_to_c_int(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
