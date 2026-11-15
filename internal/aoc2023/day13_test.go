package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day13p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#`),
			Want: "405",
		},
		{
			Input: aoc.FileInput(t, 2023, 13),
			Want:  "33047",
		},
	}
	aoc.AOCTest(t, day13p01, tests)
}

func Test_day13p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#`),
			Want: "400",
		},
		{
			Input: aoc.FileInput(t, 2023, 13),
			Want:  "28806",
		},
	}
	aoc.AOCTest(t, day13p02, tests)
}

func TestDay13RejectsMalformedPatterns(t *testing.T) {
	requireRejected(t, []string{
		"",
		"##\n#",
		"#x\n#.",
		"##\n##\n\n\n##\n##",
	}, day13p01, day13p02)
	requireRejected(t, []string{"#.#\n..#\n##."}, day13p01)
}
