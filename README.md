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

### BobuSumisu's benchmark

The benchmark below was executed on Macbook Pro M1 Max using https://github.com/BobuSumisu/aho-corasick-benchmark.
Please note, that the benchmark below sends a rather large input to the library, which might give an unfair advantage to my library.
This is because this library uses CGO, which means that the input is copied from Go to Rust and back.
This is not the case for the other libraries, and they all run natively in Go without paying the penalty of crossing the FFI boundary.

```
          name    patterns        build     search    matches       alloc
    anknown           1000       0.70ms     2.33ms        407     0.06GiB
    bobusumisu        1000       0.69ms     0.52ms        407     0.07GiB
    cloudflare        1000       7.91ms     0.23ms          9     0.12GiB
    iohub             1000       0.35ms     0.39ms        407     0.12GiB
    tmikus            1000       1.30ms     0.24ms        388     0.12GiB

    anknown           2000       1.37ms     2.19ms        413     0.12GiB
    bobusumisu        2000       1.56ms     0.52ms        413     0.13GiB
    cloudflare        2000      14.55ms     0.22ms         13     0.23GiB
    iohub             2000       0.70ms     0.40ms        413     0.23GiB
    tmikus            2000       1.77ms     0.21ms        388     0.23GiB

    anknown           4000       2.73ms     2.25ms       1429     0.24GiB
    bobusumisu        4000       2.92ms     0.56ms       1429     0.26GiB
    cloudflare        4000      24.54ms     0.26ms         45     0.45GiB
    iohub             4000       1.36ms     0.46ms       1429     0.45GiB
    tmikus            4000       3.00ms     0.23ms        972     0.45GiB

    anknown           8000       5.83ms     2.34ms       3485     0.46GiB
    bobusumisu        8000       9.72ms     0.64ms       3485     0.50GiB
    cloudflare        8000      44.93ms     0.30ms         86     0.89GiB
    iohub             8000       2.95ms     0.59ms       3485     0.89GiB
    tmikus            8000       5.57ms     0.27ms       2303     0.89GiB

    anknown          16000      14.93ms     3.53ms       7977     0.92GiB
    bobusumisu       16000      13.26ms     0.88ms       7977     0.98GiB
    cloudflare       16000      87.66ms     0.34ms        173     1.77GiB
    iohub            16000       5.78ms     0.93ms       7977     1.78GiB
    tmikus           16000      10.88ms     0.34ms       5203     1.78GiB

    anknown          32000      26.44ms     2.61ms      10025     1.83GiB
    bobusumisu       32000      24.65ms     0.97ms      10025     1.97GiB
    cloudflare       32000     167.58ms     0.40ms        262     3.50GiB
    iohub            32000      11.84ms     1.29ms      10025     3.52GiB
    tmikus           32000      21.43ms     0.40ms       6280     3.52GiB

    anknown          64000      59.53ms     3.38ms      12505     3.61GiB
    bobusumisu       64000      51.98ms     1.08ms      12505     3.88GiB
    cloudflare       64000     337.31ms     0.52ms        526     6.85GiB
    iohub            64000      25.74ms     1.29ms      12505     6.89GiB
    tmikus           64000      42.25ms     0.49ms       7165     6.89GiB

    anknown         128000     153.10ms     5.02ms      39334     7.09GiB
    bobusumisu      128000     119.86ms     2.44ms      39334     7.61GiB
    cloudflare      128000     742.37ms     1.22ms       1141    13.46GiB
    iohub           128000      53.16ms     3.83ms      39300    13.54GiB
    tmikus          128000      92.45ms     1.07ms      21913    13.54GiB

    anknown         256000     333.36ms     8.87ms      59391    13.94GiB
    bobusumisu      256000     211.83ms     3.83ms      59391    14.99GiB
    cloudflare      256000    2116.44ms     1.95ms       2243    26.84GiB
    iohub           256000     117.91ms     5.95ms      58923    26.98GiB
    tmikus          256000     199.62ms     1.45ms      29451    27.00GiB

    anknown         512000     891.87ms     9.07ms      94000    27.91GiB
    bobusumisu      512000     503.17ms     6.15ms      94000    29.99GiB
    cloudflare      512000    3571.78ms     3.12ms       4490    53.68GiB
    iohub           512000     253.09ms    10.73ms      91986    53.97GiB
    tmikus          512000     467.33ms     2.17ms      38269    54.00GiB
```

