package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestDay25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 25),
			Want:  "2214608912",
		},
	}
	aoc.AOCTest(t, day25p01, tests)
}
