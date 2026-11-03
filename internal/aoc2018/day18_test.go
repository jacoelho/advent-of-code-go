package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day18p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(".#.#...|#.\n.....#|##|\n.|..|...#.\n..|#.....#\n#.#|||#|#|\n...#.||...\n.|....|...\n||...#|.#|\n|.||||..|.\n...#.|..|."),
			Want:  "1147",
		},
		{
			Input: strings.NewReader("."),
			Want:  "0",
		},
		{
			Input: strings.NewReader("#"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 18),
			Want:  "564375",
		},
	}
	aoc.AOCTest(t, day18p01, tests)
}

func Test_day18p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("."),
			Want:  "0",
		},
		{
			Input: strings.NewReader("#"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 18),
			Want:  "189720",
		},
	}
	aoc.AOCTest(t, day18p02, tests)
}

func TestDay18ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"18/1", day18p01},
		{"18/2", day18p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
