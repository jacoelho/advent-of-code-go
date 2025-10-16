package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day15p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`),
			Want:  "1320",
		},
		{
			Input: aoc.FileInput(t, 2023, 15),
			Want:  "516469",
		},
	}
	aoc.AOCTest(t, day15p01, tests)
}

func Test_day15p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`),
			Want:  "145",
		},
		{
			Input: aoc.FileInput(t, 2023, 15),
			Want:  "221627",
		},
	}
	aoc.AOCTest(t, day15p02, tests)
}
