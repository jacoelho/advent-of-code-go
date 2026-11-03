package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`0
3
0
1
-3`), Want: "5"},
		{Input: strings.NewReader("-1"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 5), Want: "394829"},
	}

	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`0
3
0
1
-3`), Want: "10"},
		{Input: strings.NewReader("-1"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 5), Want: "31150702"},
	}

	aoc.AOCTest(t, day05p02, tests)
}
