# 2023 Day 25 exact-cut comparison

Accepted: the exact fixed-source/all-target flow solver preserves both existing
answers and meets the predeclared performance thresholds. This directory records
the matched comparison from 2026-11-02.

| Workload | Baseline median | Exact solver median | Time reduction |
| --- | ---: | ---: | ---: |
| Official sample, answer 54 | 212.94 µs | 14.32 µs | 93.28% |
| Actual input, answer 614655 | 3.21410 s | 54.88 ms | 98.29% |

Actual allocation traffic fell from a median 1,083,008,496 to 1,033,640 B/op;
allocation count fell from 14,898,715.5 to 7,867 per operation. These measure
allocation traffic, not peak RSS or retained heap.

There were ten fresh process samples per implementation, with alternating A/B and
B/A order. Each process ran both fixtures with `b.Loop`, a fresh reader per
operation, `-test.benchtime=100ms`, `-test.cpu=1`, and PGO disabled. Parsing, solving,
answer formatting, and the fixed golden check are timed; fixture file loading is
excluded. The actual baseline takes longer than 100 ms, so each of its samples
contains one complete operation. Every operation passed. No samples were removed
or retried. `baseline-screen.txt` and `candidate-screen.txt` contain earlier
exploratory runs and are excluded from the comparison.

The accepted thresholds were an actual median improvement greater than 20%, a
sample median regression no greater than 10%, and correct output throughout.
Benchstat reports `p=0.000` (below 0.001) for both time comparisons. Its median
interval summaries are ±9%/±5% for actual baseline/candidate and ±8%/±12% for the
sample. Raw observations and full output are retained; no whole-suite or arbitrary
graph speedup is claimed.

Measurements used Go 1.27.3 on darwin/arm64, Apple M2 Max, GOMAXPROCS 1. Parent and
worker Go execution was stopped during measurement. CPU affinity and thermal/power
state were not controlled. The baseline uses randomized map traversal; all
observed goldens passed, but timing success does not establish its correctness on
other graphs. The exact solver has a stronger validated contract.

## Evidence

- `benchmark_test.go.txt`: identical harness used to build both confirmation binaries.
- `baseline.go.txt`, `candidate.go.txt`, and corresponding test snapshots: pinned overlays.
- `candidate.patch`: complete production/test change against the pinned baseline.
- `plan.json`: thresholds, sample plan, fixture and source identities.
- `baseline.txt`, `candidate.txt`, `order.json`: every confirmation observation and command.
- `benchstat.txt`, `summary.json`: statistical output and numeric summaries.
- `environment.json`: toolchain, machine, dirty revision, settings, tool version, and limitations.
- `SHA256.json`: artifact checksums. Input `inputs/2023/25.txt` is identified by its hash in `plan.json`.

The initial baseline screen used `screen-original-benchmark.go.txt`, whose file
loading path was absolute. Before confirmation, that untimed path was made
package-relative and both binaries were rebuilt with the same updated harness.

## Replay

From this repository, with no competing builds, tests, profilers, or load generators:

```sh
python3 migration/measurements/2023-day25/replay.py /tmp/day25-comparison-new
```

The output directory must not exist. The script verifies artifacts, module files,
and the actual fixture; builds two binaries from frozen Day 25 source/test overlays;
then performs the declared ten alternating samples. It stops on any failed golden
without retrying. It uses the current repository for all other package and helper
files, which must match the recorded migration state for an equivalent replay.
Future Go/toolchain/helper changes define a new comparison, not a reproduction of
the original result. Benchstat runs if installed on PATH; raw samples remain usable
otherwise.

`--prepare-only` verifies inputs and writes overlays without compiling or measuring.
Normal execution has 60-second build/process limits and a five-minute measurement
budget. Original confirmation completed in 37.46 seconds, excluding builds.
