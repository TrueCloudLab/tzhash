# Demo

[![asciicast](https://asciinema.org/a/IArEDLTrQyabI3agSSpINoqNu.svg)](https://asciinema.org/a/IArEDLTrQyabI3agSSpINoqNu)

**In project root:**

```bash
$ make
...
$ ./demo.sh

```

# Homomorphic hashing in golang

Package `tz` contains pure-Go (with some Assembly) implementation of hashing
function described by [Tillich and
Zémor](https://link.springer.com/content/pdf/10.1007/3-540-48658-5_5.pdf).

There are [existing implementations](https://github.com/srijs/hwsl2-core)
already, however they are written in C.

Package `gf127` contains arithmetic in `GF(2^127)` with `x^127+x^63+1` as reduction polynomial.

# Description

TZ Hash can be used instead of Merkle-tree for data-validation, because
homomorphic hashes are concatenable: hash sum of data can be calculated based on
hashes of chunks.

The example of how it works can be seen in tests and demo.

# Benchmarks

## go vs AVX vs AVX2 version

```
BenchmarkSum/AVX_digest-8             308       3889484 ns/op          25.71 MB/s         5 allocs/op
BenchmarkSum/AVXInline_digest-8       457       2455437 ns/op          40.73 MB/s         5 allocs/op
BenchmarkSum/AVX2_digest-8            399       3031102 ns/op          32.99 MB/s         3 allocs/op
BenchmarkSum/AVX2Inline_digest-8      602       2077719 ns/op          48.13 MB/s         3 allocs/op
BenchmarkSum/PureGo_digest-8           68       17795480 ns/op          5.62 MB/s         5 allocs/op
```

# Contributing

Feel free to contribute to this project after reading the [contributing
guidelines](CONTRIBUTING.md).

Before starting to work on a certain topic, create a new issue first, describing
the feature/topic you are going to implement.

# Makefile

``` bash
  Usage:

    make <target>

  Targets:

    all       Just `make` will build all possible binaries
    clean     Print version
    dep       Pull go dependencies
    help      Show this help prompt
    test      Run Unit Test with go test
    version   Print version
```

# References

- [https://link.springer.com/content/pdf/10.1007/3-540-48658-5_5.pdf](https://link.springer.com/content/pdf/10.1007/3-540-48658-5_5.pdf)
- [https://github.com/srijs/hwsl2-core](https://github.com/srijs/hwsl2-core)
