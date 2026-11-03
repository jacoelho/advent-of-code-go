# 2025 Day 10 machine limits

Day 10 now reports infeasible machines, bounds exact search and proves its integer arithmetic safe before enumeration. This completes F13 in the [2024–2025 review](2024-2025.md). The required critic accepted the final code, proof and QA-approved regression evidence. Shared matrix code and input fixtures remain unchanged.

## Result and contracts

The subset light solver and exact RREF/free-variable joltage solver remain. Unreachable light configurations now return errors instead of `buttonCount+1`; inconsistent, fractional-only and negative-only joltage solutions return errors instead of contributing a sentinel to the total. Parse, arithmetic-range, aggregate-overflow and work-exhaustion errors reach the entrypoint with a one-based machine number. No partial optimum is returned after exhaustion.

The parser uses `convert.ScanNumber` for comma-separated indices/counts. It validates 1–64 lights, equal counter width, nonnegative int64 counters, unique in-range button indices, at most 24 buttons per machine, and at most 1024 machines per input. Masks are explicitly uint64, including bit 63. Empty buttons `()` and zero-button machines remain valid. Blank/malformed records return errors before indexing or shifting. The shared scanner retains its line-token bound and preserves reader errors.

The duplicate DFS leaf feasibility/summation loop is deleted: every candidate uses `validateAndSum`. The former free-variable wrapper and narrative comments were removed. No generic optimization framework, floating solver or shared helper changes were added.

## Exactness and arithmetic proof

Shared `matrix.RREF` uses arbitrary-precision rationals. Its extracted solution is

`x_i = (b_i + sum(c_ki * t_k)) / D`, with `D > 0`.

`ScaleToIntegers` already checks that D, each b and each c fit int64. Those checks alone do not cover downstream products or partial sums; the new admission step does.

A nonempty button cannot be pressed more than the minimum target of any counter it affects, because every button contribution is nonnegative. This gives an initial upper bound U for each free variable. A no-op button can be assigned zero in a minimum solution. For each negative coefficient c_ki, the necessary inequality is

`t_k <= floor((b_i + sum(max(c_ji,0) * U_j, j != k)) / -c_ki)`.

Limit tightening now evaluates that numerator, negation and division with `big.Int`; it cannot overflow even before admission. Other variables use their original physical upper bounds, so the tightening does not depend on iteration order and cannot remove a feasible solution. A negative numerator retains zero as a harmless candidate which the feasibility check rejects; it cannot conceal a valid nonnegative parameter. The final limits L remain nonnegative int64 values.

Before any int64 search arithmetic, admission computes exact nonnegative envelopes

`B_i = abs(b_i) + sum(abs(c_ki) * L_k)`.

Every B must fit int64. Therefore each product, complete expression and every partial sum used in pruning fits int64, regardless of cancellation or summation order. Admission also requires `sum(floor(B_i/D)) <= MaxInt64`, proving every candidate objective sum fits. The admission check is deliberately conservative: a feasible family can be refused as an unsupported arithmetic range if its envelope is too large. `[.] (0) (0) {9223372036854775807}` is the tested example; it returns a range error rather than claiming infeasibility.

DFS enumerates integer parameters from zero through L. Its partial feasibility check adds the maximum possible remaining positive contributions; if the result is negative, no continuation can make that component nonnegative. Accepted leaves require every component to be nonnegative and divisible by D, then minimize their sum. Each loop breaks at its upper bound before incrementing, including an int64 maximum bound. Thus pruning and enumeration retain the optimum for every admitted machine completed within the work budget.

Part-two aggregation separately checks `presses <= MaxInt64-total` before adding nonnegative machine costs. Part-one aggregation is bounded by 1024 machines times 24 buttons, so it cannot overflow. In the light model, pressing a button twice cancels and adds cost; some optimum uses each button at most once, making subset enumeration exact.

## Work and memory limits

Each entrypoint owns one 250,000,000-unit budget shared by all machines. It is never reset between machines. Every subset costs `buttonCount+1` units. Every DFS candidate/prune or leaf costs `1 + buttonCount*(freeCount+1)` units, covering its affine-term work. Matrix preprocessing is charged `1 + counterCount*buttonCount*(buttonCount+1)` units; its dimensions and integer input bit widths are independently bounded by the parser. Exact bound preparation is also finite under those same dimensions. These are declared work units, not elapsed-time measurements.

