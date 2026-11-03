package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const caveExample = `#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######`

func Test_day15p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("EG"),
			Want:  "134",
		},
		{
			Input: strings.NewReader("EGE"),
			Want:  "9933",
		},
		{
			Input: strings.NewReader(caveExample),
			Want:  "27730",
		},
		{
			Input: strings.NewReader(`#######
#G..#E#
#E#E.E#
#G.##.#
#...#E#
#...E.#
#######`),
			Want: "36334",
		},
		{
			Input: strings.NewReader(`#######
#E..EG#
#.#G.E#
#E.##E#
#G..#.#
#..E#.#
#######`),
			Want: "39514",
		},
		{
			Input: strings.NewReader(`#######
#E.G#.#
#.#G..#
#G.#.G#
#G..#.#
#...E.#
#######`),
			Want: "27755",
		},
		{
			Input: strings.NewReader(`#######
#.E...#
#.#..G#
#.###.#
#E#G#G#
#...#G#
#######`),
			Want: "28944",
		},
		{
			Input: strings.NewReader(`#########
#G......#
#.E.#...#
#..##..G#
#...##..#
#...#...#
#.G...G.#
#.....G.#
#########`),
			Want: "18740",
		},
		{
			Input: aoc.FileInput(t, 2018, 15),
			Want:  "319410",
		},
	}

	aoc.AOCTest(t, day15p01, tests)
}

func Test_day15p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("EG"),
			Want:  "134",
		},
		{
			Input: strings.NewReader(caveExample),
			Want:  "4988",
		},
		{
			Input: aoc.FileInput(t, 2018, 15),
			Want:  "63168",
		},
	}

	aoc.AOCTest(t, day15p02, tests)
}

func TestCaveErrors(t *testing.T) {
	for _, input := range []string{"", "E?G", "E.\nG", "E#G"} {
		if _, err := day15p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

func TestDay15ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"15/1", day15p01},
		{"15/2", day15p02},
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
