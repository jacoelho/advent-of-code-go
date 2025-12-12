package aoc2025

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2025, 12),
			Want:  "505",
		},
	}
	aoc.AOCTest(t, day12p01, tests)
}
