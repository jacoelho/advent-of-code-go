package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 5),
			Want:  "8332629",
		},
	}
	aoc.AOCTest(t, day5p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 5),
			Want:  "8805067",
		},
	}
	aoc.AOCTest(t, day5p02, tests)
}
