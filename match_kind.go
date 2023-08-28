package ahocorasick

// MatchKind is a knob for controlling the match semantics of an Aho-Corasick automaton.
//
// There are two generally different ways that Aho-Corasick automatons can report matches. The first way is the
// “standard” approach that results from implementing most textbook explanations of Aho-Corasick. The second way is to
// report only the leftmost non-overlapping matches. The leftmost approach is in turn split into two different ways of
// resolving ambiguous matches: leftmost-first and leftmost-longest.
//
// The MatchKindStandard match kind is the default and is the only one that supports overlapping matches and stream searching.
// (Trying to find overlapping or streaming matches using leftmost match semantics will result in an error in fallible APIs and a panic when using infallibe APIs.)
// The MatchKindStandard match kind will report matches as they are seen. When searching for overlapping matches, then all possible matches are reported.
// When searching for non-overlapping matches, the first match seen is reported. For example, for non-overlapping matches,
// given the patterns abcd and b and the haystack abcdef, only a match for b is reported since it is detected first.
// The abcd match is never reported since it overlaps with the b match.
//
// In contrast, the leftmost match kind always prefers the leftmost match among all possible matches. Given the same
// example as above with abcd and b as patterns and abcdef as the haystack, the leftmost match is abcd since it begins
// before the b match, even though the b match is detected before the abcd match. In this case, the b match is not
// reported at all since it overlaps with the abcd match.
//
// The difference between leftmost-first and leftmost-longest is in how they resolve ambiguous matches when there are
// multiple leftmost matches to choose from. Leftmost-first always chooses the pattern that was provided earliest,
// whereas leftmost-longest always chooses the longest matching pattern. For example, given the patterns a and ab and
// the subject string ab, the leftmost-first match is a but the leftmost-longest match is ab. Conversely, if the patterns
// were given in reverse order, i.e., ab and a, then both the leftmost-first and leftmost-longest matches would be ab.
// Stated differently, the leftmost-first match depends on the order in which the patterns were given to the Aho-Corasick automaton.
// Because of that, when leftmost-first matching is used, if a pattern A that appears before a pattern B is a prefix of B,
// then it is impossible to ever observe a match of B.
//
// If you’re not sure which match kind to pick, then stick with the standard kind, which is the default. In particular,
// if you need overlapping or streaming matches, then you must use the standard kind. The leftmost kinds are useful in
// specific circumstances. For example, leftmost-first can be very useful as a way to implement match priority based
// on the order of patterns given and leftmost-longest can be useful for dictionary searching such that only the longest
// matching words are reported.
//
// # Relationship with regular expression alternations
//
// Understanding match semantics can be a little tricky, and one easy way to conceptualize non-overlapping matches
// from an Aho-Corasick automaton is to think about them as a simple alternation of literals in a regular expression.
// For example, let’s say we wanted to match the strings Sam and Samwise, which would turn into the regex Sam|Samwise.
// It turns out that regular expression engines have two different ways of matching this alternation. The first way,
// leftmost-longest, is commonly found in POSIX compatible implementations of regular expressions (such as grep).
// The second way, leftmost-first, is commonly found in backtracking implementations such as Perl. (Some regex engines,
// such as RE2 and Rust’s regex engine do not use backtracking, but still implement leftmost-first semantics in an effort
// to match the behavior of dominant backtracking regex engines such as those found in Perl, Ruby, Python, Javascript and PHP.)
//
// That is, when matching Sam|Samwise against Samwise, a POSIX regex will match Samwise because it is the longest
// possible match, but a Perl-like regex will match Sam since it appears earlier in the alternation.
// Indeed, the regex Sam|Samwise in a Perl-like regex engine will never match Samwise since Sam will always have
// higher priority. Conversely, matching the regex Samwise|Sam against Samwise will lead to a match of Samwise in both
// POSIX and Perl-like regexes since Samwise is still longest match, but it also appears earlier than Sam.
//
// The “standard” match semantics of Aho-Corasick generally don’t correspond to the match semantics of any large
// group of regex implementations, so there’s no direct analogy that can be made here. MatchKindStandard match semantics are
// generally useful for overlapping matches, or if you just want to see matches as they are detected.
//
// The main conclusion to draw from this section is that the match semantics can be tweaked to precisely match either
// Perl-like regex alternations or POSIX regex alternations.
type MatchKind int

const (
	// Use standard match semantics, which support overlapping matches. When used with non-overlapping matches, matches are reported as they are seen.
	MatchKindStandard MatchKind = 1
	// Use leftmost-longest match semantics, which reports leftmost matches. When there are multiple possible leftmost matches, the longest match is chosen.
	//
	//This does not support overlapping matches or stream searching. If this match kind is used, attempting to find overlapping matches or stream matches will fail.
	MatchKindLeftMostLongest MatchKind = 2
	// Use leftmost-first match semantics, which reports leftmost matches. When there are multiple possible leftmost matches, the match corresponding to the pattern that appeared earlier when constructing the automaton is reported.
	//
	//This does not support overlapping matches or stream searching. If this match kind is used, attempting to find overlapping matches or stream matches will fail.
	MatchKindLeftMostFirst MatchKind = 3
)
