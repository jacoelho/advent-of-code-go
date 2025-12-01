package aoc2025

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`),
			Want: "3",
		},
		{
			Input: aoc.FileInput(t, 2025, 01),
			Want:  "1055",
		},
	}
	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`),
			Want: "6",
		},
		{
			Input: aoc.FileInput(t, 2025, 01),
			Want:  "6386",
		},
	}
	aoc.AOCTest(t, day01p02, tests)
}
