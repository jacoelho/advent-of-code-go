package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`3   4
4   3
2   5
1   3
3   9
3   3`),
			Want: "11",
		},
		{
			Input: aoc.FileInput(t, 2024, 1),
			Want:  "1110981",
		},
	}

	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`3   4
4   3
2   5
1   3
3   9
3   3`),
			Want: "31",
		},
		{
			Input: aoc.FileInput(t, 2024, 1),
			Want:  "24869388",
		},
	}

	aoc.AOCTest(t, day01p02, tests)
}
