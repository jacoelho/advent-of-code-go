package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day18p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(".#.#...|#.\n.....#|##|\n.|..|...#.\n..|#.....#\n#.#|||#|#|\n...#.||...\n.|....|...\n||...#|.#|\n|.||||..|.\n...#.|..|."),
			Want:  "1147",
		},
		{
			Input: strings.NewReader("."),
			Want:  "0",
		},
		{
			Input: strings.NewReader("#"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 18),
			Want:  "564375",
		},
	}
	aoc.AOCTest(t, day18p01, tests)
}

func Test_day18p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("."),
			Want:  "0",
		},
		{
			Input: strings.NewReader("#"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 18),
			Want:  "189720",
		},
	}
	aoc.AOCTest(t, day18p02, tests)
}
