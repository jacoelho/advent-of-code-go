package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`0 <-> 2
1 <-> 1
2 <-> 0, 3, 4
3 <-> 2, 4
4 <-> 2, 3, 6
5 <-> 6
6 <-> 4, 5`), Want: "6"},
		{Input: aoc.FileInput(t, 2017, 12), Want: "169"},
	}

	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`0 <-> 2
1 <-> 1
2 <-> 0, 3, 4
3 <-> 2, 4
4 <-> 2, 3, 6
5 <-> 6
6 <-> 4, 5`), Want: "2"},
		{Input: strings.NewReader("0 <-> 0"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 12), Want: "179"},
	}

	aoc.AOCTest(t, day12p02, tests)
}
