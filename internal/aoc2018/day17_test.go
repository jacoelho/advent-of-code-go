package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("x=495, y=2..7\ny=7, x=495..501\nx=501, y=3..7\nx=498, y=2..4\nx=506, y=1..2\nx=498, y=10..13\nx=504, y=10..13\ny=13, x=498..504"),
			Want:  "57",
		},
		{
			Input: aoc.FileInput(t, 2018, 17),
			Want:  "39367",
		},
	}
	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("x=495, y=2..7\ny=7, x=495..501\nx=501, y=3..7\nx=498, y=2..4\nx=506, y=1..2\nx=498, y=10..13\nx=504, y=10..13\ny=13, x=498..504"),
			Want:  "29",
		},
		{
			Input: aoc.FileInput(t, 2018, 17),
			Want:  "33061",
		},
	}
	aoc.AOCTest(t, day17p02, tests)
}
