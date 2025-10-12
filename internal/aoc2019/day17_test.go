package aoc2019

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p1(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`..#..........
..#..........
#######...###
#.#...#...#.#
#############
..#...#...#..
..#####...^..`),
			Want: "76",
		},
	}
	aoc.AOCTest(t, day17p1, tests)
}

func Test_day17p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 17),
			Want:  "6448",
		},
	}
	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 17),
			Want:  "914900",
		},
	}
	aoc.AOCTest(t, day17p02, tests)
}
