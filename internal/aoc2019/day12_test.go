package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 12),
			Want:  "8362",
		},
	}
	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 12),
			Want:  "478373365921244",
		},
	}
	aoc.AOCTest(t, day12p02, tests)
}
