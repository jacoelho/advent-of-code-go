package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day24p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2023, 24),
			Want:  "26657",
		},
	}
	aoc.AOCTest(t, day24p01, tests)
}

func Test_day24p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`19, 13, 30 @ -2,  1, -2
18, 19, 22 @ -1, -1, -2
20, 25, 34 @ -2, -2, -4
12, 31, 28 @ -1, -2, -1
20, 19, 15 @  1, -5, -3`),
			Want: "47",
		},
		{
			Input: aoc.FileInput(t, 2023, 24),
			Want:  "828418331313365",
		},
	}
	aoc.AOCTest(t, day24p02, tests)
}
