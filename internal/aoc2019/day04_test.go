package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day04p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 4),
			Want:  "1660",
		},
	}
	aoc.AOCTest(t, day4p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 4),
			Want:  "1135",
		},
	}
	aoc.AOCTest(t, day4p02, tests)
}
