# 2024 Day 24 exact circuit repair

Day 24 now computes and certifies four disjoint output swaps. This closes F3 and the circuit portions of F11/F12 in the [2024–2025 review](2024-2025.md). The required critic accepted the final implementation, tests and proof. Shared immediate/partial reader-error cases are integrated and pass.

## General circuit boundary

The parser owns one representation: initial boolean wires plus concrete binary gates. It validates names, 0/1 values, AND/OR/XOR operations, the arrow, one producer per wire (including initial values), and references to defined wires. Strict convert.ScanNumber preserves numeric error causes. LF, final newline, and CRLF are accepted. Input is capped at 4 MiB, 4096 initial wires, 4096 gates, and 64-character alphanumeric wire names.

Part 1 retains general finite circuit evaluation. Each pass either evaluates at least one pending gate or returns a dependency-cycle error; undefined references are rejected earlier. Evaluation copies the initial state and pending gate slice, preserving the parsed circuit. Output z-bit indices must be unique and in 0–4095; math/big assembles the answer without fixed-width shift overflow. The repaired 64-bit adder's 65th output bit therefore fits the same boundary.

Part 2 does not require the corrupted circuit to be acyclic: a swapped output can create a dependency cycle that a valid repair removes.

## Supported repair class and certificate

The named advisor approved canonical ripple-adder role matching, not greedy first-mismatch repairs. The supported class has 3–64 x/y input bits, exactly 5N−3 gates, unique commutative gate signatures, contiguous x00..x(N−1), y00..y(N−1), and z00..zN outputs, no extra z outputs, and all gates consumed by the canonical half-adder/full-adder roles. Other shapes return unsupported errors rather than a heuristic wire list.

Primary propagate/generate gates are fixed by their x/y signatures. The search fixes the first XOR to z00 and tracks the first AND as carry. At each subsequent bit it enumerates every unused nonprimary XOR/AND pair sharing unordered inputs, both assignments of the bit's propagate and incoming carry to those labels, every unused OR gate, and both assignments of generate/half-carry outputs to that OR's inputs. It fixes the sum output to the corresponding z bit and recurses with the chosen OR as carry. The final carry must produce zN.

Output assignments are one constrained involution: each label is unchanged or belongs to exactly one transposition. A demanded assignment either fixes an unchanged label or binds both ends of a swap; later assignments must agree. Four transpositions are required exactly, and every gate/output is assigned. Branches are rolled back with a local mutation journal; no search state escapes except a copied repaired witness.

Every canonical repair has the primary gates, a matching sum/half-carry pair, and an OR at each stage; both commutative orientations are enumerated. Thus this search is complete for the declared class within its explicit work limit. It does not assume that first observations, local operation warnings, or current input values identify the answer.

An independent final structural walk checks the candidate mapping again: bit 0 produces XOR and AND; each subsequent stage produces propagate XOR, generate AND, sum XOR with carry, half-carry AND, and next carry OR; the final carry is zN. Every gate must be consumed exactly once. Induction on those relations proves z=x+y for every assignment of the input bits. Sampled arithmetic alone is not treated as that proof.

## Safe pruning, ambiguity, and limits

The immutable consumer roles classify output labels as sum sinks, inputs to one paired XOR/AND, or inputs to one OR. A producer whose current output label cannot serve its necessary role must move in every valid repair. More than eight such labels proves that four swaps are insufficient. Exactly eight force every other label to remain fixed. During search, the number of unassigned mandatory moved labels cannot exceed twice the unused pair allowance. These are necessary conditions, not a heuristic suspect whitelist.

Every search-node and branch-candidate check shares a 2000000-unit allowance; final certificate work is also charged. Exceeding the allowance returns an explicit resource error even if a candidate has already been found. Search never returns a partial answer as a uniqueness claim.

Completed mappings are deduplicated by canonical transposition set. Different pairings that yield the same sorted eight-wire public answer are acceptable, with one repaired witness retained. A second valid eight-wire answer set is an ambiguity error. Exhaustive successful completion establishes uniqueness of the public answer within the supported class. With one answer set, at most its 105 four-pair matchings can remain in the completion cache; a different answer terminates immediately.

## Primary and local references

