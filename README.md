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
    match := ac.FindFirst(haystack)
    fmt.Println(haystack[match.Start:match.End])
    // Output: 
    // Samwise
}
```

In addition to leftmost-first semantics, this library also supports
leftmost-longest semantics, which match the POSIX behavior of a regular
expression alternation. See `MatchKind` in the docs for more details.

## Benchmarks

The benchmark below was executed on Macbook Pro M1 Max using https://github.com/BobuSumisu/aho-corasick-benchmark.
Please note, that the benchmark below sends a rather large input to the library, which might give an unfair advantage to my library.
This is because this library uses CGO, which means that the input is copied from Go to Rust and back.
This is not the case for the other libraries, and they all run natively in Go without paying the penalty of crossing the FFI boundary.

```
          name    patterns        build    search    matches       alloc
    anknown           1000       0.75ms    2.34ms        407     0.06GiB
    bobusumisu        1000       0.72ms    0.55ms        407     0.07GiB
    cloudflare        1000       6.43ms    0.22ms          9     0.12GiB
    iohub             1000       0.37ms    0.41ms        407     0.12GiB
    tmikus            1000       1.50ms    0.23ms        388     0.12GiB

    anknown           2000       1.49ms    2.26ms        413     0.12GiB
    bobusumisu        2000       1.82ms    0.54ms        413     0.13GiB
    cloudflare        2000      18.82ms    0.27ms         13     0.23GiB
    iohub             2000       0.72ms    0.40ms        413     0.23GiB
    tmikus            2000       1.97ms    0.22ms        388     0.23GiB

    anknown           4000       3.11ms    2.31ms       1429     0.24GiB
    bobusumisu        4000       3.44ms    0.57ms       1429     0.26GiB
    cloudflare        4000      29.64ms    0.25ms         45     0.45GiB
    iohub             4000       1.39ms    0.48ms       1429     0.45GiB
    tmikus            4000       3.43ms    0.24ms        972     0.45GiB

    anknown           8000       6.73ms    2.44ms       3485     0.46GiB
    bobusumisu        8000      10.16ms    0.67ms       3485     0.50GiB
    cloudflare        8000      51.83ms    0.30ms         86     0.89GiB
    iohub             8000       2.77ms    0.63ms       3485     0.89GiB
    tmikus            8000       6.56ms    0.29ms       2303     0.89GiB

    anknown          16000      16.11ms    2.63ms       7977     0.92GiB
    bobusumisu       16000      19.34ms    0.86ms       7977     0.99GiB
    cloudflare       16000      98.52ms    0.35ms        173     1.77GiB
    iohub            16000       6.27ms    1.06ms       7977     1.78GiB
    tmikus           16000      12.71ms    0.36ms       5203     1.78GiB

    anknown          32000      28.66ms    2.68ms      10025     1.83GiB
    bobusumisu       32000      31.09ms    0.93ms      10025     1.97GiB
    cloudflare       32000     169.21ms    0.43ms        262     3.50GiB
    iohub            32000      12.86ms    1.27ms      10025     3.52GiB
    tmikus           32000      25.36ms    0.42ms       6280     3.52GiB

    anknown          64000      60.89ms    2.82ms      12505     3.62GiB
    bobusumisu       64000      48.42ms    1.12ms      12505     3.88GiB
    cloudflare       64000     333.61ms    0.53ms        526     6.85GiB
    iohub            64000      25.67ms    1.42ms      12505     6.89GiB
    tmikus           64000      49.33ms    0.43ms       7165     6.90GiB

    anknown         128000     147.33ms    4.54ms      39334     7.09GiB
    bobusumisu      128000     121.79ms    2.46ms      39334     7.62GiB
    cloudflare      128000     761.86ms    1.31ms       1141    13.46GiB
    iohub           128000      56.81ms    3.51ms      39300    13.54GiB
    tmikus          128000     104.42ms    1.22ms      21913    13.55GiB

    anknown         256000     341.93ms    6.23ms      59391    13.95GiB
    bobusumisu      256000     209.34ms    3.91ms      59391    14.99GiB
    cloudflare      256000    1645.95ms    2.03ms       2243    26.84GiB
    iohub           256000     122.38ms    5.90ms      58923    26.99GiB
    tmikus          256000     235.84ms    1.79ms      29451    27.01GiB

    anknown         512000     851.71ms    8.61ms      94000    27.93GiB
    bobusumisu      512000     419.63ms    6.07ms      94000    30.01GiB
    cloudflare      512000    4240.75ms    3.07ms       4490    53.69GiB
    iohub           512000     264.89ms    9.36ms      91986    53.98GiB
    tmikus          512000     526.98ms    2.02ms      38269    54.03GiB
```
