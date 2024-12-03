package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`),
			Want:  "161",
		},
		{
			Input: aoc.FileInput(t, 2024, 3),
			Want:  "173517243",
		},
	}

	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`),
			Want:  "48",
		},
		{
			Input: aoc.FileInput(t, 2024, 3),
			Want:  "100450138",
		},
	}

	aoc.AOCTest(t, day03p02, tests)
}
