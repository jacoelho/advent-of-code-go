package aoc2025

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`),
			Want: "50",
		},
		{
			Input: aoc.FileInput(t, 2025, 9),
			Want:  "4743645488",
		},
	}
	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`),
			Want: "24",
		},
		{
			Input: aoc.FileInput(t, 2025, 9),
			Want:  "1529011204",
		},
	}
	aoc.AOCTest(t, day09p02, tests)
}