DFS depth is at most 24. Input count, masks and matrix dimensions are bounded. Resource exhaustion returns a distinct error; the solver does not silently impose a guessed answer, drop a machine or reset its budget.

## Independent evidence and reference

QA assessed the coherent boundary/arithmetic/search batch before tests were written. Existing official results 7/33 and personal results 545/22430 passed. The existing free-variable example remains 3, and existing aggregation-failure regressions remain in place.

New tests include unreachable no-op lights; zero buttons; unique half-integer and negative-only joltage solutions; strict parsing/width/button caps; successful bit-63 masks; MaxInt64 single-machine success; two-machine aggregate overflow; conservative arithmetic-range refusal; and the existing solver budget supplied with zero remaining work. Tests assert result/error categories, not exact charged operations.

The independent tiny oracle uses ordinary count assignments: all 64 triples of two-counter button masks, all four light targets, and all nine joltage targets in 0..2. It compares 256 light and 576 joltage results. Joltage expectations enumerate all 27 count triples in 0..2; no nonempty button can need more than the maximum target, and no-op buttons optimally use zero. Light expectations enumerate binary button counts and apply parity. The oracle does not call matrix, scaling, pruning or production feasibility helpers.

Compared the already selected `reference/.repos/NickyMeuleman--scrapyard/advent_of_code/2025/solutions/src/day_10.rs` at verified revision `1d6d8502daa017f4885c3b3110ea9ff2a3c62de1`. It uses BFS for lights and an external floating-point ILP solver for joltage, with an unwrapped solve result. Neither is needed to repair the existing smaller exact Go model. The reviewed F13 notes and shared exact-matrix contract support retaining RREF and finite free-variable enumeration.

Input SHA256: `76631ed00d473af2940118a406d3fff9444167bf3d1e3c9b7335ac773ed5c632`. Fixture bytes did not change. Go codebase/modern-Go skills were applied; gopls was unavailable, so scoped caller/source inspection and Go checks were used.

## Verification

All commands below passed with the agreed caches:

```sh
env GOCACHE=/tmp/aoc-performance-cache GOMODCACHE=/tmp/aoc-performance-modcache STATICCHECK_CACHE=/tmp/aoc-staticcheck-cache go test ./internal/aoc2025 -run '^(Test_day10p0[12]|TestJoltageFailureStopsAggregation)$' -count=1 -race -shuffle=on -timeout=3m

env GOCACHE=/tmp/aoc-performance-cache GOMODCACHE=/tmp/aoc-performance-modcache STATICCHECK_CACHE=/tmp/aoc-staticcheck-cache go test ./internal/aoc2025 -run '^(TestMachineFeasibility|TestMachineArithmeticLimits|TestMachineWorkLimit|TestMachineParsing|TestMachineSmallExhaustive|Test2025ReaderErrors)$' -count=1 -race -shuffle=on -timeout=3m

env GOCACHE=/tmp/aoc-performance-cache GOMODCACHE=/tmp/aoc-performance-modcache STATICCHECK_CACHE=/tmp/aoc-staticcheck-cache staticcheck ./internal/aoc2025

git diff --check -- internal/aoc2025/day10.go internal/aoc2025/day10_test.go
```

The shared reader table includes both Day 10 rows to the shared reader table using `[.] (0) {0}\n`. The focused command below also passed, preserving the original sentinel through both immediate and valid-prefix failures with `errors.Is`:

```sh
env GOCACHE=/tmp/aoc-performance-cache GOMODCACHE=/tmp/aoc-performance-modcache STATICCHECK_CACHE=/tmp/aoc-staticcheck-cache go test ./internal/aoc2025 -run '^Test2025ReaderErrors$/^day10p0[12]$' -count=1 -race -shuffle=on -timeout=3m
```

Both files were formatted with gofmt. No benchmarks or performance claims. The required critic found no material issues; final repository checks are recorded in the migration entry point.
