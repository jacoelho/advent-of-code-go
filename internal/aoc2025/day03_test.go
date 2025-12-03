package aoc2025

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`987654321111111
811111111111119
234234234234278
818181911112111`),
			Want: "357",
		},
		{
			Input: aoc.FileInput(t, 2025, 03),
			Want:  "17445",
		},
	}
	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`987654321111111
811111111111119
234234234234278
818181911112111`),
			Want: "3121910778619",
		},
		{
			Input: aoc.FileInput(t, 2025, 03),
			Want:  "173229689350551",
		},
	}
	aoc.AOCTest(t, day03p02, tests)
}
