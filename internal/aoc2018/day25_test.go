package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("0,0,0,0\n3,0,0,0\n0,3,0,0\n0,0,3,0\n0,0,0,3\n0,0,0,6\n9,0,0,0\n12,0,0,0"),
			Want:  "2",
		},
		{
			Input: strings.NewReader("0,0,0,0\n0,0,0,0\n3,0,0,0\n6,0,0,0\n10,0,0,0"),
			Want:  "2",
		},
		{
			Input: aoc.FileInput(t, 2018, 25),
			Want:  "383",
		},
	}
	aoc.AOCTest(t, day25p01, tests)
}
