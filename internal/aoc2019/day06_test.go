package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 6),
			Want:  "204521",
		},
	}
	aoc.AOCTest(t, day6p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 6),
			Want:  "307",
		},
	}
	aoc.AOCTest(t, day6p02, tests)
}
