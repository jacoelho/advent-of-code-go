# Go migration

Goal: consolidate all Advent of Code solutions for 2015–2025 in Go 1.27, improve idioms, investigate slow days, and compare the collected references for simpler algorithms. Coverage includes all 513 puzzle parts: 49 per edition through 2024 and 23 for 2025. The final day of each edition has one computational part.

[Coverage and source inventory](coverage.json) records each part and the exact upstream revisions inspected. Counts describe implementation presence, not completed correctness, performance, or simplicity reviews.

| Edition | Implemented parts | Remaining |
| --- | ---: | ---: |
| 2015 | 49 | 0 |
| 2016 | 49 | 0 |
| 2017 | 49 | 0 |
| 2018 | 49 | 0 |
| 2019 | 49 | 0 |
| 2020 | 49 | 0 |
| 2021 | 49 | 0 |
| 2022 | 49 | 0 |
| 2023 | 49 | 0 |
| 2024 | 49 | 0 |
| 2025 | 23 | 0 |
| Total | 513 | 0 |

## Sources

- [Original repository](https://github.com/jacoelho/advent-of-code/tree/6522adb786cf29d1322d193234698ed2c429b89b): F# 2020 and earlier Rust work.
- [Rust repository](https://github.com/jacoelho/advent-of-code-r/tree/c420b42f08c6e8687d19f7e446c04b3feab88b04): 2019, 2022, and 2023.
- [Clojure repository](https://github.com/jacoelho/advent-of-code-clj/tree/601804012ca87c06ff004daeb94bc3fd3064bf03): 2015–2019 and 2021.
- [Collected reference repositories](../reference/README.md): examples to compare after establishing the puzzle contract and independent expected answers.

Paths inside each source inventory are relative to that repository. The local directory named `advent-of-code` contains other people's reference repositories; it is not the original source above. The local Clojure `update-dev` branch has the same day files as the inspected current `main` revision.

## Completed batches

- Raised `go.mod` from 1.25 to 1.27; CI already selects its toolchain from that file. Verified with installed Go 1.27.3. [Release notes](https://go.dev/doc/go1.27).
- Ported both parts of 2020 days 1–6 with original inputs, official examples, independently calculated full-input answers, focused regressions, and propagated reader errors. See [oracle provenance](2020-days01-06.md).
- Finished 2020 days 7–19, completing all 49 parts for that year. Source bugs were corrected using independently calculated answers; [batch evidence and design decisions](2020-days07-19.md) include the reference comparisons.
- Ported all 49 parts for 2015, with original inputs, independent oracles, source-defect regressions and reference comparisons. See [2015 evidence and design decisions](2015.md).
- Ported all 49 parts for 2016, including executable implementations of day 23 part 2 and day 25. A shared interpreter handles all three Assembunny days; [2016 evidence](2016.md) records source corrections, proofs and search limits.
- Ported all 49 parts for 2017, including exact particle collisions, shared knot hashing and assembly arithmetic, and a verified composite-count shortcut. See [2017 evidence](2017.md).
- Ported all 49 parts for 2018, with complete-program VM guards, exact spatial searches and independent global-alignment evidence. See [2018 evidence and design decisions](2018.md).
- Ported all 49 parts for 2021, with finite cave-state counting, exact scanner and reactor geometry, infinite-image backgrounds and validated ALU constraints. See [2021 evidence and verification status](2021.md).
- Ported all 49 parts for 2022, including the original crate drawing, computed CRT text, exact rock-cycle state, validated robot/equation searches, generic cube folding and periodic valley routing. See [2022 evidence and verification status](2022.md).
- Added a shared six-row display decoder. Both 2016 day 8 and 2019 day 8 now derive their text from pixels; the latter no longer prints an image and returns a hardcoded answer.
- Improved the shared signed-number parser to reject lone signs and overflow, preserve error causes, and support explicit plus signs. New ports reuse the existing parsing, grid, set and heap helpers.
- Replaced full-pattern allocation in 2019 day 16 with prefix sums over alternating nonzero blocks. See [measurements and reference comparison](performance.md).
- Replaced the 2023 day 25 edge-ranking heuristic with deterministic exact capped flow, checking the global minimum and validating the resulting partition. Independent exhaustive graph tests and controlled benchmarks cover the change; see [contract, proof and measurements](performance.md#2023-day-25-exact-three-wire-partition).
- Replaced both 2023 day 17 searches with one Dijkstra search over legal straight runs. It reuses the scanner and heap helpers, validates the dense heat grid and handles zero costs and minimum stopping lengths. Independent step-state tests and controlled benchmarks cover both parts; see [proof and measurements](performance.md#2023-day-17-one-search-over-legal-strides).
- Corrected 2019 day 18's lost alternate routes using explicit entrance/key/door vertices and a memoized search over key pickups. Independent cell-by-cell checks cover one and four robots; controlled benchmarks improve both actual-input cases. See [contracts, proof and measurements](performance.md#2019-day-18-exact-routes-between-pickup-events).
- Corrected shared Intcode parsing, zero-extended parameter reads and continuation ownership. Painting text is computed through shared OCR, arcade counts final tiles, and droid/beam failures propagate. The existing binary-search helper now preserves predicate errors and handles integer extrema. [Contracts and regression evidence](reviews/2019-intcode.md) record the bounds and review.
- Resolved all 12 concrete findings in the existing 2019/2020 review, including token-safe movement compression, checked reaction search, affine composition, complete adventure continuations and empty-inventory trials. [Correction contracts and checks](reviews/2019-2020-fixes.md) record the remaining supported-domain limits.
- Corrected 2023 race counting, workflow ranges, hiking compression and reader failures. Hailstones now use exact determinants and shared rational elimination with complete candidate validation. The matrix helper owns arbitrary-precision immutable rationals and checked integer conversion; its 2025 joltage caller propagates failures. [Contracts, supported domains and checks](reviews/2023-exact.md) record the five resolved review findings.
- Corrected 2023 navigation using complete finite-state schedules and exact congruences. Repeated gardens now validate a finite shortest-distance certificate before counting, replacing the unchecked geometric shortcut. [Proofs and regression evidence](reviews/2023-navigation-garden.md) cover transient alignment, parity, isolation and supported limits.
- Replaced 2023 pulse first-HIGH extrapolation with complete branch cycles and verified FIFO event alignment, retaining exact whole-network simulation for other topologies. [Pulse contracts and proof](reviews/2023-pulses.md) complete all seven 2023 review findings.
- Corrected 2024 multiplication grammar, induced page ordering, disk gaps and primitive antinode lines, plus 2025 dial counting and exiting timeline branches. Verified the two formerly open part-two contracts against complete statement text. [Correction evidence](reviews/2024-2025-corrections.md) records independent oracles, checked boundaries and review.
- Corrected 2024 calibration row ownership, zero/concatenation semantics and claw press constraints, including dependent vectors. [Arithmetic evidence](reviews/2024-arithmetic.md) records independent oracles and the corrected full-input calibration benchmark.
- Corrected 2025 tile rectangles with discrete containment, and gift packing with sound certificates and bounded exact search. [Geometry evidence](reviews/2025-geometry.md) records the narrow-notch distinction, official packing example, independent checks and supported limits.
- Corrected the final secret-price window, worksheet trailing-space handling and pathfinding failure boundaries. [Boundary evidence](reviews/2024-2025-boundaries.md) records frozen independent expectations, shared helper reuse and error contracts.
- Replaced circuit warning heuristics with computed four-pair repairs and an all-input addition certificate. [Circuit evidence](reviews/2024-circuits.md) covers complete bounded search, ambiguity and independently checked repairs.
- Bounded robot-frame and chronospatial searches, preserving the verified personal tree specialization and admitting only proven octal quine loops. [Search evidence](reviews/2024-bounded-searches.md) records the domains, interpreter semantics and tests.
- Completed checked machine arithmetic and bounded exact optimization for 2025 day 10, preserving the shared rational solver. [Machine evidence](reviews/2025-machines.md) records the arithmetic proof and 832 independent small cases.
- Corrected the timing command: uncached runs, exact puzzle-test names, complete package results, and nonzero exit on execution, parsing, skipped-test or cancellation failures. Bounded stdout/stderr capture preserves diagnostics. Focused race tests, Staticcheck and a selected-2025 command run passed; all 23 part timings were reported.

Puzzle functions use `dayNNpPP` names in `internal/aocYEAR`, with input fixtures in `inputs/YEAR` and reader-based tests at that boundary. The 16 older unpadded 2019 names and their coverage records now follow that convention. Six existing factories configure reader-based solvers: 2024 day 14 part 1, both parts of days 18 and 20, and 2025 day 8 part 1. Their dimensions, input counts and saving thresholds are legitimate sample/personal parameters. All puzzle entrypoint bodies use the multiline layout of 2023/2024; 66 thin wrappers were expanded without changing their non-whitespace content. Ports use ordinary values and loops, standard-library operations, and existing parsing helpers where they make the code clearer.

## Completion

All 513 parts are implemented with canonical solver/test names and 262 nonempty input fixtures. Full-input assertions cover every part: 512 canonical test bodies load their fixture directly, while 2018 day 14 part 2 uses the deliberate `TestRecipeFullInput` test. Independent QA supplied the final missing 2023 day 2 goldens, 2476 and 54911.

All year-by-year reference comparisons and concrete findings are complete. The [2019/2020 review](reviews/2019-2020.md) identified 12 findings, the [2023 review](reviews/2023.md) identified seven, and the [2024–2025 review](reviews/2024-2025.md) identified 13 plus two contract questions; all are resolved with linked correction evidence and accepted specialist reviews. Supported search, arithmetic and input domains remain explicit in those records. The four ranked original slow cases are complete; [controlled measurements](performance.md) record the improvements.

The [ported-package style pass](reviews/ported-style.md) covered 175 production files. The [final style and coverage reconciliation](reviews/final-style.md) covered the remaining 95 production files and touched shared helpers and commands against the 2023/2024 baseline. It removed 11 redundant comments and added phase spacing in three functions; a Go lexer check verified that these final style edits preserve every executable token. No additional abstraction or broad naming rewrite was needed.

The timing command, canonical names, six configured factories and reference indexes were reconciled. The 2019 day 25 timing report includes all 49 parts, and the test-only scaffold wrapper was removed in favor of exercising its real parser/alignment behavior directly.

`go run ./cmd/test-timings -year YEAR` reports uncached top-level puzzle-test durations. It fails before printing a report if any requested year fails or produces incomplete output. These test durations rank investigation candidates; use isolated, repeated benchmarks for performance comparisons.

## Verification

Final checks passed on 2026-11-03 with Go 1.27.3 and Staticcheck 2026.2.2 (v0.8.2):

```sh
go mod tidy
staticcheck ./...
go test -race -shuffle=on -count=1 -timeout=3m ./...
```

The repository-wide race/shuffle run passed for every package and all 11 editions, without test caching. Formatting and `git diff --check` also passed. The final integration critic accepted the completion evidence and documented limits, with no further required work identified. The installed Staticcheck binary was upgraded. Sandbox verification used temporary writable Go build/module/Staticcheck caches; project dependencies remain unchanged apart from the Go directive. Performance evidence comes from the separate controlled benchmark runs, not these test durations.
