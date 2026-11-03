package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const amphipodExample = `#############
#...........#
###B#C#B#D###
  #A#D#C#A#
  #########`

func Test_day23p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(amphipodExample), Want: "12521"},
		{Input: strings.NewReader("#############\n#...........#\n###A#B#C#D###\n  #A#B#C#D#\n  #########"), Want: "0"},
		{Input: strings.NewReader("#############\n#A.........B#\n###.#.#C#D###\n  #A#B#C#D#\n  #########"), Want: "73"},
		{Input: strings.NewReader("#############\n#AB.........#\n###.#.#C#D###\n  #A#B#C#D#\n  #########"), Want: "43"},
		{Input: strings.NewReader("#############\n#...........#\n###B#A#C#D###\n  #A#B#C#D#\n  #########"), Want: "46"},
		{Input: aoc.FileInput(t, 2021, 23), Want: "16157"},
	}
	aoc.AOCTest(t, day23p01, tests)
}

func Test_day23p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(amphipodExample), Want: "44169"},
		{Input: aoc.FileInput(t, 2021, 23), Want: "43481"},
	}
	aoc.AOCTest(t, day23p02, tests)
}

func TestBurrowInvalidInput(t *testing.T) {
	tests := []string{
		"",
		strings.Replace(amphipodExample, "#...........#", "#..A........#", 1),
		strings.Replace(amphipodExample, "#A#D#C#A#", "#.#D#C#A#", 1),
		strings.Replace(amphipodExample, "###B#C#B#D###", "###B#C#B#E###", 1),
		strings.Replace(amphipodExample, "#############", "############", 1),
	}
	for _, input := range tests {
		if _, err := day23p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected invalid burrow error for %q", input)
		}
	}
	input := "#############\n#A..........#\n###.#B#C#D###\n  #A#B#C#D#\n  #########"
	if _, err := day23p02(strings.NewReader(input)); err == nil {
		t.Fatal("unfolding an occupied hallway must fail")
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
