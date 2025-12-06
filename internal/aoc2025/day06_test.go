package aoc2025

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  `),
			Want: "4277556",
		},
		{
			Input: aoc.FileInput(t, 2025, 06),
			Want:  "3261038365331",
		},
	}
	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +`),
			Want: "3263827",
		},
		{
			Input: aoc.FileInput(t, 2025, 06),
			Want:  "8342588849093",
		},
	}
	aoc.AOCTest(t, day06p02, tests)
}
