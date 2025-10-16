package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 9),
			Want:  "2870072642",
		},
	}
	aoc.AOCTest(t, day9p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 9),
			Want:  "58534",
		},
	}
	aoc.AOCTest(t, day9p02, tests)
}
