package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`125 17`),
			Want:  "55312",
		},
		{
			Input: aoc.FileInput(t, 2024, 11),
			Want:  "189092",
		},
	}

	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`125 17`),
			Want:  "65601038650482",
		},
		{
			Input: aoc.FileInput(t, 2024, 11),
			Want:  "224869647102559",
		},
	}

	aoc.AOCTest(t, day11p02, tests)
}
