package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day15p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 15),
			Want:  "216",
		},
	}
	aoc.AOCTest(t, day15p01, tests)
}

func Test_day15p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 15),
			Want:  "326",
		},
	}
	aoc.AOCTest(t, day15p02, tests)
}
