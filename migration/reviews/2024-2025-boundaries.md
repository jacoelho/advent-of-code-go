# Input boundaries, price windows and route failures

These corrections close F10 and the listed parts of F11/F13 in the [2024–2025 review](2024-2025.md). The required critic found no material issues in the final code, tests, contracts or supplied verification evidence. Puzzle algorithms remain local; shared parsing and search helpers retain their existing ownership.

## Complete price sequence

2024 Day 22 now supplies the initial price plus all 2000 generated prices to the four-change window. Each buyer still contributes only its first occurrence of a sequence, and the existing bounded concurrent map remains unchanged. Reader and numeric errors reach both entrypoints.

QA independently generated the recurrence using multiplication, integer division and modulus. For seed 1, prices at indices 1996–2000 are 6,0,3,2,9. The sequence (-6,3,-1,7) first occurs at the final price and sells for 9. A frozen assertion on that per-buyer result catches the missing window; the personal maximum remains 1405 with or without it and cannot detect the defect. Independent official/personal totals remain 37327623 / 23 and 12759339434 / 1405.

## Worksheet spacing

2025 Day 6 pads missing trailing cells to the maximum row width before extracting problems. Leading spaces and internal columns remain significant. Input needs numeric rows and an operator row; invalid glyphs, missing or repeated operators and empty worksheets return errors before indexing. LF and CRLF have the same line semantics.

Independent fixtures cover `12`, `3`, `+` on separate lines, with both answers 15, and a two-problem layout whose row interpretation totals 35 while its column interpretation totals 69. Padding variants preserve both results. Personal answers remain 3261038365331 / 8342588849093. This correction does not extend the existing fixed-width integer arithmetic domain.

## Route failures and falling bytes

2024 Days 16 and 20 preserve the existing orientation-aware maze search and single-track cheat calculation. Missing starts and failed route searches now return errors. A valid track with zero qualifying cheats still returns zero. Day 16's private state and search helpers now use reindeer terminology.

Day 18 validates exact comma-separated coordinates with `convert.ScanNumber`, inclusive maximum-coordinate dimensions 0–1024, at most 100000 falling-byte events, coordinates inside the grid, and prefix counts inside the input. Start and destination must both be uncorrupted. Duplicate coordinates remain separate events, preserving input-prefix indices.

Part 2 requires the configured low prefix to be reachable. The existing `search.BinarySearch` finds the first unreachable prefix in `[low+1,len(events)]`; the answer is the event at that prefix minus one. No blocking prefix, an already-blocked low prefix, and invalid configuration are errors. The private route helper uses Manhattan distance on the bounded grid and the existing A* implementation.

QA checked the boundary expectations before test edits: blocked start and destination; the first or last event blocking a 2×2 grid; no blocker; a duplicated-event plateau; a zero-dimensional single-cell grid; invalid prefix/configuration/coordinate input. Day 16's reversed `ES` route costs 2001 and visits two unique tiles. Day 20's valid `SE` route has zero cheats at a positive saving threshold. Both distinguish those results from `S#E`, which has no route.

## Remaining 2024 readers

The remaining Day 1/2/4/6/10/11/12/15/19/21/23/25 entrypoints now return read/conversion errors. Days 1/2/11/21 reuse `convert.ScanNumber`. Day 2 consumes its scanner before inspecting the final error. Original blank reports are invalid; derived reports with fewer than two levels satisfy the adjacent-level constraints vacuously, avoiding the shared window helper's partial-window indexing trap. Day 11 preserves scanner errors and rejects negative stone labels.

Day 6 requires one north-facing guard and runs obstacle trials in a direct loop. The original position/direction cycle state remains authoritative; only one trial's visited set is live at a time. There is no per-obstacle goroutine or measured speed claim.

Days 15 and 19 consume a bounded input before parsing sections, preserving late reader errors and normalizing CRLF. Each has a 128 KiB total-input cap. Warehouse validation requires rectangular enclosing walls, valid tiles/moves and one robot, preventing an unenclosed map from causing an endless walk. Towel patterns must be nonempty. Day 21 creates validated keypad records with a parsed numeric value. Day 23 admits simple non-self edges before calling the existing clique implementation. Other solving algorithms remain unchanged.

The critic identified two integration defects: uncapped whole-input reads and singleton report windows. Both were corrected and accepted in the follow-up review. The shared 2024 table now covers all 49 entrypoints with immediate and valid-prefix reader failures; configured solvers retain their existing factories.

## Verification

Immediate and valid-prefix reader failures use `errors.Is` at all affected entrypoints in the shared year tables. Existing examples and personal goldens pass. Focused counting/worksheet tests passed race detection and randomized order (2024 4.559s, 2025 1.317s); focused route tests and the shared reader table passed under the same settings (4.005s). No performance claim is made.

The remaining 23 entrypoints' existing answers and QA-approved parser regressions also passed race/shuffle checks and Staticcheck. After review fixes, focused warehouse/towel checks passed in 1.402s and short-report/shared-reader checks passed in 1.344s. The required critic accepted the completed reader batch.
