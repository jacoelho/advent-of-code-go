# 2019 Intcode corrections

The shared VM and its painting, arcade, droid and beam consumers now preserve failures and compute their results from execution. This resolves F1, F2, F5, F9 and the VM-consumer portion of F10 in the [2019/2020 review](2019-2020.md). Other F10 parsers are tracked there separately.

## Ownership and contracts

- `ParseIntcodeInput` trims surrounding whitespace on comma-separated words and uses the existing checked `convert.ScanNumber` helper. Empty words, malformed numbers and reader failures return errors. Input is limited to 16 MiB and one million words.
- The VM reads every instruction and parameter cell through its zero-extending memory owner. Unwritten nonnegative addresses read zero without allocation. Negative addresses and relative-address overflow fail. Writes are limited to one million words, buffered output to one million values, and each `Run` to ten million instructions. Trusted program/input setters retain their signatures; `Run` checks admitted memory/input sizes. Arithmetic opcodes retain machine-integer semantics.
- `New` and `SetInput` snapshot caller slices. `Clone` retains the full continuation and independent mutable buffers. Droid branches use this owner, preserving the instruction pointer and relative base rather than restarting a copied memory image. Exploration requires one movement reply followed by an input request, rejects invalid/missing/multiple oxygen results, and stops after 100,000 visited tiles.
- Painting keeps sensor color and paint history in one panel record. A seed color is not a paint event; final pairs are processed even when the machine halts. White pixels determine the display bounds, and shared separated-glyph OCR now supports cropped six-row letters. Rendering is limited to one million pixels; simulation to 100,000 commands.
- Arcade output has one validated triple owner. Tile updates replace prior values at a coordinate; block counting uses the final board. Scores remain distinct from tile IDs. Limits are 100,000 input requests and tile positions, one million updates, and coordinates up to one million.
- Beam queries require one binary result followed by halt. Errors flow through both answer functions and the existing `search.BinarySearch` helper. Its callback and result now carry errors; invalid intervals and absent matches fail explicitly. Inclusive midpoint calculation supports the full signed/unsigned range without overflow. The beam search retains its existing widening-beam assumptions and finite coordinate bounds; the first possible bottom row is correctly 99 for a 100-row square. This does not establish arbitrary beam geometry support.

The useful reference ideas remain direct position/direction painting, coordinate-indexed arcade updates and complete VM continuations. No additional interpreter or protocol framework was introduced. The reference revision and source paths remain in the original review.

## Verification

QA assessed fixtures before test changes. New assertions cover final LF/CRLF and malformed comma words, zero-extended instruction parameters, negative/overflowed addresses, the official relative-base quine, large output, owned inputs and independent suspended clones. A self-modifying two-room droid program pauses with a nonzero relative base and produces independently known distances of one for both parts.

An independently traced 19-command painting program yields `I` and 18 painted positions. Arcade regressions cover cleared/repeated blocks and joystick input echoed as score. Both consumers reject invalid protocol output. Beam cases cover an all-pulled field, no solution, and invalid VM replies. Shared-reader checks preserve an injected error with `errors.Is`. Binary-search tests cover endpoints, full integer ranges and predicate errors without asserting probe order.

`go test -race -shuffle=on -count=1 -timeout=2m ./pkg/search ./internal/aoc ./internal/aoc2019` passed: 1.660 s, 1.374 s and 13.401 s. Scoped Staticcheck passed. The named completion critic found no material issues after inspecting the final changes and this verification evidence. Existing personal answers remain unchanged. These are correctness checks, not performance measurements; no additional speed claim is made.
