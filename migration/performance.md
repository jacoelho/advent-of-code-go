# AoC performance investigations

## Decision

Accept the 2019 day 16 part 1 prefix-sum rewrite. It removes repeated full-pattern allocation and iterator dot products; each nonzero pattern run becomes one prefix-sum subtraction. Existing part 2 suffix summation remains unchanged. Complexity of the first half changes from quadratic to O(n log n) per phase, using two linear-sized arrays rather than one pattern allocation per output digit. The implementation is 14 lines shorter overall.

Source idea: `reference/.repos/CameronAavik--AdventOfCode/fsharp/2019/Solutions/Day16.fs`, `applyFFTPhase`, which sums alternating nonzero blocks using a prefix array. Upstream: https://github.com/CameronAavik/AdventOfCode/blob/1f1313a67b959e21509300a54925ad0ec7a2a23b/fsharp/2019/Solutions/Day16.fs (local manifest pins the collected revision).

## Measurement

Apple M2 Max, darwin/arm64, Go 1.27.3, module go1.27. Baseline revision e4f5a9acdef63cc203de5578a8eb1b546e6c6b54 production day16 with the same new benchmark harness as candidate; only day16 implementation differs in these binaries. PGO disabled for both. Explicit benchmark GOMAXPROCS=1. Parent granted an exclusive timing window; external desktop load and power/thermal state were not instrumented.

Predeclared acceptance: >20% lower real-input time, no >10% sample regression. Ten samples per version, alternating A/B then B/A; 100ms requested duration, fixed golden outputs checked on every operation. The slower baseline real-input fixture ran one operation per sample; the framework gives candidate enough iterations to fill its duration. One operation includes parsing a fresh reader, 100 FFT phases, formatting the answer; fixture file reading is setup only.

| Workload | Baseline median | Candidate median | Change | Baseline allocation | Candidate allocation |
| --- | --- | --- | --- | --- | --- |
| Official 32-digit sample | 3.51223 ms (±2%) | 26.34 µs (±6%) | -99.25%, p=0.000, n=10 | 1100.82 KiB; 20,904 allocations | 53.93 KiB; 204 allocations |
| Actual input, 650 digits | 1.2570423 s (±2%) | 670.9 µs (±6%) | -99.95%, p=0.000, n=10 | 179.544 MiB; 422,607 allocations | 1.033 MiB; 207 allocations |

The percentages/ranges are benchstat's reported comparison and confidence ranges. These are end-to-end solver benchmark results on these fixtures, not whole-suite speedup claims. No profiler was mixed into acceptance timing.

Commands:

```sh
GOCACHE=/tmp/aoc-performance-cache GOMODCACHE=/tmp/aoc-performance-modcache go test -c -pgo=off -o /tmp/aoc-performance-before.test ./internal/aoc2019
# Built candidate after production edit, with identical test and benchmark source:
GOCACHE=/tmp/aoc-performance-cache GOMODCACHE=/tmp/aoc-performance-modcache go test -c -pgo=off -o /tmp/aoc-performance-after.test ./internal/aoc2019
# Each binary runs from internal/aoc2019, ten times in alternating order:
/tmp/aoc-performance-before.test -test.run='^$' -test.bench='^BenchmarkDay16Part1$' -test.benchmem -test.benchtime=100ms -test.count=1 -test.cpu=1
benchstat /tmp/aoc-performance-before.txt /tmp/aoc-performance-after.txt
```

Preserved evidence: [baseline samples](measurements/day16-before.txt), [candidate samples](measurements/day16-after.txt), [benchstat](measurements/day16-benchstat.txt), [execution order](measurements/day16-order.json), [environment and hashes](measurements/day16-environment.json), [baseline source](measurements/day16-before.go.txt), and [implementation diff](measurements/day16.patch). The temporary binaries can be rebuilt from these sources and the committed benchmark harness.

## Tests and scope

QA consulted before tests. Retained three official 100-phase examples and actual input for each part. Added the official four-phase sequence, odd-length clipped-run case `123456789 -> 583250479`, and singleton case. Added one end-to-end benchmark with sample and actual fixtures and golden answer checks. No fixed-time/allocation-count assertions, production instrumentation, parser-contract change, or unrelated helper edits. Empty/malformed input semantics are outside this behavior-preserving optimization.

Changed:
- `internal/aoc2019/day16.go`
- `internal/aoc2019/day16_test.go`
- `internal/aoc2019/day16_bench_test.go` (new)

