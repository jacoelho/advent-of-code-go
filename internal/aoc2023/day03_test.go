package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`),
			Want: "4361",
		},
		{
			Input: aoc.FileInput(t, 2023, 3),
			Want:  "512794",
		},
	}
	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`),
			Want: "467835",
		},
		{
			Input: aoc.FileInput(t, 2023, 3),
			Want:  "67779080",
		},
	}
	aoc.AOCTest(t, day03p02, tests)
}
