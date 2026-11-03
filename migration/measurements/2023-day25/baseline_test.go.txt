package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`jqt: rhn xhk nvd
rsh: frs pzl lsr
xhk: hfx
cmg: qnr nvd lhk bvb
rhn: xhk bvb hfx
bvb: xhk hfx
pzl: lsr hfx nvd
qnr: nvd
ntq: jqt hfx bvb xhk
nvd: lhk
lsr: lhk
rzs: qnr cmg lsr rsh
frs: qnr lhk lsr`),
			Want: "54",
		},
		{
			Input: aoc.FileInput(t, 2023, 25),
			Want:  "614655",
		},
	}
	aoc.AOCTest(t, day25p01, tests)
}
