# 2023 navigation and repeated gardens

This batch resolves day 8's portion of F3 and F6 in the [2023 review](2023.md). Day 20 pulse synchronization remains open. Both changes retain private day-local models and existing scanner, grid and search helpers; neither imports a reference implementation's input-specific assumptions.

## Navigation schedules

Day 8's state is `(node, instruction index)`. Part one stops at the first target or a repeated state. Part two records every start's transient target hits, cycle-entry time, period and target residues. Let T be the latest cycle-entry time. Any shared hit below T occurs among that latest-entry walk's transient hits; those candidates are tested against every schedule. From T onward, generalized Chinese remaindering combines all compatible residues, including noncoprime periods. Each resulting residue is lifted to its first occurrence at least T, and the minimum is returned.

Periods and combined answers use `math/big`. Existing machine-integer LCM cannot satisfy that contract. Other exact congruence implementations are private to earlier puzzle packages; no cross-package access or general scheduling framework was added. The collected CameronAavik day 8 reference at `1f1313a67b959e21509300a54925ad0ec7a2a23b` explicitly assumes special cycle structure, which is not a valid replacement for the observed defect.

Parsing checks L/R instructions, unique alphanumeric names, two destinations and complete references. Bounds are 4 MiB, 100000 instructions, 100000 nodes, 128-character names and 1024 starts. State expansion is limited to one million states, shared across part-two walks. Alignment permits one million candidate checks and at most 10000 combined residues. Exhaustion returns a resource error; incompatible schedules return a no-solution error.

QA reviewed the tests before edits. Independent cases cover the original 11-versus-6 counterexample, transient-only alignment, multiple residues, noncoprime periods, cycle-entry lower bounds, instruction phase, unreachable targets and reader failures. Ten literal prime cycles have an independently calculated first common hit of 421336842070675358938, beyond int64. Existing personal results remain 13771 and 13129439557681.

## Certified garden distances

The [part-one statement](https://adventofcode.com/2023/day/21) asks for positions reachable in exactly the requested steps. The [reproduced part-two statement](https://github.com/AntonKosov/advent-of-code-2023/blob/61cdf160b4dd/day21/part2/README.txt), verified at blob `521ab5ee9747bae25a429bcbddde6a5a0a36ac3a`, repeats the garden infinitely with only the original start active.

Finite part one retains shared BFS and distance/parity counting. An isolated start has no positive-length walk, so it is excluded at every positive target. Parsing now requires a rectangular ASCII grid of `.#S` with exactly one start and at most one million cells. Reader errors retain their cause.

Part two supports odd square tiles of side L at most 201 whose distances pass a fixed certificate. It does not require a centered start or clear axes/borders. A shared BFS in a 5×5 repeated patch supplies candidate distances in the nine representative tiles with coordinates −1, 0 and 1. For a cell p in tile q, the proposed distance is:

```text
D(q,p) = d(clamp(qx,-1,1), clamp(qy,-1,1), p)
         + L * (max(|qx|-1,0) + max(|qy|-1,0))
```

Absent base distances remain infinite. All open cells in the 25 tile regimes `[-2,2]²` are checked against all four neighbors, including neighbors in tiles at coordinate ±3. The checks require the actual start to be the only zero, adjacent finite labels to differ by at most one, no edge between finite and infinite labels, and a predecessor of distance D−1 for every positive finite label. Periodic neighbors use signed offsets directly; the existing nonnegative-only `Neighbours4` remains appropriate for finite part one.

Every farther one-hop neighborhood translates one of these checked regimes and adds the same constant to every finite label. Therefore the checks hold globally. Nonnegative descent reaches the unique zero in D steps; the edge bound makes every path at least D steps long. No finite/infinite edge can omit a reachable cell. The bounded BFS proposes labels; this certificate proves them. A failed certificate returns an unsupported-distance-pattern error, not an unreachable result or a fitted answer. The side bound limits BFS to 1010025 cells.

Counting partitions all tiles into the center, four rays and four quadrants. For a representative distance d≤N, let `m=(N-d)/L` and `p=(N-d)%2`. Since L is odd, eligible extra tile distances have parity p. If m≥p, set `q=(m-p)/2+1`: a ray contributes q and a quadrant q(q+p). The center contributes only when its distance has the target parity. Products and totals use `big.Int`. An isolated periodic start is handled separately because distance/parity counting requires an edge for two-step backtracking.

The advisor checked the global certificate, counting proof and supported boundary. The personal input satisfies the certificate and retains 3598/601441063166538. The old formula's clear-axis/border/diamond assumptions, also present in the collected CameronAavik reference, are no longer prerequisites or unchecked claims. There is no polynomial fitting, increasing-radius search or special-case personal answer.

QA independently checked literal fixtures with exact step-by-step frontiers before tests were written. Required accepted cases include the single-cell plane, an off-center open tile, a line, a bounded pair, and isolated/unreachable components. The one-cell result is 702322399865956 at puzzle steps and 16000000008000000001 at four billion steps. A diagonal staircase must return an unsupported error: its distances do not fit this model. All 256 centered 3×3 wall patterns are also checked against independent raw frontiers for steps 0..12 when admitted; rejected cases must identify the unsupported distance pattern. Explicit acceptance and rejection fixtures prevent blanket rejection from passing that comparison. The published finite six-step example remains 16.

## Checks and review

Day 8's focused uncached race/shuffle tests passed in 1.418s. Garden and shared reader-error tests passed in 3.660s; the existing personal garden checks also passed before the new tests. Staticcheck passed for the complete 2023 package after both changes. These are test durations, not performance measurements. All checks used Go 1.27.3, Staticcheck 2026.2.2 and temporary writable caches.

Gopls was unavailable; scoped caller inspection, compilation, tests and Staticcheck supplied verification. The critic reviewed production, callers, tests and this record, finding no material issues. Repository-wide checks remain scheduled after the remaining corrections.
