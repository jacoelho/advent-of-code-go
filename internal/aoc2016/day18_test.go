package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day18p01(t *testing.T) {
	aoc.AOCTest(t, day18p01, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 18), Want: "2013"}})
}
func Test_day18p02(t *testing.T) {
	aoc.AOCTest(t, day18p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 18), Want: "20006289"}})
}

func TestSafeTileRows(t *testing.T) {
	for _, tc := range []struct {
		input string
		rows  int
		want  string
	}{{".^^.^.^^^^", 10, "38"}, {"..^^.", 3, "6"}, {"^", 2, "1"}, {".", 3, "3"}} {
		got, err := safeTiles(strings.NewReader(tc.input), tc.rows)
		if err != nil || got != tc.want {
			t.Errorf("%s rows%d = %q,%v; want%s", tc.input, tc.rows, got, err, tc.want)
		}
	}
}
