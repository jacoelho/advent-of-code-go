# Robot frames and chronospatial quines

These changes complete the remaining Day 14/17 search and reader-boundary corrections from F11/F12 in the [2024–2025 review](2024-2025.md). The advisor approved the supported domains and proofs before implementation; QA independently assessed the tests. The required critic found no material issues in the final code, tests, supported domains or supplied verification evidence.

## Robot frames

Integer robot motion on the 101×103 torus repeats after 10403 steps. The candidate finder examines `[0,10403)`, including the original frame. More robots than cells proves that collision freedom is impossible. No candidate or multiple candidates returns an identification error; the old unbounded wait is removed.

Part 2 retains an explicit specialization: the unique collision-free frame is assumed to be the intended tree. Uniqueness alone does not prove tree geometry. A full board constructed from stationary robots and one moving robot is a counterexample, so no generic recognition claim is made. For the personal input, an independent Python enumeration of all 10403 times found exactly one candidate, 7861. The [saved 101×103 raster](2024-robots-7861.txt), with trailing blank cells omitted, visibly contains the framed tree; the primary agent and advisor both inspected it.

The parser uses checked coordinate conversion and preserves reader errors. It admits at most 100000 robots. Part-1 dimensions must be between 1 and one million. Positions and velocities are reduced modulo each dimension before motion, avoiding overflowing the original coordinate/velocity products.

A separate candidate-finder test uses a 2×3 torus with five stationary robots and one moving robot; independent enumeration proves that only time zero is collision-free. It makes no tree claim. Identical stationary robots, one stationary robot and excess robot count exercise absent, ambiguous and impossible candidates. Existing part-1 examples and personal answers 229868730 / 7861 remain unchanged.

## Chronospatial execution

One value-owned interpreter implements byte-index jumps, including odd targets. Fetching beyond the program halts; an in-range opcode without an operand fails. Literal or ignored operand 7 is valid, while combo operand 7 is reserved. Nonnegative registers are limited to signed-int64 input values; shifts use unsigned values, including shift counts at least 64. Reader and numeric causes propagate through strict field parsing. The input limit is 64 KiB and 1024 three-bit program words.

Part 1 has a one-million-instruction limit. Part 2 shares 100 million instruction steps and one million candidate trials across the entire search. Exhaustion returns an error, never a partial or guessed answer. Output storage is bounded by executed instructions. The unused disassembler, file-wide lint suppression, mutable clone/IP ownership and unbounded greedy search were removed.

## Supported quine class

The quine solver requires a straight loop with exactly one division of A by eight, one output and a sole terminal `jnz 0`. A forward def-use check requires every B/C read to follow a definition within that iteration; unused registers need no definition. Reads are checked before writes, including a register used in its own defining instruction. Both the official example and personal program satisfy these rules.

For an incoming A of at least eight, execution emits one value and then has exactly the output of an iteration starting with `A/8`. Previous B/C values are irrelevant. Starting from the suffix seed zero, the search extends every survivor with all eight octal digits and compares its complete output with the corresponding target suffix. This includes every supported quine; retaining all survivors avoids greedy branch loss. The minimum positive complete match is returned.

A register in the supported range has at most 21 octal digits. Longer quine targets, checked-growth failures, unsupported loop structure, exhausted budgets and absence of a quine are distinguished. Unsupported programs such as `5,0` terminate with an error rather than entering the old endless probing loop.

QA independently verified a further supported program, `2,0,0,3,5,4,3,0`, whose minimum quine is 7516176. Literal tests cover valid operand 7, reserved combo 7, odd jumps, jumps past the program, missing operands, large shifts, bounded loops, register read-before-write rejection and a supported loop with no quine. Existing official and personal answers remain unchanged, including personal minimum 37221274271220.

## Verification

All added tests were assessed by the named QA specialist before edits. Focused Day 14/17 tests and the complete 2024 immediate/partial reader-error table passed race detection and randomized order (1.722s). Staticcheck passed for 2024. No benchmark or speed claim accompanies the changes.
