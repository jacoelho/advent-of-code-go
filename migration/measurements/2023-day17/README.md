# 2023 Day 17 stride-search comparison

Accepted: one dense Dijkstra search over cell/axis states preserves all four
benchmark answers and meets the predeclared performance thresholds. Measurements
were collected on 2026-11-02.

| Complete solver operation | Baseline median | Candidate median | Time reduction |
| --- | ---: | ---: | ---: |
| Part 1 sample, answer 102 | 1.583721 ms | 55.8675 µs | 96.47% |
| Part 1 actual, answer 755 | 393.435229 ms | 9.156420 ms | 97.67% |
| Part 2 sample, answer 94 | 1.141211 ms | 41.7265 µs | 96.34% |
| Part 2 actual, answer 881 | 1.280625 s | 14.101063 ms | 98.90% |

Actual median allocation traffic fell from 391,060,584 to 2,617,488 B/op for part
one and from 989,249,544 to 3,495,872 B/op for part two. Allocation counts fell from
2,646,304 to 133,073 and from 7,988,203 to 170,056 per operation. These are allocation
traffic and counts, not retained heap or peak RSS.

Ten fresh process samples per implementation ran serially, alternating A/B and B/A.
Each process ran all four fixtures using `b.Loop`, fresh readers, fixed golden
checks on every operation, `-test.benchtime=100ms`, `-test.cpu=1`, and PGO disabled.
Parsing, search, answer formatting, and result validation are included; fixture
loading is excluded. Each baseline actual sample contains one complete operation
because that operation exceeds 100 ms. Every operation passed. No samples were
dropped or retried. The earlier three `1x` screens per implementation are retained
separately and excluded from acceptance statistics.

The declared thresholds required each actual median to improve by more than 20%,
each sample median to regress no more than 10%, and correct output throughout.
Benchstat reports `p=0.000` (below 0.001) for all four time comparisons. Its median
interval summaries are ±3%/±1% for part one sample, ±3%/±2% for part one actual,
±2%/±1% for part two sample, and ±1%/±1% for part two actual. Raw observations and
full statistical output are retained.

The host was an Apple M2 Max running darwin/arm64 and Go 1.27.3. Other agents
stopped Go execution during measurement. CPU affinity and thermal/power state were
not controlled. These results concern complete solver calls for these fixtures;
they do not establish whole-suite speedups or performance on every grid shape.
The positive-digit sample and actual inputs are shared valid-domain anchors.
New parser validation and unreachable errors are covered by separate correctness
tests; their changed behavior is not counted as a performance gain.

## Evidence

- `benchmark_test.go.txt`: the identical harness used to build both binaries.
- `baseline.go.txt`, `candidate.go.txt`, and matching test snapshots: frozen overlays.
- `candidate.patch`: production and test changes against the baseline snapshots.
- `plan.json`: predeclared thresholds, sample plan, source and fixture identities.
- `baseline.txt`, `candidate.txt`, `order.json`: all confirmation observations and commands.
- `benchstat.txt`, `summary.json`: statistical output and numeric summaries.
- `environment.json`: toolchain, machine, dirty revision, settings, tool version, and limitations.
- `SHA256.json`: artifact checksums. The actual fixture is identified by its hash in `plan.json`.

## Replay

From this repository, with no competing builds, tests, profilers, or load generators:

```sh
python3 migration/measurements/2023-day17/replay.py /tmp/day17-comparison-new
```

The output directory must not exist. The script verifies artifacts, module files,
and the actual fixture; builds two binaries from the frozen Day 17 source/test
overlays; then executes the same ten alternating samples. Failed goldens stop the
run without retrying. Other package/helper files come from the current repository,
so equivalent replay requires the recorded migration state. Toolchain or helper
changes define a new comparison. Benchstat runs if installed on PATH; raw samples
remain available otherwise.

`--prepare-only` verifies inputs and creates overlays without Go execution. Normal
execution has 60-second build/process limits and a five-minute measurement budget.
Original confirmation completed in 24.24 seconds, excluding builds.
