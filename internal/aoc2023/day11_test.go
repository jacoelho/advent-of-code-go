package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....`),
			Want: "374",
		},
		{
			Input: aoc.FileInput(t, 2023, 11),
			Want:  "9795148",
		},
	}
	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....`),
			Want: "82000210",
		},
		{
			Input: aoc.FileInput(t, 2023, 11),
			Want:  "650672493820",
		},
	}
	aoc.AOCTest(t, day11p02, tests)
}
