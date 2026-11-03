package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("0 2 7 0"), Want: "5"},
		{Input: aoc.FileInput(t, 2017, 6), Want: "11137"},
	}

	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("0 2 7 0"), Want: "4"},
		{Input: aoc.FileInput(t, 2017, 6), Want: "1037"},
	}

	aoc.AOCTest(t, day06p02, tests)
}
