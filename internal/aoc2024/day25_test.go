package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`#####
.####
.####
.####
.#.#.
.#...
.....

#####
##.##
.#.##
...##
...#.
...#.
.....

.....
#....
#....
#...#
#.#.#
#.###
#####

.....
.....
#.#..
###..
###.#
###.#
#####

.....
.....
.....
#....
#.#..
#.#.#
#####`),
			Want: "3",
		},
		{
			Input: aoc.FileInput(t, 2024, 25),
			Want:  "3525",
		},
	}
	aoc.AOCTest(t, day25p01, tests)
}

func TestSchematicLineEndings(t *testing.T) {
	const lock, key = "#####\n#####\n.....\n.....\n.....\n.....\n.....", ".....\n.....\n.....\n.....\n.....\n#####\n#####"
	for _, input := range []string{lock + "\n\n" + key + "\n", strings.ReplaceAll(lock+"\n\n"+key+"\n", "\n", "\r\n")} {
		got, err := day25p01(strings.NewReader(input))
		if err != nil || got != "1" {
			t.Errorf("input %q: got %q, %v; want 1", input, got, err)
		}
	}
}

func TestSchematicInvalidInput(t *testing.T) {
	for _, tt := range []struct{ name, input string }{
		{"empty", ""},
		{"six rows", "#####\n.....\n.....\n.....\n.....\n....."},
		{"short row", "#####\n.....\n....\n.....\n.....\n.....\n....."},
		{"invalid character", "#####\n..x..\n.....\n.....\n.....\n.....\n....."},
		{"neither lock nor key", "####.\n.....\n.....\n.....\n.....\n.....\n....#"},
		{"both full rows", "#####\n.....\n.....\n.....\n.....\n.....\n#####"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := day25p01(strings.NewReader(tt.input)); err == nil {
				t.Fatalf("got %q, want schematic error", got)
			}
		})
	}
}

func TestDay25ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day25p01, "#####\n.####\n.####\n.####\n.#.#.\n.#...\n.....\n"},
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
