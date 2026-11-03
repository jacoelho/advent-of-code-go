# Ported-package style consistency review

The ported packages broadly match the intended 2023/2024 conventions. No material naming, accidental public API, comment, helper-reuse, or file-organization problem was identified in this bounded pass. The five phase-spacing improvements below have been applied. Their non-whitespace content was verified unchanged.

Scope: /Volumes/External/Projects/advent-of-code-go/internal/aoc2015 through aoc2018, aoc2020 days 01–19, and aoc2021–aoc2022. The reviewer performed no Go execution, tests, builds, benchmarks, reference execution or repository edits. The integrating agent subsequently applied the recommended blank lines and checked their non-whitespace equivalence. Active 2019 work and 2020 days 20–25 were excluded.

Basis: [migration/README.md](/Volumes/External/Projects/advent-of-code-go/migration/README.md), root CLAUDE.md, the user-supplied AGENTS instructions, and representative [2023/day01.go](/Volumes/External/Projects/advent-of-code-go/internal/aoc2023/day01.go), [2023/day12.go](/Volumes/External/Projects/advent-of-code-go/internal/aoc2023/day12.go), [2024/day01.go](/Volumes/External/Projects/advent-of-code-go/internal/aoc2024/day01.go), and [2024/day16.go](/Volumes/External/Projects/advent-of-code-go/internal/aoc2024/day16.go). These establish domain names, private day-local helpers, explicit parsing/solving/formatting blocks, standard import groups, and comments explaining non-obvious rules. Their unrelated legacy defects are not style requirements.

## Coverage and limits

The static scan covered file organization, top-level type/function/method declarations, helper references, and all comment lines in 175 production files: 169 day files and six shared capability files. Targeted body inspection concentrated on representative simple solvers, shared helpers, and longer multi-phase functions. This was not a line-by-line reread of every solver, a new correctness review, or a test-quality audit.

| Package scope | Production files scanned | Targeted body inspection |
| --- | ---: | --- |
| 2015, all days | 25 | wrappingSupplies, containerCombinations, medicine parsing/search preparation |
| 2016, all days plus Assembunny | 26 | bathroomCode, Assembunny parser/state/acceleration, cleaningRoute |
| 2017, all days plus assembly/knot hash | 27 | spiralPositions and consumers, dancePrograms, knotRound/knotHash |
| 2018, all days plus Elfcode | 26 | cave combat parser/neighbors/movement/battle, Elfcode validation/execution |
| 2020, days 01–19 | 19 | expense parsing/pair search, passport validation, docking input, ticket parsing/matching |
| 2021, all days plus digit grid | 26 | depthIncreases, commaSeparatedNumbers, crabPositions, digit-grid parsing, burrow moves/search |
| 2022, all days plus input boundary | 26 | puzzleInput, rucksackPriority, hillDistance, rock-chamber trimming, cube-map parsing/folding/wrapping |

Representative boundary-test layout was sampled in 2015/02, 2016/02, 2017/03, 2018/15, 2020/01, 2021/01 and 2022/22. Compact one-line table cases and expanded multiline examples are both readable; no test-table reformat is warranted.

The already verified 513 solver/test names, fixtures, and 66 expanded wrappers were accepted as existing evidence and not inventoried again. Previously completed correctness, performance and reference decisions were not reopened.

## One bounded readability improvement

Longer functions sometimes place separate algorithm phases together without the blank lines used in the baseline packages. Add a blank line at the following existing boundaries; keep the cohesive loops and error checks together. No new comments, helpers, files, or behavioral changes are needed.

| File and symbol | Concrete insertion points in the inspected file |
| --- | --- |
| [2015/day19.go:137](/Volumes/External/Projects/advent-of-code-go/internal/aoc2015/day19.go:137), medicineParseCost | Before line 159 (unit-rule closure setup), 178 (binary-rule lookup), 182 (parse chart), and 219 (final result). These separate grammar construction, closure, chart evaluation, and result lookup. |
| [2016/day24.go:25](/Volumes/External/Projects/advent-of-code-go/internal/aoc2016/day24.go:25), cleaningRoute | Before line 62 (pairwise-distance phase), 85 (subset route DP), and 106 (final route selection). Parsing and validation remain one block. |
| [2020/day16.go:32](/Volumes/External/Projects/advent-of-code-go/internal/aoc2020/day16.go:32), parseTicketNotes; day16p02 at line 162 | Parser: before 41, 68, 82, 90 to distinguish rules, ticket parser, own ticket, and nearby tickets. P2: before 185, 198, 207 to distinguish candidate construction, unique matching, and result aggregation. |
| [2021/day23.go:156](/Volumes/External/Projects/advent-of-code-go/internal/aoc2021/day23.go:156), burrow.moves | Before line 162 (hallway-to-room moves) and 173 (room departures). These are different transition families sharing the local move closure. |
| [2022/day22.go:173](/Volumes/External/Projects/advent-of-code-go/internal/aoc2022/day22.go:173), foldCube | Before line 185 (face-size derivation), 192 (face discovery), 220 (orientation traversal), and 249 (normal uniqueness check). The entire function still owns one capability. |

