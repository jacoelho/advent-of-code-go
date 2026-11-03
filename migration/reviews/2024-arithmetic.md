# 2024 calibration and claw arithmetic

This batch resolves F6 and F8 in the [2024–2025 review](2024-2025.md), plus the F11 error boundaries within days 7 and 13. The remaining F11 boundaries were subsequently completed in the [reader and route corrections](2024-2025-boundaries.md) and the linked circuit/search records.

## Calibration reachability

Each input row is an independent equation, so repeated targets are retained and successful rows all contribute. A memoized backward predicate owns both parts' operator rules. Its state is the operand index and required value. The base case compares the first operand directly; addition, positive multiplication and decimal concatenation are inverted exactly. Zero has one decimal digit. Parsed values are numerically normalized, including leading zeros.

Multiplication by zero is the important exception: a required zero is reachable from any nonempty mathematical prefix. That prefix always has an expression, even if its intermediate value exceeds int64. Accepting this inverse branch preserves overshoot-then-reset behavior without storing huge forward accumulators. All nonzero inverse states stay at most the original target. Thus the solver handles mathematical nonnegative intermediates while constraining parsed literals and the exported total to int64.

The [puzzle statement](https://adventofcode.com/2024/day/7) establishes independent equations and left-to-right evaluation. The collected zelark solver at `432c59f6ccc52aae7da46cdd8ca0b3ffa8f8aede` retains rows, but its complete operator enumeration is unnecessary. The old A* heuristic, zero goal marker and suffix handling were removed. Parsing reuses the strict number helper and preserves error causes. Bounds are 10000 rows, 64 operands per row, the scanner's line limit and two million state expansions shared across the solve. Recursion removes an operand on every edge.

QA independently recomputed personal results using forward Python accumulator sets: 945512582195/271691107779347. Tests cover repeated rows (20/20), zero reset (10/10), unreachable zero-containing equations, concatenation ending in zero (0/120), and an overflowing mathematical prefix followed by a reset (10/10). A separate forward big.Int enumerator checks all operand tuples of length 1..4 over values 0..3 against targets 0..50 in both modes. Numeric error causes and total-overflow boundaries are also tested.

## Exact claw optimization

One solver owns both parts' feasibility and token cost. Press counts must be integral and nonnegative; part one additionally requires each count to be at most 100. Part two adds 10000000000000 to each prize coordinate without that ceiling. These constraints follow the [button costs and press limit](https://adventofcode.com/2024/day/13). Unwinnable machines contribute zero.

For independent button vectors, exact Cramer numerators must divide the determinant, then the common bounds apply. For dependent vectors, cross products first establish that the prize lies on the button line. A nonzero coordinate reduces the problem to `x*a+y*b=target`. Zero buttons and vertical movement have explicit cases. Otherwise extended GCD produces one solution and the complete family `(a+k*y/g,b-k*x/g)`. Nonnegative and optional upper bounds intersect into a closed integer interval for k. The slope of `3*a+b` selects the minimizing endpoint.

All coordinates enter as checked nonnegative int64 values; subsequent offsets, products, divisions, press counts, costs and totals use `math/big`. Parsed values are never mutated by calculations. Input is bounded at 4 MiB and 10000 complete three-line machines, with normal final-newline and CRLF handling. The same pinned reference's rational formula assumes independent vectors and misses the required press bounds, so it was not adopted as the complete solver.

QA independently recomputed the actual answers with exact Python determinants: 28138/108394825772874. Regressions cover negative and excessive presses, the inclusive 100-press boundary, dependent minima favoring either button, a bounded dependent optimum, gcd incompatibility, an off-line prize, zero/vertical buttons, a determinant above MaxInt64 with a small valid answer, and the part-two offset without a ceiling. Existing published results remain 480/875318608908. Shared reader tests preserve immediate and valid-prefix error causes for both days.

## Verification and benchmark harness

All new tests received QA assessment before edits. The focused uncached race/shuffle batch passed in 1.440s; the full 2024 package passed in 8.938s alongside the pulse integration. Staticcheck and formatting/whitespace checks passed. Verification used Go 1.27.3, Staticcheck 2026.2.2 and temporary writable caches. Gopls was unavailable; scoped caller inspection, compilation and runtime checks supplied verification.

The existing `Benchmark_day07` reused an exhausted file reader. After QA assessment, the named Go performance agent changed it to load the unchanged fixture once, create a fresh reader per iteration and check both errors and the independently verified 271691107779347 result. Exactly one isolated smoke invocation passed: `go test -pgo=off -run '^$' -bench '^Benchmark_day07$' -benchtime=1x -count=1 -cpu=1 ./internal/aoc2024`. No optimization or before/after performance claim is part of this batch.

The critic reviewed production, callers, tests, the corrected benchmark and this record, finding no material issues. Repository-wide checks and final style reconciliation remain after the other corrections.
