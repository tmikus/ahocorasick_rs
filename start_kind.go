package ahocorasick

// The kind of anchored starting configurations to support in an Aho-Corasick searcher.
//
// Depending on which searcher is used internally by AhoCorasick, supporting both unanchored and anchored searches can be quite costly. For this reason, AhoCorasickBuilder::start_kind can be used to configure whether your searcher supports unanchored, anchored or both kinds of searches.
//
// This searcher configuration knob works in concert with the search time configuration Input::anchored. Namely, if one requests an unsupported anchored mode, then the search will either panic or return an error, depending on whether you’re using infallible or fallibe APIs, respectively.
//
// AhoCorasick by default only supports unanchored searches.
type StartKind int

const (
	StartKindBoth       StartKind = 1 // Support both anchored and unanchored searches.
	StartKindUnanchored StartKind = 2 // Support only unanchored searches. Requesting an anchored search will return an error in fallible APIs and panic in infallible APIs.
	StartKindAnchored   StartKind = 3 // Support only anchored searches. Requesting an unanchored search will return an error in fallible APIs and panic in infallible APIs.
)
