package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	aoc.AOCTest(t, day25p01, []aoc.TestInput{
		{Input: strings.NewReader("row 1, column 1"), Want: "20151125"},
		{Input: strings.NewReader("row 2, column 1"), Want: "31916031"},
		{Input: strings.NewReader("row 1, column 2"), Want: "18749137"},
		{Input: strings.NewReader("row 6, column 6"), Want: "27995004"},
		{Input: aoc.FileInput(t, 2015, 25), Want: "9132360"},
	})
	for _, input := range []string{"", "row 0, column 1", "row 1, column 0"} {
		if _, err := day25p01(strings.NewReader(input)); err == nil {
			t.Errorf("accepted %q", input)
		}
	}
}
