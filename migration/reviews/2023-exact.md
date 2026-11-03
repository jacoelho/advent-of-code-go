# 2023 arithmetic, hiking and reader corrections

This batch addresses F1, F2, F4, F5 and F7 from the [2023 comparison](2023.md). Days 8/20 synchronization and day 21's garden specialization remain separate work. The useful reference ideas were compared against CameronAavik's C# solvers at `1f1313a67b959e21509300a54925ad0ec7a2a23b`; their input-specific topology and arithmetic assumptions were not adopted.

## Races and workflows

Day 6 finds the first winning hold in the increasing half of the race, then counts its symmetric interval. The strict comparison `hold > distance/(time-hold)` avoids both floating-point rounding and multiplication overflow. The shared error-aware binary search remains the only search implementation. Parsing validates matching time/distance fields; combined fields and the part-one product must fit in int64. A later zero count correctly cancels an overflowing prefix product.

Day 19 owns four-element ratings and ranges in private workflow types. Each split classifies none/all/partial matches before changing an endpoint, so it cannot widen the incoming range or overflow at extreme thresholds. Parsing validates destinations, final default rules, and complete ratings in 1..4000. Fixed-rating evaluation detects revisited workflows; range evaluation rejects reachable cycles while ignoring empty branches. Limits are 4 MiB, 1024 workflows, 16384 rules, 100000 parts and one million rule visits. The count is bounded by 4000⁴.

Independent regressions include the 2⁵⁴ zero-distance race (18014398509481983 winning holds), strict ties, zero-result product cancellation, nested workflow splits (576000000000), extreme thresholds, and a 256-rating enumeration independent of range splitting. Existing sample and personal answers remain unchanged.

## Exact hailstones and shared rational arithmetic

The [part-one puzzle](https://adventofcode.com/2023/day/24) asks about future XY paths within inclusive bounds. Determinants and rational comparisons now handle vertical paths, exact boundaries and independent nonnegative times. The explicit interpretation for degenerate projections is intersection of the two future path sets: stationary points and overlapping collinear rays count once when their intersection meets the square. This interpretation is not a claimed published degenerate example.

The [reproduced part-two statement](https://github.com/AntonKosov/advent-of-code-2023/blob/61cdf160b4dd/day24/part2/README.txt) requires integer rock position and velocity; collision times need not be integers. Its text was verified at blob `6dea521b55f6cf7152da4fe7a862342ee4114680`. No implementation from that repository was used.

For rock position r and velocity u, subtracting the anchor stone's cross-product equation from each other stone gives:

```text
r × (vi − v0) + (pi − p0) × u = pi × vi − p0 × v0
```

These anchor differences span every pair difference. All rows enter the shared RREF solver. Full rank supplies one candidate. With one free component, substituting the parametric solution into the anchor's original cross product yields a quadratic; its integer roots provide the complete finite candidate set. The free parameter is itself a rock component, so integrality is necessary. Every candidate must have six integer components and collide with every stone at a common nonnegative time in all three axes. Two valid candidates produce an ambiguity error; none produces an infeasibility error.

Rank at most four, or a rank-five family whose anchor polynomial vanishes identically, returns an explicit unsupported-family error. The implementation does not claim these cases are impossible or necessarily ambiguous. Input is bounded at 4 MiB and 2000 stones, with the scanner's existing per-line bound. Coordinates enter as checked int64 values, then all differences, products, elimination and output sums use arbitrary precision.

`pkg/matrix.Rat` now owns immutable `math/big.Rat` values instead of unchecked int64 fractions. Its zero value is zero. Constructors and numerator/denominator accessors copy mutable integers; arithmetic never changes either operand. Matrix clones can therefore safely share rational values. `Cmp` is the value comparison; direct equality is disabled. `ScaleToIntegers` computes the common denominator exactly and reports conversion overflow. The existing 2025 day 10 caller propagates those errors and infeasible joltage systems with machine context instead of adding -1 to the total. That change does not close the other 2025 optimization findings.

The bounded local quadratic calculation was retained instead of adding a second linear solver or a general symbolic framework. The advisor checked the equation signs, candidate validation, immutable helper contract and rank-five extension. QA independently verified the literal rock fixtures using rational equations before tests were added.

Regression evidence includes the published sample (2/47), a rock with velocity outside ±500 (position sum 60), duplicate initial rows, contradictory later stones, fractional positive collision time, past-only collisions, fractional rock components despite an integral sum, and rank-five unique/ambiguous/no-future cases. Small XY cases cover stationary, parallel, collinear and boundary behavior. Matrix tests cover overflow intermediates, zero values, alias isolation, dependent/inconsistent systems and checked scaling. Personal answers remain 26657 and 828418331313365.

## Hiking and input boundaries

Day 23 uses one corridor graph for both slope conditions. Junctions are all cells of underlying degree other than two, plus start/end. Walking stops on membership of the current cell in that set, preserving adjacent junctions, degree-two endpoints and parallel corridors. Slopes constrain departure; they do not prohibit entering from the side. Corridor interiors have degree two, so a simple junction path expands to a simple cell path. One owned visited slice removes the former 64-node restriction.

Graph construction is deterministic. Limits are one million stored cells/rows, 4096 junctions and 50 million DFS states. Exceeding a bound returns a distinct error, never a partial longest distance. The reference's fixed junction-grid optimization was rejected because its topology is not part of this solver's contract.

QA verified literal tiny mazes independently. A separate cell-level oracle checks 256 wall/slope cases without using production graph helpers. A 70-node chain exercises the old bit-width defect; an open grid exercises the junction bound. Published and personal results remain 94/154 and 2306/6718.

Days 3, 22 and 23 now propagate scanner failures. Day 15 reads one bounded payload up to 4 MiB, accepting an optional final LF/CRLF and rejecting extra lines. It preserves part one's opaque-text hashing; part two uses the checked number helper and validates focusing-power arithmetic. A single reader-error table covers both parts of days 3/15/22/23/24 with immediate and valid-prefix failures, preserving the original cause through `errors.Is`. A payload longer than 64 KiB verifies day 15 no longer silently loses its input.

## Verification and review

Focused tests passed before the combined uncached run:

```sh
go test -race -shuffle=on -count=1 -timeout=3m ./pkg/matrix ./internal/aoc2023 ./internal/aoc2025
staticcheck ./pkg/matrix ./internal/aoc2023 ./internal/aoc2025
```

The combined run passed (matrix 1.193s, 2023 12.578s, 2025 5.432s), using Go 1.27.3 and Staticcheck 2026.2.2. These are test durations, not benchmark claims. Temporary writable build/module/checker caches were used. Gopls tools were unavailable; scoped caller inspection, compilation, tests and Staticcheck supplied verification. The critic reviewed production, callers, tests, this record and the integration log, finding no material issues in the batch. Repository-wide checks remain scheduled after the remaining corrections.
