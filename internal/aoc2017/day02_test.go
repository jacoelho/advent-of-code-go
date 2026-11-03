package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`5 1 9 5
7 5 3
2 4 6 8`), Want: "18"},
		{Input: aoc.FileInput(t, 2017, 2), Want: "30994"},
	}

	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`5 9 2 8
9 4 7 3
3 8 6 5`), Want: "9"},
		{Input: strings.NewReader("2 2"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 2), Want: "233"},
	}

	aoc.AOCTest(t, day02p02, tests)
}
