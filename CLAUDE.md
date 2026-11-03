# CLAUDE.md

## Project Overview

Advent of Code solutions in Go.

## Development Commands

### Testing
- `make test` - Run all tests with race detection and shuffling
- `make test-2024` - Run tests for specific year (e.g., 2024)
- `go test -race -shuffle=on -timeout=2m -v ./internal/aoc2024/...` - Run tests for specific package

### Code Quality
- `make fmt` - Format all Go code
- `make staticcheck` - Run static analysis

### Creating New Solutions
- `make tmpl-2024-01` - Generate template files for year 2024, day 01
  - Creates `internal/aoc<year>/day<day>.go` with solution stubs
  - Creates `internal/aoc<year>/day<day>_test.go` with test structure

## Architecture

### Solution Structure

Each day's solution follows a consistent pattern:
- Solutions are in `internal/aoc<year>/day<dd>.go`
- Tests are in `internal/aoc<year>/day<dd>_test.go`
- Each day has two functions: `day<dd>p01(r io.Reader) (string, error)` and `day<dd>p02(r io.Reader) (string, error)`
- Input files are stored in `inputs/<year>/<dd>.txt`

### Testing Pattern

Tests use the `aoc.AOCTest` helper with `aoc.TestInput` slices:
```go
tests := []aoc.TestInput{
    {
        Input: strings.NewReader(`example input`),
        Want: "expected output",
    },
    {
        Input: aoc.FileInput(t, 2024, 1),
        Want:  "actual answer",
    },
}
aoc.AOCTest(t, day01p01, tests)
```

### Package Structure

**Public Packages (`pkg/`):**
Reusable utilities available to all code:

- **Core Data Structures:**
  - `collections` - Deque, Heap, Set, Stack
  - `grid` - 2D grid operations with `Position2D[T]` generic type, includes Neighbours4/Neighbours8 iterators
  - `graph` - Graph algorithms (clique finding, etc.)

- **Algorithms:**
  - `search` - A*, BFS, DFS implementations
  - `xmath` - Math utilities (Abs, GCD, LCM, etc.)

- **Utilities:**
  - `xiter` - Custom iterator utilities, heavily used for functional-style operations
  - `xslices` - Extended slice operations
  - `xmaps` - Extended map operations
  - `xstrings` - String manipulation
  - `convert` - Type conversion helpers (e.g., `digits`)
  - `scanner` - Custom scanning with `scanner.NewScanner(r, parseFunc)` pattern
  - `funcs` - Higher-order functions (memoization, etc.)
  - `conc` - Concurrent data structures
  - `assert` - Test assertions

**Internal Packages (`internal/`):**
AOC-specific code not meant for external use:

- `aoc` - AOC helpers and test utilities
  - `aoc.Must[T](t T, err error) T` - Panic on error, useful for parsing
  - `aoc.MustAtoi(s string) int` - Panic on Atoi error
  - `aoc.FileInput(t, year, day)` - Load puzzle input file
- `aoc<year>` - Solutions for specific years (e.g., `aoc2024`)

### Key Patterns

1. **Iterators**: The codebase heavily uses Go 1.23+ `iter.Seq[T]` for functional-style operations
2. **Generics**: Internal packages leverage generics extensively (e.g., `Position2D[T constraints.Signed]`)
3. **Scanner pattern**: Use `scanner.NewScanner(r, parseFunc)` for custom line parsing with `.Values()` iterator
4. **Test-first**: Tests include both example inputs and actual puzzle inputs

## Development Guidelines

1. **Add helpers to existing packages** - When you need utility functions, add them to the appropriate package in `pkg/` rather than implementing inline or creating standalone utilities unless they will not be reused
2. **Keep code simple** - Prefer clear, straightforward implementations over clever solutions
3. **Minimal comments** - Only add comments to explain complex algorithms; let the code speak for itself
4. **Functional style** - Prefer functional patterns with `iter.Seq[T]`, range-over-func, and pipeline-style operations using `xiter` utilities
5. **Consistent solution structure** - Match the established pattern of `dayXXp01/p02` functions returning `(string, error)`
6. **Check for existing utilities before implementing** - Before writing manual loops for common operations, check if a `pkg/` package already provides the functionality:
   - `xslices.Map` - Transform slice elements
   - `xslices.Sum` - Sum numeric slices
   - `xslices.Filter` - Filter slice elements
   - `xiter` - Iterator operations and transformations
   - Review the relevant package docs when working with slices, maps, strings, or iterators
7. Use the problem domain language when naming

## Common Implementation Patterns

### Parsing Input

**Prefer `scanner.NewScanner` over `bufio.Scanner`:**
```go
// Good - uses scanner pattern with iterator
s := scanner.NewScanner(r, parseFunc)
for item := range s.Values() {
    // process item
}
if err := s.Err(); err != nil {
    return "", err
}

// Avoid - manual bufio.Scanner loops
scanner := bufio.NewScanner(r)
for scanner.Scan() {
    line := scanner.Text()
    // manual parsing
}
return scanner.Err()
```

The `scanner.NewScanner` pattern:
- Takes an `io.Reader` and a parse function `func([]byte) (T, error)`
- Returns a `*Scanner[T]` that provides an iterator via `.Values()`
- Enables functional-style pipelines with `xiter` utilities
- Centralizes error handling with `.Err()`

### Counting and Aggregation

**Use `xslices.CountFunc` and `xslices.Sum` instead of manual loops:**
```go
// Good - use helpers
count := xslices.CountFunc(predicate, slice)
total := xslices.Sum(numbers)

// Avoid - manual counting
count := 0
for _, v := range slice {
    if predicate(v) { count++ }
}

// Avoid - manual summation
total := 0
for _, v := range numbers {
    total += v
}
```

### Functional Pipelines with Iterators

**Use `xiter.Map` and `xiter.Sum` for transformations:**
```go
// Good - functional pipeline with iterators
total := xiter.Sum(xiter.Map(transform, iterator))

// Also good - collect to slice then sum
items := slices.Collect(iterator)
total := xslices.Sum(items)

// Avoid - manual accumulation
total := 0
for item := range iterator {
    total += transform(item)
}
```

### String to Number Conversion

**Use `aoc.MustAtoi` for cleaner code:**
```go
// Good - clear panic intent for invalid input (acceptable in AoC context)
num := aoc.MustAtoi(str)

// Avoid - ignoring errors silently
num, _ := strconv.Atoi(str)
```

### Bit Manipulation for Powers of 2

**Use bit shifting for clarity:**
```go
// Good - clear intent for power of 2
value := 1 << n  // equivalent to 2^n

// Avoid - manual doubling loop
value := 1
for i := 0; i < n; i++ {
    value *= 2
}
```