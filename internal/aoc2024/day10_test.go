package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`),
			Want: "36",
		},
		{
			Input: aoc.FileInput(t, 2024, 10),
			Want:  "482",
		},
	}

	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`),
			Want: "81",
		},
		{
			Input: aoc.FileInput(t, 2024, 10),
			Want:  "1094",
		},
	}

	aoc.AOCTest(t, day10p02, tests)
}
