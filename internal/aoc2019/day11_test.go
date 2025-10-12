package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 11),
			Want:  "2041",
		},
	}
	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 11),
			Want:  "ZRZPKEZR",
		},
	}
	aoc.AOCTest(t, day11p02, tests)
}
