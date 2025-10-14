package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day8p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 8),
			Want:  "2684",
		},
	}
	aoc.AOCTest(t, day8p01, tests)
}

func Test_day8p02(t *testing.T) {
	tests := []aoc.TestInput{

		{
			Input: aoc.FileInput(t, 2019, 8),
			Want:  "YGRYZ",
		},
	}
	aoc.AOCTest(t, day8p02, tests)
}
