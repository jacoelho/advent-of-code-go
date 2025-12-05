package aoc2025

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`3-5
10-14
16-20
12-18

1
5
8
11
17
32`),
			Want: "3",
		},
		{
			Input: aoc.FileInput(t, 2025, 05),
			Want:  "674",
		},
	}
	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`3-5
10-14
16-20
12-18

1
5
8
11
17
32`),
			Want:  "14",
		},
		{
			Input: aoc.FileInput(t, 2025, 05),
			Want:  "352509891817881",
		},
	}
	aoc.AOCTest(t, day05p02, tests)
}
