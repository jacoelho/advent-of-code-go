package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day08p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`b inc 5 if a > 1
a inc 1 if b < 5
c dec -10 if a >= 1
c inc -20 if c == 10`), Want: "1"},
		{Input: strings.NewReader("a dec 1 if a == 0"), Want: "-1"},
		{Input: strings.NewReader("a dec 1 if b == 0"), Want: "0"},
		{Input: aoc.FileInput(t, 2017, 8), Want: "5849"},
	}

	aoc.AOCTest(t, day08p01, tests)
}

func Test_day08p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`b inc 5 if a > 1
a inc 1 if b < 5
c dec -10 if a >= 1
c inc -20 if c == 10`), Want: "10"},
		{Input: strings.NewReader("a dec 1 if a == 0"), Want: "0"},
		{Input: aoc.FileInput(t, 2017, 8), Want: "6702"},
	}

	aoc.AOCTest(t, day08p02, tests)
}
