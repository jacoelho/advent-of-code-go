package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Time:      7  15   30
Distance:  9  40  200`),
			Want: "288",
		},
		{
			Input: aoc.FileInput(t, 2023, 6),
			Want:  "1195150",
		},
	}
	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Time:      7  15   30
Distance:  9  40  200`),
			Want: "71503",
		},
		{
			Input: aoc.FileInput(t, 2023, 6),
			Want:  "42550411",
		},
	}
	aoc.AOCTest(t, day06p02, tests)
}
