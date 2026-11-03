# 2025 tile rectangles and gift packing

Day 9 now checks every included integer tile. Day 12 retains gift geometry and distinguishes a proved fit, a proved failure, and an unsupported or exhausted search. These changes close F1/F2 and the Day 12 final-newline boundary in the [2024–2025 review](2024-2025.md). The required critic found no material issues in the final code, tests, contracts or supplied verification evidence.

## Rectangles

Part 2 compresses each axis into vertex-coordinate singletons and representatives of nonempty integer gaps. Color is constant within each Cartesian block. For each representative row, paired half-open vertical crossings describe the interior; horizontal boundary spans and inclusive vertical boundary points restore the boundary tiles. Merged integer intervals populate a prefix count of uncolored blocks. A candidate is contained exactly when that count is zero. Area uses original coordinates, so compression does not change inclusive tile counts.

This supersedes the original review's proposed continuous segment predicate. The wide U-shaped notch has answers 25/10, but a notch between adjacent integer columns excludes no tiles and has answers 16/16. Rejecting every segment entering the continuous rectangle would reject that valid narrow-notch result. Single-row, single-column and single-tile queries use the same discrete rule.

Parsing accepts 2–1024 distinct signed-int64 points; part 2 additionally requires a simple closed orthogonal loop with at least four vertices. Nonadjacent edges cannot touch, cross or overlap. Forward collinear red vertices remain candidates. Compression is limited to four million cells. Inclusive spans and products are checked before returning an int64 area; part 2 checks containment first. Numeric and reader errors retain their causes.

QA independently checked the tile contract and the actual answers, 4743645488 / 1529011204. The winning part-2 corners are (5106,67466) and (94699,50401). Tests compare compressed containment with direct integer-tile enumeration and ray classification for every red-corner pair in small literal polygons, including reversal, translation, transposition and degenerate rectangles. Invalid loops and overflowing areas return errors.

The [primary statement](https://adventofcode.com/2025/day/9) supplies part 1. Complete part-2 text was checked in czechbol/advent-of-code, `2025/09.livemd`, blob `13302b1b8c3024698a14c95a9be108fb9419718e`, and independently in zertovitch/hac, `exm/aoc/2025/aoc_2025_09_questions.txt`, blob `7c7a34d2f32f89a0126d1fdac342a16ac0802d69`. The selected NickyMeuleman/scrapyard reference at `1d6d8502daa017f4885c3b3110ea9ff2a3c62de1` uses continuous interval rejection; that shortcut was not copied.

## Packing

Occupied cells are normalized and their rotations/reflections deduplicated. Excess occupied area or a piece with no fitting orientation proves failure. A regular array of disjoint boxes, each large enough for any requested piece in a suitable orientation, proves a fit. These certificates run before the exact-search size limit.

Otherwise, regions of at most 64 cells use exact placement masks. A state consists of occupied cells and remaining shape counts. Search chooses the shape with the fewest legal placements and explores every placement. Failed-state memoization depends only on that complete state. Per-shape placement counts and unions, and the union across all remaining shapes, provide necessary-condition pruning. No ordering restriction excludes interchangeable-piece arrangements.

The public invocation permits at most 16 nonempty shapes, each drawn within 32×32 cells, 10000 regions, dimensions/counts at most one million, and 4 MiB of input. Exact search shares a budget of 100 million placement/state checks and one million failed states across regions. Uncertified regions larger than 64 cells and exhausted budgets return errors with region context; neither becomes an impossible-packing answer. Unused region cells are permitted. Reader errors and conventional LF/CRLF endings are handled at parsing.

The [official example](https://adventofcode.com/2025/day/12) produces 2: its individual regions are feasible, feasible, impossible. Tests also cover a shape whose dimensions cannot fit, rotation, reflection, interlocking shapes, empty demand, padding and large certified regions. An independent set-of-cells oracle with literal orientations checks small boards and shape counts without reusing production normalization, masks, memoization or pruning. QA assessed these expectations and resource-error cases before tests were written; the advisor approved the certificate/search boundary and state model.

The actual answer remains 505. Independent inspection found 495 regions rejected by occupied area and 505 accepted by disjoint 3×3 boxes; none requires exact search. NickyMeuleman/scrapyard's Day 12 implementation at the same pinned revision uses this sufficient positive certificate. The old Go area-only answer was removed rather than retained as another solver.

## Verification

Focused Day 9 tests and the shared 2025 immediate/partial reader-error table passed with race detection and randomized order. Focused Day 12 tests, including the official impossible region and independent small-board oracle, passed with the same settings. Staticcheck passed for the complete 2025 package. No benchmark or performance claim accompanies these corrections.
