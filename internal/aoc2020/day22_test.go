package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Player 1:
9
2
6
3
1

Player 2:
5
8
4
7
10`),
			Want: "306",
		},
		{
			Input: aoc.FileInput(t, 2020, 22),
			Want:  "31455",
		},
	}
	aoc.AOCTest(t, day22p01, tests)
}

func Test_day22p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Player 1:
9
2
6
3
1

Player 2:
5
8
4
7
10`),
			Want: "291",
		},
		{
			Input: aoc.FileInput(t, 2020, 22),
			Want:  "32528",
		},
	}
	aoc.AOCTest(t, day22p02, tests)
}
