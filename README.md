# Advent of Code in Go

Solutions for 2015–2025, using Go 1.27 or later.

## Development

```sh
make test          # All tests, with race detection and shuffled order
make test-2024     # Tests for one year
make staticcheck  # Static analysis
make fmt          # Format Go source
```

Report uncached puzzle-test timings for a year:

```sh
go run ./cmd/test-timings -year 2024
```

Omit `-year` to run every year. Timings include sample and personal-input tests;
they are useful for locating slow solutions, not controlled benchmark comparisons.

Create a solution and its tests with `make tmpl-2024-01`. Fetch an input with
`make input-fetch-2024-01`, using the `AOC_SESSION_COOKIE` environment variable.
Both commands preserve existing files.
