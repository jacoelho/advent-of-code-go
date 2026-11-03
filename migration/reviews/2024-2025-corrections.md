# 2024–2025 ordering, boundaries and counting corrections

This batch resolves F4, F5, F7, F9 and the two contract questions Q1/Q2 from the [2024–2025 review](2024-2025.md). F11 is only partially addressed: reader failures now propagate through the touched 2024 days 3/5/8/9. Subsequent year-wide corrections and final status are linked from the [review](2024-2025.md).

## Multiplication, ordering and disk movement

2024 day 3 recognizes only multiplication operands with one to three ASCII digits, as required by the [instruction grammar](https://adventofcode.com/2024/day/3). Corrupted instruction-like text is ignored. The existing enable/disable state remains sequential. Reads are bounded at 4 MiB; that bound and the operand ceiling keep the int64 sum safe. The shared strict numeric helper replaces panic-based conversion.

Day 5 stores directed rules without losing contradictory edges and topologically orders each invalid update's induced graph. Among available pages, original update order breaks ties. Only rules involving pages present in the update apply, matching the [puzzle contract](https://adventofcode.com/2024/day/5). Part one excludes cyclic updates; part two reports the induced cycle. Parsing checks nonnegative unique pages and an odd update length; numerical errors retain their causes. Bounds are 100000 lines, 1023 pages per update and ten million aggregate page pairs. Middle-page sums use checked int64 arithmetic.

Day 9 coalesces physically adjacent free intervals across zero-length files, preserving their original IDs. Part two retains whole-file movement in descending ID order; vacated space is always to the right of files still under consideration. Part one retains cell compaction. One decimal line with optional final LF/CRLF/CR is validated before either algorithm. Limits are 100000 digits, at most 900000 expanded cells and one hundred million free-span checks; checksums fit int64 under these limits. See the [disk format](https://adventofcode.com/2024/day/9).

QA independently approved regressions before edits: the ordering chain repairs to middle page 2, induced-only rules and deterministic ties have literal answers, and cycles fail. Disk `12023` gives 12/12, `1202023` gives 18/18, and `1200203` gives 38/60. The last two distinguish preserved file IDs and whole-file movement from cell movement. The personal disk has no zero-length files, so its previous golden is unchanged for an independently established reason.

The collected zelark references at `432c59f6ccc52aae7da46cdd8ca0b3ffa8f8aede` repeat the unbounded regex and partial-order sorting assumptions. Those shortcuts were rejected. Their mutable block contents were unnecessary for the existing interval model.

## Dial position and click counts

2025 day 1 now counts zero after each rotation in part one, matching the [statement](https://adventofcode.com/2025/day/1). Part two counts actual arrivals at zero: a zero-distance rotation adds no clicks. Full turns and residual movement are separated, and only movement modulo 100 is added to the current position, avoiding overflow. Parsing accepts unsigned decimal magnitudes through MaxInt64, with L/R owning direction. Input is limited to 100000 rotations and the final selected password is checked for int64 overflow.

QA compared the click helper with literal single-click movement for positions 0..99 and movements −300..300. `R50` gives 1/1, while `R50` followed by `R0` gives 2/1. One hundred maximum-distance right rotations reach exactly MaxInt64 clicks; 101 fail. The NickyMeuleman reference at `1d6d8502daa017f4885c3b3110ea9ff2a3c62de1` corroborates after-rotation counting; its per-click loop is useful as an independent test model, not as the production algorithm.

## Collinear antennas

Q1 is resolved by the [reproduced complete statement](https://github.com/zertovitch/hac/blob/33feba40659058fc2981c07f11b2eaeeeb80e8ae/exm/aoc/2024/aoc_2024_08_questions.txt), verified at blob `13ddd2383ebce28255d1fa8b83f81c491d7a89f8`. Part two includes every lattice position on a line through two same-frequency antennas, regardless of their separation. Dividing their displacement by its coordinate GCD gives the primitive step; the shared GCD helper owns that arithmetic.

The parser now requires a nonempty rectangular ASCII map of alphanumeric frequencies and dots. Frequency remains case-sensitive. Dimensions come from rows rather than occupied cells, and reader errors propagate. Input is bounded at one million cells. Search is bounded at one million pair/step operations; exhaustion returns an error instead of a partial count.

Independent examples include `A.A` (0/3), a five-cell diagonal with antennas two steps apart (1/5), a reversed diagonal (0/3), and `A.a` (0/0). QA independently recomputed the personal part-two result, 1308, by checking determinants for every grid cell rather than stepping along lines. The malformed old part-one sample's 13-character final row was corrected to the published 12 dots; its expected result remains 14.

## Quantum beam exits

Q2 is resolved by the [reproduced complete statement](https://github.com/zertovitch/hac/blob/33feba40659058fc2981c07f11b2eaeeeb80e8ae/exm/aoc/2025/aoc_2025_07_questions.txt), verified at blob `ee894f83929cf88e72221a40a0938cb56289ab80`. Every encountered splitter creates two timeline branches; a branch leaving the manifold is still a completed journey. Part two now evaluates both branches, including side exits. The existing memoized downward-trace model remains; counts use `big.Int`, and addition never mutates a cached value.

Parsing requires a rectangular `.^S` diagram with exactly one start, at most 100000 cells and 4096 rows. This prevents unknown glyphs or duplicate starts from leaving the trace loop without progress. Part one continues to count each reached splitter once, independently of the number of histories reaching it.

Literal tests cover one or both branches exiting, nested exits (2 splits/3 timelines), and converging histories (4/6). A compact repeated diagram gives 192 splits and 55340232221128654846 timelines, beyond int64. QA's independent weighted row propagation preserved the actual 1533/10733529153890 results. Same-row adjacent-splitter interpretation was not established by this batch; the existing behavior of advancing downward from each emitted side position is retained.

## Verification and remaining scope

All existing personal results remain unchanged, including 2024 days 3/5/9 (173517243/100450138, 4996/6311, 6332189866718/6353648390778) and 2025 day 1 (1055/6386). Shared per-year reader tables check immediate and valid-prefix failures with `errors.Is`; redundant dial-only reader tests were removed.

Focused uncached race/shuffle tests passed for the four initial corrections (2024 1.681s; 2025 1.327s), antenna/manifold corrections and shared readers (2024 1.747s; 2025 2.466s), and the expanded 2024 reader table (1.510s). Staticcheck passed for both complete packages. Gofmt and diff whitespace checks passed. These are test durations, not performance measurements. Verification used Go 1.27.3, Staticcheck 2026.2.2 and temporary writable caches; gopls was unavailable.

The critic reviewed the six implementations, their callers, tests and this record, finding no material issues in scope. Full repository checks and the final style reconciliation remain after the other corrections. This batch does not close packing, polygon containment, circuit repair, calibration, claw constraints, final-price generation, remaining reader/search failures or bounded-search findings.