### My benchmark

To better showcase the differences between these implementations I prepared a modified version of the benchmark,
which compares the performance of the libraries against different length of input string.
You can find the source code of the benchmark at https://github.com/tmikus/aho-corasick-benchmark


```
          name    patterns    input len       build     search    matches      alloc
    anknown         128985            6    188.30ms     0.01ms          7    0.21GiB
    bobusumisu      128985            6    173.07ms     0.00ms          7    0.58GiB
    cloudflare      128985            6    734.51ms     0.01ms          7    3.90GiB
    iohub           128985            6     65.75ms     0.01ms          7    0.07GiB
    tmikus          128985            6    100.24ms     0.00ms          2    0.01GiB

    anknown         128985           19    185.35ms     0.01ms         24    0.21GiB
    bobusumisu      128985           19    191.83ms     0.01ms         24    0.58GiB
    cloudflare      128985           19    729.61ms     0.01ms         23    3.90GiB
    iohub           128985           19     63.90ms     0.01ms         17    0.07GiB
    tmikus          128985           19     93.11ms     0.01ms          7    0.01GiB

    anknown         128985           41    181.05ms     0.02ms         49    0.21GiB
    bobusumisu      128985           41    185.44ms     0.01ms         49    0.58GiB
    cloudflare      128985           41    716.34ms     0.02ms         47    3.90GiB
    iohub           128985           41     63.01ms     0.02ms         45    0.07GiB
    tmikus          128985           41     94.33ms     0.01ms         11    0.01GiB

    anknown         128985           73    187.96ms     0.02ms         76    0.21GiB
    bobusumisu      128985           73    179.61ms     0.02ms         76    0.58GiB
    cloudflare      128985           73    737.76ms     0.03ms         67    3.90GiB
    iohub           128985           73     65.48ms     0.03ms         73    0.07GiB
    tmikus          128985           73     91.60ms     0.10ms         22    0.01GiB

    anknown         128985          146    186.37ms     0.03ms        153    0.21GiB
    bobusumisu      128985          146    181.59ms     0.03ms        153    0.58GiB
    cloudflare      128985          146    734.40ms     0.03ms        136    3.90GiB
    iohub           128985          146     66.32ms     0.03ms        146    0.07GiB
    tmikus          128985          146     92.33ms     0.01ms         43    0.01GiB

    anknown         128985          279    198.18ms     0.05ms        287    0.21GiB
    bobusumisu      128985          279    179.76ms     0.05ms        287    0.58GiB
    cloudflare      128985          279    729.76ms     0.04ms        234    3.90GiB
    iohub           128985          279     72.50ms     0.06ms        261    0.07GiB
    tmikus          128985          279     90.33ms     0.06ms         78    0.01GiB

    anknown         128985          534    188.12ms     0.08ms        495    0.21GiB
    bobusumisu      128985          534    178.65ms     0.08ms        495    0.58GiB
    cloudflare      128985          534    713.20ms     0.07ms        357    3.90GiB
    iohub           128985          534     63.94ms     0.07ms        445    0.07GiB
    tmikus          128985          534     88.42ms     0.10ms        137    0.01GiB

    anknown         128985         1118    189.47ms     0.14ms       1066    0.21GiB
    bobusumisu      128985         1118    181.79ms     0.15ms       1066    0.58GiB
    cloudflare      128985         1118    763.51ms     0.11ms        665    3.90GiB
    iohub           128985         1118     68.88ms     0.17ms        941    0.07GiB
    tmikus          128985         1118     91.90ms     0.09ms        299    0.01GiB

    anknown         128985         2233    177.17ms     0.27ms       2169    0.21GiB
    bobusumisu      128985         2233    180.50ms     0.28ms       2169    0.58GiB
    cloudflare      128985         2233    725.63ms     0.19ms       1225    3.90GiB
    iohub           128985         2233     64.07ms     0.24ms       1949    0.07GiB
    tmikus          128985         2233     98.35ms     0.24ms        617    0.01GiB

    anknown         128985         4428    180.39ms     0.51ms       4301    0.21GiB
    bobusumisu      128985         4428    181.76ms     0.48ms       4301    0.58GiB
    cloudflare      128985         4428    745.05ms     0.37ms       2055    3.90GiB
    iohub           128985         4428     64.06ms     0.44ms       3861    0.07GiB
    tmikus          128985         4428     94.88ms     0.22ms       1242    0.01GiB

    anknown         128985         8975    188.04ms     1.15ms       8718    0.21GiB
    bobusumisu      128985         8975    181.72ms     0.89ms       8718    0.58GiB
    cloudflare      128985         8975    758.63ms     0.69ms       3556    3.90GiB
    iohub           128985         8975     68.82ms     1.04ms       7784    0.07GiB
    tmikus          128985         8975     89.70ms     0.31ms       2518    0.01GiB

    anknown         128985        17837    181.85ms     1.88ms      17214    0.21GiB
    bobusumisu      128985        17837    180.63ms     1.73ms      17214    0.58GiB
    cloudflare      128985        17837    710.73ms     1.23ms       6141    3.90GiB
    iohub           128985        17837     70.92ms     2.11ms      15430    0.07GiB
    tmikus          128985        17837     89.48ms     0.33ms       4937    0.01GiB

    anknown         128985        35720    180.62ms     3.91ms      34232    0.22GiB
    bobusumisu      128985        35720    181.36ms     3.28ms      34232    0.58GiB
    cloudflare      128985        35720    752.59ms     2.45ms      10486    3.90GiB
    iohub           128985        35720     64.43ms     3.50ms      30701    0.08GiB
    tmikus          128985        35720     90.05ms     0.52ms       9877    0.01GiB

    anknown         128985        71332    185.42ms     7.15ms      68256    0.22GiB
    bobusumisu      128985        71332    189.95ms     6.84ms      68256    0.58GiB
    cloudflare      128985        71332    685.10ms     4.54ms      17431    3.90GiB
    iohub           128985        71332     61.38ms     6.37ms      61258    0.08GiB
    tmikus          128985        71332     95.16ms     0.92ms      19673    0.01GiB

    anknown         128985       142563    178.64ms    13.81ms     136922    0.22GiB
    bobusumisu      128985       142563    182.05ms    12.97ms     136922    0.59GiB
    cloudflare      128985       142563    712.75ms     8.64ms      28924    3.90GiB
    iohub           128985       142563     65.50ms    11.67ms     122868    0.09GiB
    tmikus          128985       142563     96.77ms     1.77ms      39564    0.01GiB

    anknown         128985       285779    178.54ms    24.49ms     274116    0.23GiB
    bobusumisu      128985       285779    177.26ms    26.34ms     274116    0.60GiB
    cloudflare      128985       285779    711.18ms    15.52ms      46529    3.90GiB
    iohub           128985       285779     66.47ms    24.92ms     246073    0.10GiB
    tmikus          128985       285779     93.05ms     2.95ms      79029    0.01GiB

    anknown         128985       571011    178.51ms    52.57ms     545790    0.26GiB
    bobusumisu      128985       571011    179.90ms    51.04ms     545790    0.63GiB
    cloudflare      128985       571011    735.84ms    28.38ms      70251    3.90GiB
    iohub           128985       571011     60.61ms    50.41ms     490157    0.13GiB
    tmikus          128985       571011     90.21ms     6.37ms     157663    0.02GiB
```

## Conclusion

Based on the results above, the answer to the question "Which library is the fastest?" is "It depends on the input length".

For short inputs (up to around 3000 characters) it doesn't really matter which library you use as they all perform similarly.

For longer inputs you're best off using my library, as it is the fastest.
