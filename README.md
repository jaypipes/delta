[![Go Reference](https://pkg.go.dev/badge/github.com/jaypipes/delta.svg)](https://pkg.go.dev/github.com/jaypipes/delta)
[![Test Status](https://github.com/jaypipes/delta/actions/workflows/test-unit.yml/badge.svg?branch=main)](https://github.com/jaypipes/delta/actions)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

`delta` is a Go library that facilitates efficient type inspection, comparison
and diffing of things.

# Type inspection, casting and conversion

Package `types` contains functions for type inspection. Subpackages
`types/cast` and `types/convert` contain functions for type casting and type
conversion, respectively.

# Comparisons

Package `cmp` contains functions for determining if two things are comparable,
and comparing those things for equality and sorting.

# Differences

Package `diff` contains functions for collecting differences between two
things.

# Developers

[Contributions](CONTRIBUTING.md) to `delta` are welcomed! Fork the repo on GitHub
and submit a pull request with your proposed changes. Or, feel free to log an
issue for a feature request or bug report.