These are low-priority readability edits, not completion-blocking correctness findings. Line numbers are from the inspected shared worktree; use the named symbols/statements if concurrent edits move them. If applied, inspect a whitespace-insensitive diff and include the files in the existing final formatting/diff check. Additional behavioral tests are unnecessary for blank-line-only changes.

## Retain the existing design and style

- **Comments:** The inspected comments explain proofs or constraints. Keep examples such as [2015/day17.go:29](/Volumes/External/Projects/advent-of-code-go/internal/aoc2015/day17.go:29), which explains descending iteration over physical containers; [2016/assembunny.go:161](/Volumes/External/Projects/advent-of-code-go/internal/aoc2016/assembunny.go:161), which explains rechecking toggled instructions; [2020/day16.go:129](/Volumes/External/Projects/advent-of-code-go/internal/aoc2020/day16.go:129), which explains uniqueness checks; and [2022/day22.go:275](/Volumes/External/Projects/advent-of-code-go/internal/aoc2022/day22.go:275), which explains the cube-edge direction. No declaration/loop narration worth deleting was found. There is no reason to impose a comment count.

- **Names and API:** Declaration inspection found no exported solver-local types, functions, or methods in this scope. Names such as medicineGrammar, caveBattle, ticketNotes, haltingRecurrence, burrow, cubeFrame and blizzardValley identify their domain. Short coordinates, loop indices, raw input/data variables, and scoped receiver methods are understandable in context. No blanket renames are warranted. Tiny predicate/index helpers need not be expanded merely because entrypoint wrappers were expanded.

- **Organization:** Each day retains its capability. The shared Assembunny, assembly, knot-hash, Elfcode, digit-grid and input files have concrete consumers and responsibilities. The largest inspected production file is 324 lines, 2022/day22.go, and its parsing, walking, and cube-folding code all serve the same day. Splitting files solely by line count would obscure that ownership.

- **Parsing helper variants:** [2021/digit_grid.go:19](/Volumes/External/Projects/advent-of-code-go/internal/aoc2021/digit_grid.go:19) owns a buffer/cell bound and rectangular digit validation; [2022/input.go:9](/Volumes/External/Projects/advent-of-code-go/internal/aoc2022/input.go:9) owns bounded whole-input reading and newline normalization. They are not redundant replacements for the existing scanner helper. [2020/day04.go:50](/Volumes/External/Projects/advent-of-code-go/internal/aoc2020/day04.go:50) enforces decimal-only passport fields, a stricter contract than signed number conversion. Keep that local policy. Existing straight numeric parsing already reuses convert.ScanNumber and scanner where suitable.

- **Search/grid variants:** [2016/day24.go:75](/Volumes/External/Projects/advent-of-code-go/internal/aoc2016/day24.go:75) directly reuses search.BFSDistances. [2018/day15.go:77](/Volumes/External/Projects/advent-of-code-go/internal/aoc2018/day15.go:77) deliberately returns reading-order flat indices for combat; replacing it with generic grid neighbor order would obscure the tie-break rule. [2021/day23.go:198](/Volumes/External/Projects/advent-of-code-go/internal/aoc2021/day23.go:198) reuses the heap but owns stored/expanded-state budgets and returns only a cost; generic AStar retains paths and has no equivalent budget contract. These are justified local algorithms, not missed helper reuse.

- **Loops and layout:** Direct loops that combine validation, errors, mutation and accumulation are clearer than callbacks introduced solely for stylistic uniformity. The representative 2017 and 2018 files already separate cohesive phases well. Standard-library slices/maps, integer ranges and existing collection helpers are used without a new puzzle framework.

This pass found no additional material style issue beyond the bounded phase spacing above. Subsequent correctness fixes and the remaining style reconciliation are complete; see the [migration entry point](../README.md) for their evidence and final verification.

