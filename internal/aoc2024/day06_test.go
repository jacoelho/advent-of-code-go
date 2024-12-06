package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`),
			Want: "41",
		},
		{
			Input: aoc.FileInput(t, 2024, 6),
			Want:  "5208",
		},
	}

	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`),
			Want: "6",
		},
		{
			Input: aoc.FileInput(t, 2024, 6),
			Want:  "1972",
		},
	}

	aoc.AOCTest(t, day06p02, tests)
}
