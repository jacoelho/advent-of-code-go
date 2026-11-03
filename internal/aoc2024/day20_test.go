package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01_example(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############`),
			Want: "44",
		},
	}
	aoc.AOCTest(t, day20p01(2), tests)
}

func Test_day20p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 20),
			Want:  "1445",
		},
	}
	aoc.AOCTest(t, day20p01(100), tests)
}

func Test_day20p02_example(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############`),
			Want: "285",
		},
	}
	aoc.AOCTest(t, day20p02(50), tests)
}

func Test_day20p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 20),
			Want:  "1008040",
		},
	}
	aoc.AOCTest(t, day20p02(100), tests)
}

func TestRaceRoutes(t *testing.T) {
	for part, solve := range []func(io.Reader) (string, error){day20p01(1), day20p02(1)} {
		for _, input := range []string{"S#E", ".#E"} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("part %d input %q: expected route error", part+1, input)
			}
		}
		got, err := solve(strings.NewReader("SE"))
		if err != nil || got != "0" {
			t.Fatalf("part %d short route = %q, %v; want 0", part+1, got, err)
		}
	}
}

func TestDay20ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day20p01(1), "SE\n"},
		{"part2", day20p02(1), "SE\n"},
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
