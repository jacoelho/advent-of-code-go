package aoc2019

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("R8,U5,L5,D3\nU7,R6,D4,L4"),
			Want:  "6",
		},
		{
			Input: aoc.FileInput(t, 2019, 3),
			Want:  "386",
		},
	}
	aoc.AOCTest(t, day3p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("R8,U5,L5,D3\nU7,R6,D4,L4"),
			Want:  "30",
		},
		{
			Input: aoc.FileInput(t, 2019, 3),
			Want:  "6484",
		},
	}
	aoc.AOCTest(t, day3p02, tests)
}
