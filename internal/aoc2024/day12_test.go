package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`AAAA
BBCD
BBCC
EEEC`),
			Want: "140",
		},
		{
			Input: strings.NewReader(`RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`),
			Want: "1930",
		},
		{
			Input: aoc.FileInput(t, 2024, 12),
			Want:  "1457298",
		},
	}

	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`AAAA
BBCD
BBCC
EEEC`),
			Want: "80",
		},
		{
			Input: strings.NewReader(`EEEEE
EXXXX
EEEEE
EXXXX
EEEEE`),
			Want: "236",
		},
		{
			Input: aoc.FileInput(t, 2024, 12),
			Want:  "921636",
		},
	}

	aoc.AOCTest(t, day12p02, tests)
}

func TestDay12ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day12p01, "A\n"},
		{"part2", day12p02, "A\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				r := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(want))
				if _, err := tt.solve(r); !errors.Is(err, want) {
					t.Errorf("prefix %q: got %v, want reader error", prefix, err)
				}
			}
		})
	}
}