Focused tests, gofmt, and targeted race/shuffle tests completed successfully (`go test -race -count=1 -shuffle=on -timeout=2m -run '^(Test_day16p0[12]|TestApplyFFTPhase)$' ./internal/aoc2019`, 2.082s). The full repository race/shuffle suite also passed after integration; completion review follows the batch. gopls MCP operations were unavailable, so reference discovery used scoped text inspection; generatePattern had only one caller in this file and no direct tests.

## Initial suite screening

Before this optimization all five existing AoC packages passed uncached serial tests using `go test -count=1 -p=1 -parallel=1 -timeout=110s -json ./internal/aocYEAR`; 120-second outer cap per package, initial measurement budget under 8 minutes. Successful baseline used the same installed Go1.27.3 after the module directive upgrade. Initial attempt could not download x/exp/write the default module cache; rerun used temporary caches and authorized download. No suite timed out.

| Package | Test package duration |
| --- | --- |
| 2019 | 4.736s |
| 2020 | 2.203s |
| 2023 | 5.719s |
| 2024 | 1.909s |
| 2025 | 0.715s |

Slowest top-level tests in initial screening: 2023/25p1 2.69s; 2019/18p2 1.53s; 2023/17p2 1.14s; 2019/16p1 1.13s. Three uncached repetitions gave ranges 3.16–3.75s, 1.56–1.69s, 1.50–1.73s, and 1.13–1.15s respectively. Times include each test's samples plus actual input. Initial suite screening is ranking evidence, not a controlled optimization comparison.

The [screening summary](measurements/screening-summary.json) records the package/test timings. Raw screening JSON streams remain in temporary files; the controlled before/after benchmark evidence is preserved above.

## 2023 day 25: exact three-wire partition

The former solver ranked edges by all-source BFS path frequency, then tried triples among the twenty most popular edges. That restriction was incomplete and depended on map traversal order. Checking only for two components after deleting three edges could also accept a bridge plus two irrelevant removals.

The replacement assigns sorted integer vertex IDs and runs unit-capacity residual flow from one fixed source to every other vertex. Four successful augmentations prove that a target's minimum cut is at least four. An exhausted search after three yields a three-cut certificate; exhaustion earlier rejects the graph. Every global cut separates the fixed source from some target, so checking all targets establishes the global minimum. A discovered three-cut does not end that verification. Reverse residual capacities can reach two and preserve cancellation. The retained partition is checked against the original graph for exactly three crossing wires and two nonempty connected components.

The supported input is a simple undirected graph whose global minimum edge cut is exactly three. Repeated and reciprocal declarations describe the same wire; isolated declarations are preserved and self-loops rejected. When several minimum partitions exist, sorted graph traversal selects one deterministically; their products need not be equal. This deliberately rejects accidental successes of the old three-removal heuristic. Input is bounded to 4 MiB, 4,096 vertices, 20,000 wires and 128-byte names, with the shared scanner's line limit. Search is bounded to 200 million adjacency visits. A limit produces an error, never a heuristic answer. For fixed cut size three, time is O(V(V+E)) and working memory O(V+E).

