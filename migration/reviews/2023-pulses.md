# 2023 pulse synchronization

This batch resolves the remaining day 20 portion of F3 in the [2023 review](2023.md). Day 8's navigation schedules remain separate; only exact congruence arithmetic is shared.

## Network ownership and simulation

One private network owns flip-flop state, conjunction input memory and ordered destinations. A single FIFO simulator serves both parts and both part-two strategies. It preserves repeated destination occurrences, initializes one conjunction memory bit per sending module and drains each button press before advancing. Undefined destinations are valid sinks. The [module and FIFO rules](https://adventofcode.com/2023/day/20) and the [reproduced part-two statement](https://github.com/AntonKosov/advent-of-code-2023/blob/61cdf160b4dd/day20/part2/README.txt), verified at blob `fe6d09dcad23e273134669cbab28e56728e410bb`, establish the required first LOW to rx.

The former exported module interface, repeated structs/methods, first-HIGH LCM and declaration narration were deleted. Parsing validates unique definitions and names, preserves reader errors, and requires the broadcaster. Module names become stable invocation-local indices; no mutable global state is introduced.

For general topology, complete state keys at drained press boundaries detect cycles. Repetition without an observed LOW proves that rx cannot receive one. Every exact prefix press is observed directly. A non-draining press or exhausted resource limit returns an error; neither yields a partial answer.

## Verified independent branches

The large-answer path requires a unique conjunction G feeding undefined rx, with every G output terminating at an undefined sink. Nothing may feed the broadcaster. Removing the broadcaster and G must leave disjoint stateful branches, each with exactly one broadcaster destination occurrence and one source feeding G. All other defined modules must belong to these branches. Shared undefined sinks are harmless; shared stateful modules, duplicate roots, cross-branch connections and feedback through G invalidate the reduction and select exact whole-state simulation.

For each branch, the cycle key contains every flip-flop bit, every internal conjunction input-memory bit, and that branch's input-memory bit in G. The last bit matters even though G does not feed back: a later press may produce no new pulse before another branch triggers G. A repeat after completed presses a and b makes traces for presses a+1 through b periodic, with period b−a. Each trace stores its initial G-memory bit and every event sent to G with causal depth and pulse level.

FIFO processes causal depths in increasing order. At any depth, independent branches remain grouped in original broadcaster destination order: their children remain in the same branch or terminate at sinks. Induction establishes this ordering even with internal feedback and multiple pulses.

Every LOW from G must be triggered by an incoming HIGH. For each distinct triggering branch/depth pair, the solver computes every qualifying cycle phase. The trigger branch must send a HIGH at that depth. For another branch, G's memory at that instant is its last event at a smaller depth, also including equal-depth events when its root precedes the trigger; otherwise it is the trace's initial memory. Equal-depth trigger events are contiguous relative to other branches, so any HIGH at that depth suffices. Merely starting a press with all memories HIGH does not emit a pulse.

Generalized Chinese remaindering intersects the complete allowed residue sets. Every candidate is lifted strictly beyond the fully simulated prefix, and the minimum across all cuts is returned. Empty intersections mean that cut cannot trigger LOW; a resource error on any unfinished cut aborts the entire search. The private `combineCongruences` and inclusive `firstCongruentAtLeast` helpers now own arithmetic used by both navigation and pulses. They use arbitrary-precision integers and do not mutate inputs. Navigation passes its cycle-entry lower bound; pulses pass the last simulated press plus one.

The advisor checked this topology, event-order and cycle proof before implementation. The collected CameronAavik solver at `1f1313a67b959e21509300a54925ad0ec7a2a23b` directly decodes input-specific counters. That shortcut was rejected because topology alone does not establish its arithmetic or pulse ordering. No separate counter-recognition path was added.

## Bounds and evidence

Input is limited to 4 MiB, 4096 modules/sinks, 16384 edges and 64-character names. Simulation permits 100000 presses, one million events per press, ten million total events and 100000 queued events. Cycle storage has a 64 MiB accounting allowance for encoded keys and trace records, including per-record charges; this is not a measured Go heap limit. Network size, slice capacities and map metadata are additionally bounded by the declared counts. Alignment permits ten million aggregate phase/event/pair checks, 10000 distinct event cuts and 10000 surviving residues. Exhaustion is explicit. Combined answer magnitude is not limited to int64.

An independent Python FIFO probe confirmed the personal topology and cycles. Broadcaster roots `js,ng,lb,gr` lead to four disjoint branches, feeding G through `ss,fh,mf,fz`. All cycles enter at boundary 1, with periods 3881, 3847, 3761 and 3793. Each has one HIGH-bearing phase at depth 4. After 3882 exact presses, ten compatible-pair checks suffice to produce 212986464842911. These are correctness/work-size observations, not benchmark measurements.

QA independently used a separate literal FIFO simulation before new tests were written. Regressions include remembered HIGH state (first LOW at 3), duplicated roots whose order changes the answer from 2 to 1, multiple rx predecessors, absent/never-LOW rx, invalid definitions, reader failures and a non-draining press stopped by a small event budget. A seven-module network has complete branch cycles by press 5 and first LOW at 7, establishing a small independently verified case beyond the observed prefix. Literal trace cases produce 4 versus 3 when root order changes, protecting equal-depth memory and trigger semantics. Shared arithmetic tests cover incompatible and dividing moduli, inclusive lifting and independent result ownership.

Existing sample and personal results remain unchanged: part-one samples 32000000/11687500, personal 787056720, and part two 212986464842911. Existing navigation tests also pass after arithmetic extraction. The focused uncached race/shuffle batch passed in 1.671s; the final full 2023 package passed in 15.848s alongside the 2024 package in 8.938s. Staticcheck and whitespace checks passed. Verification used Go 1.27.3 and Staticcheck 2026.2.2 with temporary writable caches; gopls was unavailable.

The critic reviewed production, callers, tests and this record, finding no material issues. Repository-wide checks remain after the remaining 2024–2025 corrections and final style pass.
