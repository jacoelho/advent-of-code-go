package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day23p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`389125467`),
			Want:  "67384529",
		},
		{
			Input: aoc.FileInput(t, 2020, 23),
			Want:  "43769582",
		},
	}
	aoc.AOCTest(t, day23p01, tests)
}

func Test_day23p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`389125467`),
			Want:  "149245887792",
		},
		{
			Input: aoc.FileInput(t, 2020, 23),
			Want:  "264692662390",
		},
	}
	aoc.AOCTest(t, day23p02, tests)
}