The [AxlLind reference](https://github.com/AxlLind/AdventOfCode2023/blob/7502ebced52fa59cf17016e571f2116f7c7430e9/src/bin/25.rs) supplies the capped-flow idea. Its early acceptance of any cut at most three was not retained. The [CameronAavik reference](https://github.com/CameronAavik/AdventOfCode/blob/1f1313a67b959e21509300a54925ad0ec7a2a23b/csharp/2023/Solvers/Day25.cs) illustrates indexed residual state, but its farthest-pair assumption has acknowledged counterexamples. Dense Stoer–Wagner would add quadratic storage and cubic work; randomized contraction would introduce bounded-search failure or unbounded retries. Neither is needed for this fixed-cut contract.

Named advisor and QA reviewed the proof and tests before implementation. Official and actual answers remain 54 and 614655. Independent partition enumeration checks all 1,024 simple five-vertex graphs; further cases cover a proven unique K5/K6 three-wire join, smaller/larger cuts, a later smaller cut after an initial three-cut, residual cancellation, isolated vertices, and deterministic reordered/duplicate declarations. Focused race/shuffle tests passed in 2.021s, followed by the full 2023 package in 17.269s. Repository-wide Staticcheck passed, and completion code review found no material issue.

Accept the replacement. On Apple M2 Max, darwin/arm64, Go 1.27.3, ten samples per implementation used identical end-to-end benchmark code, PGO disabled and GOMAXPROCS=1. Each operation parses a fresh reader, solves and formats the answer, checking the fixed golden; fixture loading is untimed. Predeclared acceptance required more than 20% actual-input improvement, no sample regression above 10%, and correct answers throughout. Alternating A/B then B/A runs used a requested 100ms benchmark duration; the slower baseline actual case naturally ran one operation per sample. Exploratory one-operation screens are preserved separately and excluded from these statistics.

| Workload | Baseline median | Candidate median | Change | Baseline allocation traffic | Candidate allocation traffic |
| --- | --- | --- | --- | --- | --- |
| Official sample | 212.94 µs (±8%) | 14.32 µs (±12%) | -93.28% | 122.69 KiB; 1,816 allocations | 14.45 KiB; 102 allocations |
| Actual input, 1,568 vertices / 3,503 wires | 3.21410 s (±9%) | 54.88 ms (±5%) | -98.29% | 1,057,625.5 KiB; 14,898,716 allocations | 1,009.4 KiB; 7,867 allocations |

All 20 confirmation processes passed every golden. Benchstat reports p=0.000 with ten samples per side for both timing comparisons; the percentages and ranges above are its output. The actual fixture improves about 58.6-fold. Allocation figures measure traffic per operation, not peak memory. Measurements ran without competing agent workloads; desktop load, affinity and thermal/power state were not controlled. These two fixtures establish solver improvements, not a whole-suite speedup. The baseline remains a map-order-dependent heuristic, even though all observed answers passed.

Preserved evidence: [benchmark comparison](measurements/2023-day25/benchstat.txt), [baseline samples](measurements/2023-day25/baseline.txt), [candidate samples](measurements/2023-day25/candidate.txt), [execution order](measurements/2023-day25/order.json), [predeclared plan](measurements/2023-day25/plan.json), [environment and source provenance](measurements/2023-day25/environment.json), [matched benchmark](measurements/2023-day25/benchmark_test.go.txt), [implementation diff](measurements/2023-day25/candidate.patch), and [artifact hashes](measurements/2023-day25/SHA256.json). Both implementation and test snapshots are retained beside these files; [replay instructions](measurements/2023-day25/README.md#replay) rebuild them through temporary overlays.

## 2023 day 17: one search over legal strides

Accept one dense Dijkstra search using the existing scanner and heap helpers. The previous solver stored position/direction/run-length maps, allocated neighbours and predecessor paths, and searched twice for the two initial directions. The new state is `(cell, next axis)`; each edge spans a complete straight run of 1–3 cells in part one or 4–10 in part two, charging every entered cell. Both origin axes enter one queue at zero cost. Strict improvements and stale-entry rejection handle zero-cost cycles without a heuristic.

Every legal nonempty route decomposes into straight runs within those limits; every path through the compressed graph expands into such a route. At a legal endpoint, the next run must use the perpendicular axis, so the previous sign and length no longer affect its choices. Goals are accepted only at legal endpoints, preserving the ultra crucible's minimum final run. The old constant-one non-goal heuristic could overestimate a remaining zero-cost route and is not used as an oracle for zero cells.

The parser accepts nonempty rectangular ASCII digits 0–9, up to 250,000 cells, subject to the existing scanner's line limit. Malformed input and reader errors propagate. An unreachable destination now returns an error instead of the private `-1` sentinel. A singleton returns zero for part one and errors for part two, which still requires a legal final run. At N cells there are 2N states, at most 40N stride-cell visits and 28N directed transitions. These bounds also bound heap insertions and memory; no second search budget or queue implementation is needed. With fixed stride limits, time is O(N log N) and memory O(N).

The [CameronAavik reference](https://github.com/CameronAavik/AdventOfCode/blob/1f1313a67b959e21509300a54925ad0ec7a2a23b/csharp/2023/Solvers/Day17.cs) supplied the whole-stride/axis simplification. Its fixed buckets, bitsets, packed ushort states and transformed Manhattan costs were unnecessary here; ordinary Dijkstra supports zeros and unreachable grids directly. The map-backed shared `Grid2D` adds no value for this dense rectangle, so the heat grid remains a small private type.

QA approved tests before edits. Existing answers 102/94/71 and 755/881 pass. Literal cases cover both initial axes, minimum/maximum run lengths, origin heat, singleton and unreachable semantics, and a zero-cost grid whose ultra answer must be 36 rather than an illegal early stop at zero. A synchronous one-cell Bellman–Ford oracle retains direction and run length, independently checking both parts on 49 seeded rectangular grids from 1×1 through 7×7. Parser tests cover shape, digits, the cell cap and reader-error identity. Focused race/shuffle tests passed in 1.681s; the full 2023 package then passed in 9.886s. Repository-wide Staticcheck passed.

Ten fresh samples per implementation alternated A/B and B/A, using the same four-case benchmark, Go 1.27.3, PGO off and GOMAXPROCS=1 on Apple M2 Max/darwin arm64. Each operation includes a fresh reader, parsing, search, formatting and a fixed golden check; file loading is untimed. Predeclared acceptance required both actual medians to improve by more than 20%, neither sample to regress more than 10%, and correct results throughout. Requested duration was 100ms; each baseline actual case naturally ran one operation. Exploratory screens were excluded from acceptance statistics.

| Workload | Baseline median | Candidate median | Change | Baseline allocation traffic | Candidate allocation traffic |
| --- | --- | --- | --- | --- | --- |
| Part 1 sample | 1.58372 ms (±3%) | 55.87 µs (±1%) | -96.47% | 2016.66 KiB | 29.80 KiB |
| Part 1 actual, 141×141 | 393.435 ms (±3%) | 9.156 ms (±2%) | -97.67% | 372.944 MiB | 2.496 MiB |
| Part 2 sample | 1.14121 ms (±2%) | 41.73 µs (±1%) | -96.34% | 1594.47 KiB | 29.41 KiB |
| Part 2 actual, 141×141 | 1.28062 s (±1%) | 14.10 ms (±1%) | -98.90% | 943.422 MiB | 3.334 MiB |

All 20 confirmation processes passed every golden, without discarded or retried samples. Benchstat reports p=0.000 with ten samples per side for all four timing comparisons; ranges and percentages above are its output. Allocation figures are traffic, not peak memory. Other agent workloads stopped during measurement; affinity and thermal/power state were not controlled. These are complete solver measurements for the supplied fixtures, not whole-suite speedups. Changed invalid/unreachable behaviour is covered separately and does not contribute to the performance comparison.

Preserved [evidence and replay instructions](measurements/2023-day17/README.md) include raw samples, the matched harness, frozen production/test snapshots, the predeclared plan, execution order, environment, hashes and [benchstat output](measurements/2023-day17/benchstat.txt).

## 2019 day 18: exact routes between pickup events

Accept the corrected semantic graph and memoized collection search. The old graph retained one shortest route and door mask per entrance/key pair; a locked shortest route could hide a longer unlocked alternative. It also ignored reader failures and returned zero when search failed. The replacement keeps every entrance, key and individual door tile as a vertex. Preprocessing BFS stops at other such vertices, so edge interiors contain only ordinary floor. Among corridors with the same endpoints, the shortest dominates because their permissions and pickup effects are identical.

For the current shared key mask, Dijkstra blocks locked doors, traverses entrances and owned keys, and stops at each first unowned key. The outer memoized search moves one robot to one newly collected key and minimizes the remaining total distance. Each transition adds a key, bounding recursion depth to 26; no outer heap or predecessor paths are needed. Robots share keys and do not block each other. A robot's partial moves can therefore be deferred until its next pickup without increasing cost or losing permissions: keys only accumulate and other robots impose no obstacles. This establishes the reduction from individual movement to pickup events, including connected four-robot layouts.

The parser requires rectangular ASCII `# . @ a-z A-Z`, unique lowercase keys and valid entrance counts. Repeated doors are separate vertices; unmatched doors stay locked. Part one requires one entrance. Part two accepts four entrances or transforms one only when its complete 3×3 neighborhood is `.../.@./...`; it cannot erase keys/doors or create passages through walls. Valid zero-key layouts return zero; malformed, unreadable and unreachable layouts return errors.

Limits are 100,000 cells, 256 semantic vertices, one million memo entries and 200 million shared work units across preprocessing, inner searches and key transitions. Exhaustion propagates as an error rather than an unreachable or partially explored answer. Preprocessing costs O(PN) for N cells and P semantic vertices, with at most P(P−1) graph edges. The exponential key-order search remains explicitly bounded; a finite 26-key universe alone is not a practical resource limit. Reachability results are not cached separately, and robot positions remain labelled. A second cache, symmetry reduction and Pareto route-label frontiers are unnecessary for the accepted performance.

The selected [Scala reference](https://github.com/maneatingape/advent-of-code-scala/blob/063a22ee01c1d489f76851e4e4971b0b9b767903/src/main/scala/AdventOfCode2019/Day18.scala) and [F# reference](https://github.com/CameronAavik/AdventOfCode/blob/1f1313a67b959e21509300a54925ad0ec7a2a23b/fsharp/2019/Solutions/Day18.fs) also summarize one route/door set, so neither validates the corrected detour semantics. The outer memoized remaining-cost model is useful; hardcoded quadrant transforms and selecting only one robot for a reachable key are not retained. Dynamic grid BFS for every key state would be exact, but the actual vault has 3,201 walkable cells and only 53 semantic vertices, or 56 after splitting the entrance. Corridor compression removes repeated floor traversal while keeping the state-dependent tiles explicit.

Named advisor and QA assessed the model before dependent tests. Existing official answers and actual 3862/1626 pass. Literal cases cover longer and equal-distance unlocked alternatives, keys before doors, traversal through owned keys, repeated doors, connected robots sharing keys, safe splitting, parser failures and work exhaustion. An independent one-cell BFS moves one labelled robot per step and checks 80 seeded tiny layouts, including connected four-robot maps. The million-entry memo ceiling was inspected statically rather than exercised with a costly synthetic allocation. Focused race/shuffle tests passed in 2.091s; the complete 2019 package passed in 11.288s. Repository-wide Staticcheck passed.

Ten fresh samples per implementation alternated A/B and B/A using the same four-case benchmark, Go 1.27.3, PGO off and GOMAXPROCS=1 on Apple M2 Max/darwin arm64. Each operation includes a fresh reader, parsing, transformation, graph construction, collection search, formatting and a fixed golden check; fixture loading is untimed. Predeclared acceptance required both actual medians to improve by more than 20%, neither sample to regress more than 10%, and correct answers throughout. Requested duration was 100ms; slower actual baselines naturally ran one operation. Exploratory screens are preserved separately and excluded from acceptance.

| Workload | Baseline median | Candidate median | Change | Baseline allocation traffic | Candidate allocation traffic |
| --- | --- | --- | --- | --- | --- |
| Part 1 branching sample, answer 136 | 66.93 ms (±5%) | 14.39 ms (±6%) | -78.49% | 30.03 MiB | 11.33 MiB |
| Part 1 actual | 274.20 ms (±8%) | 14.85 ms (±4%) | -94.58% | 111.485 MiB | 9.908 MiB |
| Part 2 four-robot sample, answer 72 | 478.5 µs (±1%) | 106.0 µs (±10%) | -77.84% | 447.1 KiB | 173.9 KiB |
| Part 2 actual | 1.75756 s (±5%) | 34.29 ms (±10%) | -98.05% | 1011.05 MiB | 42.13 MiB |

All 20 confirmation processes passed every golden, without discarded or retried samples. Benchstat reports p=0.000 with ten samples per side for all four comparisons; ranges and percentages above are its output. Part one's branching sample performs 44.68% more allocations, while its allocated bytes fall 62.27% and time falls 78.49%; this tradeoff does not justify another queue/cache implementation. Allocation traffic is not peak memory. Other agent workloads stopped during measurements; affinity and thermal/power state were not controlled. These fixture results do not imply whole-suite speedups or certify the old solver's behavior on arbitrary vaults.

Preserved [evidence and replay instructions](measurements/2019-day18/README.md) include the identical benchmark, production/test snapshots, raw samples, execution order, predeclared plan, environment, hashes and [benchstat output](measurements/2019-day18/benchstat.txt).

## Screening follow-up

The four ranked original slow cases—2019 days 16/18 and 2023 days 17/25—have now been investigated, corrected or simplified, and measured. Ported-year performance checks are recorded in their batch documents. Further measurement is warranted only when the remaining reference/idiom review exposes a concrete new cost mechanism or an implementation change needs comparison.

## Timing command

`cmd/test-timings` now invokes `go test -json -count=1 -timeout=3m` for each requested year, with a five-minute process deadline. Only exact `Test_dayNNpPP` names contribute timings. Malformed/truncated JSON, missing package completion, failed or skipped primary tests, command failure and cancellation return errors; a report is printed only after every requested year succeeds. Stdout and stderr are captured separately, capped at 16 MiB and 1 MiB, and an exceeded limit is an error. Failure diagnostics include bounded excerpts from both streams.

Focused race/shuffle tests cover the former overwritten-timing and false-success defects, incomplete output, cancellation, subprocess diagnostics and output limits. Staticcheck passed, and a selected-2025 command run returned all 23 puzzle timings successfully. The 2019 day 25 test was renamed to the established `Test_day25p01` convention; a subsequent 2019 command run returned all 49 expected part timings. Reported durations include each primary test's examples and actual fixture, use Go's event timing precision, and remain screening evidence rather than controlled benchmark measurements.
