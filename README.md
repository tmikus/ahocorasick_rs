# ahocorasick_rs

A Go wrapper for the Rust library [aho-corasick](https://github.com/BurntSushi/aho-corasick).

A library for finding occurrences of many patterns at once with SIMD acceleration in some cases.
This library provides multiple pattern search principally through an implementation of
the [Aho-Corasick algorithm](https://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_algorithm),
which builds a finite state machine for executing searches in linear time. Features include case insensitive matching,
overlapping matches, fast searching via SIMD and optional full DFA construction and search & replace in streams.

Dual-licensed under MIT or the [UNLICENSE](https://unlicense.org/).

## Documentation

https://pkg.go.dev/github.com/tmikus/ahocorasick_rs

## Installation

To build this package, you will need to have Rust installed. The minimum supported version of Rust is 1.60.0.
You can install Rust by following the instructions at https://www.rust-lang.org/tools/install.

Once Rust is installed, you can install this package with:
```bash
# Build the FFI bindings
git clone git@github.com:tmikus/aho-corasick-ffi.git
cargo build --release --manifest-path aho-corasick-ffi/Cargo.toml

# Configure env variables for the Go build. This is necessary so that the Go linker can find the Rust library.
export CGO_LDFLAGS="-L$(pwd)/aho-corasick-ffi/target/release"
export LD_LIBRARY_PATH="$(pwd)/aho-corasick-ffi/target/release"

# Install the Go package
go get -t github.com/tmikus/ahocorasick_rs

# Optional: Run the tests
go test github.com/tmikus/ahocorasick_rs
```

## Example: basic searching

This example shows how to search for occurrences of multiple patterns simultaneously. Each match includes the pattern
that matched along with the byte offsets of the match.

```go
package main

import (
    "fmt"
    "github.com/tmikus/ahocorasick_rs"
)

func main() {
    patterns := []string{"apple", "maple", "Snapple"}
    haystack := "Nobody likes maple in their apple flavored Snapple."
    ac := ahocorasick_rs.NewAhoCorasick(patterns)
    defer ac.Close() // Close the AhoCorasick instance when done.
    for _, match := range ac.FindAll(haystack) {
        fmt.Println(match.PatternIndex, match.Start, match.End)
    }
    // Output: 
    // 1, 13, 18 
    // 0, 28, 33 
    // 2, 43, 50
}

```

## Example: ASCII case insensitivity

This is like the previous example, but matches `Snapple` case insensitively using `AhoCorasickBuilder`:

```go
package main

import (
    "fmt"
    "github.com/tmikus/ahocorasick_rs"
)

func main() {
    patterns := []string{"apple", "maple", "snapple"}
    haystack := "Nobody likes maple in their apple flavored Snapple."
    ac := ahocorasick_rs.NewAhoCorasickBuilder().SetAsciiCaseInsensitive(true).Build(patterns)
    defer ac.Close() // Close the AhoCorasick instance when done.
    for _, match := range ac.FindAll(haystack) {
        fmt.Println(match.PatternIndex, match.Start, match.End)
    }
    // Output: 
    // 1, 13, 18 
    // 0, 28, 33 
    // 2, 43, 50
}
```

## Example: finding the leftmost first match

In the textbook description of Aho-Corasick, its formulation is typically
structured such that it reports all possible matches, even when they overlap
with another. In many cases, overlapping matches may not be desired, such as
the case of finding all successive non-overlapping matches like you might with
a standard regular expression.

Unfortunately the "obvious" way to modify the Aho-Corasick algorithm to do
this doesn't always work in the expected way, since it will report matches as
soon as they are seen. For example, consider matching the regex `Samwise|Sam`
against the text `Samwise`. Most regex engines (that are Perl-like, or
non-POSIX) will report `Samwise` as a match, but the standard Aho-Corasick
algorithm modified for reporting non-overlapping matches will report `Sam`.

A novel contribution of this library is the ability to change the match
semantics of Aho-Corasick (without additional search time overhead) such that
`Samwise` is reported instead. For example, here's the standard approach:

```go
package main

import (
    "fmt"
    "github.com/tmikus/ahocorasick_rs"
)

func main() {
    patterns := []string{"Samwise", "Sam"}
    haystack := "Samwise"
    ac := ahocorasick_rs.NewAhoCorasick(patterns)
    defer ac.Close() // Close the AhoCorasick instance when done.
    match := ac.FindFirst(haystack)
    fmt.Println(haystack[match.Start:match.End])
    // Output: 
    // Sam
}
```

And now here's the leftmost-first version, which matches how a Perl-like
regex will work:

```go
package main

import (
    "fmt"
    "github.com/tmikus/ahocorasick_rs"
    "github.com/tmikus/ahocorasick_rs/matchkind"
)

func main() {
    patterns := []string{"Samwise", "Sam"}
    haystack := "Samwise"
    ac := ahocorasick_rs.NewAhoCorasickBuilder().SetMatchKind(matchkind.LeftMostFirst).Build(patterns)
    defer ac.Close() // Close the AhoCorasick instance when done.
    match := ac.FindFirst(haystack)
    fmt.Println(haystack[match.Start:match.End])
    // Output: 
    // Samwise
}
```

In addition to leftmost-first semantics, this library also supports
leftmost-longest semantics, which match the POSIX behavior of a regular
expression alternation. See `MatchKind` in the docs for more details.