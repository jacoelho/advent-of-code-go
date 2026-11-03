# 2019 Day 18 exact-route comparison

Accepted: the graph with explicit entrances, keys and doors preserves all four
benchmark answers and meets the predeclared thresholds. Measurements were collected
on 2026-11-03.

| Complete solver operation | Baseline median | Candidate median | Time reduction |
| --- | ---: | ---: | ---: |
| Part 1 official branching example, answer 136 | 66.925667 ms | 14.393169 ms | 78.49% |
| Part 1 actual, answer 3862 | 274.202104 ms | 14.850923 ms | 94.58% |
| Part 2 official four-robot example, answer 72 | 478.5255 µs | 106.035 µs | 77.84% |
| Part 2 actual, answer 1626 | 1.757560 s | 34.285236 ms | 98.05% |

Actual median allocation traffic fell from 116,900,320 to 10,388,968 B/op for part
one and from 1,060,161,832 to 44,178,013 B/op for part two. Counts fell from
1,439,041.5 to 402,390 and from 3,712,494 to 1,080,537 allocations per operation.
These measure allocation traffic, not retained heap or peak RSS.

The branching sample has a tradeoff: allocation count rises 44.68%, from 291,593.5
to 421,873, while allocated bytes fall 62.27% and elapsed time falls 78.49%.
Repeated small inner searches are consistent with this tradeoff, but no allocation
profile was collected to attribute individual allocations. No extra caching or
queue implementation was added after the first candidate met the objectives.

Ten fresh process samples per implementation ran serially, alternating A/B and B/A.
Each process ran all four fixtures using `b.Loop`, fresh readers, fixed golden
checks on every operation, `-test.benchtime=100ms`, `-test.cpu=1`, and PGO disabled.
Parsing, entrance transformation, graph construction, key search, formatting and
result validation are included; fixture loading is excluded. Each baseline actual
sample contains one full operation because it exceeds 100 ms. Every operation
passed. No samples were dropped or retried. Earlier three-run `1x` screens are
retained separately and excluded from acceptance statistics.

The declared thresholds required each actual median to improve by more than 20%,
each official guardrail to regress no more than 10%, and correct output throughout.
Benchstat reports `p=0.000` (below 0.001) for all four time comparisons. Its median
interval summaries are ±5%/±6% for part one sample, ±8%/±4% for part one actual,
±1%/±10% for part two sample, and ±5%/±10% for part two actual. Raw observations
and full statistical output are retained.

The host was an Apple M2 Max running darwin/arm64 and Go 1.27.3. Other agents
stopped Go execution during measurement. CPU affinity and thermal/power state were
not controlled. These complete-solver measurements do not establish whole-suite
speedups or performance for every maze shape. Existing official/personal goldens
are shared valid-domain anchors; independent route regressions validate the
corrected alternative-door-route semantics. New validation and unreachable errors
are separate correctness changes, not shortcuts counted as performance gains.

## Evidence

- `benchmark_test.go.txt`: identical benchmark source for both binaries.
- `baseline.go.txt`, `candidate.go.txt`, and corresponding test snapshots: frozen overlays.
- `candidate.patch`: production/test changes against the frozen baseline.
- `plan.json`: predeclared thresholds, sample plan, and source/fixture identities.
- `baseline.txt`, `candidate.txt`, `order.json`: every confirmation observation and command.
- `benchstat.txt`, `summary.json`: statistical output and numeric summaries.
- `environment.json`: toolchain, machine, dirty revision, settings and limitations.
- `SHA256.json`: artifact checksums. The actual input hash is recorded in `plan.json`.

The parent's mechanical private-solver renames and plain-scaffold wrapper removal
were present in both builds. Only the Day 18 source/test overlays differed.

## Replay

From this repository, with no competing builds, tests, profilers, or load generators:

```sh
python3 migration/measurements/2019-day18/replay.py /tmp/day18-comparison-new
```

The output directory must not exist. The script verifies artifacts, module files,
and the actual fixture; builds two binaries from frozen Day 18 source/test overlays;
then executes ten alternating samples. Failed goldens stop the run without retrying.
Other package/helper files come from the current repository, so equivalent replay
requires the recorded migration state. Toolchain/helper changes define a new
comparison. Benchstat runs if installed on PATH; raw samples remain available
otherwise.

`--prepare-only` verifies inputs and creates overlays without Go execution. Normal
execution has 60-second build/process limits and a five-minute measurement budget.
Original confirmation completed in 27.71 seconds, excluding builds.
