package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

const elfExample = `....#..
..###.#
#...#.#
.#...##
#.###..
##.#.##
.#..#..`

func Test_day23p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(elfExample), Want: "110"},
		{Input: strings.NewReader("#"), Want: "0"},
		{Input: strings.NewReader("."), Want: "0"},
		{Input: aoc.FileInput(t, 2022, 23), Want: "3762"},
	}
	aoc.AOCTest(t, day23p01, tests)
}

func Test_day23p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(elfExample), Want: "20"},
		{Input: strings.NewReader("#"), Want: "1"},
		{Input: strings.NewReader("."), Want: "1"},
		{Input: aoc.FileInput(t, 2022, 23), Want: "997"},
	}
	aoc.AOCTest(t, day23p02, tests)
}

func TestElfConflictingProposals(t *testing.T) {
	elves, err := parseElves(strings.NewReader("#\n#\n.\n#\n#"))
	if err != nil {
		t.Fatal(err)
	}
	if !elfRound(elves, 0) {
		t.Fatal("outer elves should move while the middle proposals conflict")
	}
	want := collections.NewSet(
		grid.Position2D[int]{X: 0, Y: -1},
		grid.Position2D[int]{X: 0, Y: 1},
		grid.Position2D[int]{X: 0, Y: 3},
		grid.Position2D[int]{X: 0, Y: 5},
	)
	if len(elves) != len(want) || len(elves.Difference(want)) != 0 {
		t.Fatalf("got %v; want %v", elves, want)
	}
}

func TestElfInvalidInput(t *testing.T) {
	for _, input := range []string{"", "#\n##", ".x"} {
		if _, err := day23p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected invalid grove error for %q", input)
		}
	}
}

func TestDay23ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day23p01},
		{"part2", day23p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
