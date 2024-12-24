package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day21p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`029A
980A
179A
456A
379A`),
			Want: "126384",
		},
		{
			Input: aoc.FileInput(t, 2024, 21),
			Want:  "202274",
		},
	}
	aoc.AOCTest(t, day21p01, tests)
}

func Test_day21p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`029A
980A
179A
456A
379A`),
			Want: "154115708116294",
		},
		{
			Input: aoc.FileInput(t, 2024, 21),
			Want:  "245881705840972",
		},
	}
	aoc.AOCTest(t, day21p02, tests)
}
