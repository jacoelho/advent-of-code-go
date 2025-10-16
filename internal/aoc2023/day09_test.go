package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45`),
			Want: "114",
		},
		{
			Input: aoc.FileInput(t, 2023, 9),
			Want:  "1938731307",
		},
	}
	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45`),
			Want: "2",
		},
		{
			Input: aoc.FileInput(t, 2023, 9),
			Want:  "948",
		},
	}
	aoc.AOCTest(t, day09p02, tests)
}
