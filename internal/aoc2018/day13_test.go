package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day13p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`/->-\
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/`),
			Want: "7,3",
		},
		{
			Input: aoc.FileInput(t, 2018, 13),
			Want:  "57,104",
		},
	}

	aoc.AOCTest(t, day13p01, tests)
}

func Test_day13p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`/>-<\
|   |
| /<+-\
| | | v
\>+</ |
  |   ^
  \<->/`),
			Want: "6,4",
		},
		{
			Input: strings.NewReader(">><"),
			Want:  "1,0",
		},
		{
			Input: aoc.FileInput(t, 2018, 13),
			Want:  "67,74",
		},
	}

	aoc.AOCTest(t, day13p02, tests)
}

func TestMineCartErrors(t *testing.T) {
	for _, input := range []string{"", "><", "?"} {
		if _, err := day13p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}