Read the [primary Day 24 statement](https://adventofcode.com/2024/day/24) and directly fetched the [full statement mirror](https://github.com/zertovitch/hac/blob/master/exm/aoc/2024/aoc_2024_24_questions.txt), file blob `deb636a34d19cbb38e1d48d15bb3f9776c3555e3`. Part 2 requires correct addition for every input assignment and exactly four disjoint pairs of swapped gate outputs.

Inspected `reference/.repos/zelark--AoC/src/zelark/aoc_2024/day_24.clj` at verified revision `432c59f6ccc52aae7da46cdd8ca0b3ffa8f8aede`. It supplies manual personal swaps and checks one numeric instance; those are not a repair algorithm or universal addition proof. Its explicit half-adder relations are consistent with the supported structural model. No manual swap names were copied into production.

Caller inspection found no other consumers of the replaced private parser/generic load/apply/run/output functions or warning helpers. Go skills and the repository's 2023/2024 layout baseline were applied. gopls was unavailable; scoped source inspection, compilation, race tests, and Staticcheck supply verification.

## Independent evidence and tests

Named QA assessed the initial coherent batch before edits. Existing part-1 examples remain 4 and 2024, and the independently evaluated actual result remains 51107420031718. The former part-2 warning-list expectations on non-adder examples are replaced with unsupported-class assertions.

QA independently tried all 105 pairings of the eight personal golden names and found exactly one complete 222-gate/45-bit structural match: (cpm,krs), (ghp,z33), (gpr,z10), (nks,z21). This proves addition for all 45-bit inputs for that witness and verifies the sorted answer `cpm,ghp,gpr,krs,nks,z10,z21,z33`. QA's evidence alone does not assert uniqueness among arbitrary other eight-wire sets; the production necessary-role pruning and complete search own that conclusion.

A test-only adder fixture writes explicit independent ripple equations and changes only output labels. The review counterexample swaps z01/z02, z03/z04, z05/z06, and z07/z08 in a nine-bit adder; the expected answer is literally z01 through z08. It is checked with all-zero initial values and again with reversed gate order and commutative operands. The answer cannot depend on finding a mismatch in those current input values. An independent three-bit truth table covers all 64 x/y pairs through the part-1 reader boundary, expecting literal x+y.

Parser/evaluator regressions cover final newline and CRLF, invalid bits/operators/arrows, duplicate initial/output producers, initial/output collisions, undefined inputs, dependency cycles, malformed z indices, and oversized z indices. The shared year table exercises immediate and valid-prefix reader errors for both parts.

## Executed verification

All Go commands use GOCACHE=/tmp/aoc-performance-cache, GOMODCACHE=/tmp/aoc-performance-modcache, and STATICCHECK_CACHE=/tmp/aoc-staticcheck-cache.

- Existing part-1 official/personal tests under race/shuffle/count1/timeout2m: passed, package 1.333s.
- Computed actual part-2 repair under the same flags: passed, package 1.316s.
- Nine-bit sum-output counterexample and order variant: passed, package 1.314s.
- Current full owned suite, `go test ./internal/aoc2024 -run '^(Test_day24p|TestCircuit)' -race -shuffle=on -count=1 -timeout=2m`: passed, package 1.231s.
- `/Volumes/External/go/bin/staticcheck ./internal/aoc2024`: passed.
- gofmt and scoped git diff --check: passed.

QA independently enumerated all 51975 four-disjoint-swap mappings over a conventional three-bit adder and proved a unique repair for (z00,c01), (z01,g01), (z02,g02), (z03,p01). The public answer `c01,g01,g02,p01,z00,z01,z02,z03` is asserted, and the returned repaired circuit is evaluated for all 64 x/y pairs against literal x+y.

QA also independently found six distinct eight-wire answer sets for the fixture with (z00,z01), (z02,z03), (g01,h01), (p01,c01). Its test requires the semantic ambiguity error. The unique tiny fixture with zero search allowance requires the distinct resource error; it uses the same real search budget parameter as production, with no separate solver. A literal sparse z64 circuit must return 18446744073709551616, protecting arbitrary-precision output assembly.

Final verification after that supplement:

- `go test ./internal/aoc2024 -run '^(Test_day24p|TestCircuit)' -race -shuffle=on -count=1 -timeout=2m`: passed, package 1.406s.
- Staticcheck2024, gofmt, and scoped git diff --check: passed.

No benchmarks or performance measurements were run. The final shared reader table passed race/shuffle checks after integration; the critic found no material issues.
