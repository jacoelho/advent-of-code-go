package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`AAAA
BBCD
BBCC
EEEC`),
			Want: "140",
		},
		{
			Input: strings.NewReader(`RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`),
			Want: "1930",
		},
		{
			Input: aoc.FileInput(t, 2024, 12),
			Want:  "1457298",
		},
	}

	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`AAAA
BBCD
BBCC
EEEC`),
			Want: "80",
		},
		{
			Input: strings.NewReader(`EEEEE
EXXXX
EEEEE
EXXXX
EEEEE`),
			Want: "236",
		},
		{
			Input: aoc.FileInput(t, 2024, 12),
			Want:  "921636",
		},
	}

	aoc.AOCTest(t, day12p02, tests)
}
