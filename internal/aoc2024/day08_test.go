package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day08p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............`),
			Want: "14",
		},
		{
			Input: aoc.FileInput(t, 2024, 8),
			Want:  "409",
		},
	}

	aoc.AOCTest(t, day08p01, tests)
}

func Test_day08p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............`),
			Want: "34",
		},
		{
			Input: aoc.FileInput(t, 2024, 8),
			Want:  "1308",
		},
	}

	aoc.AOCTest(t, day08p02, tests)
}

func TestAntinodeLines(t *testing.T) {
	for _, tt := range []struct{ input, first, second string }{
		{"A.A", "0", "3"},
		{"A....\n.....\n..A..\n.....\n.....", "1", "5"},
		{"..A\n...\nA..", "0", "3"},
		{"A.a", "0", "0"},
		{"A\n.\nA", "0", "3"},
		{"...\n...", "0", "0"},
	} {
		for i, solve := range []func(io.Reader) (string, error){day08p01, day08p02} {
			want := [2]string{tt.first, tt.second}[i]
			got, err := solve(strings.NewReader(tt.input))
			if err != nil || got != want {
				t.Fatalf("input %q, part %d = %q, %v; want %s", tt.input, i+1, got, err, want)
			}
		}
	}
}

func TestAntennaInputErrors(t *testing.T) {
	for _, input := range []string{"", "AA\n...", "A_", "Aé"} {
		for _, solve := range []func(io.Reader) (string, error){day08p01, day08p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected error", input)
			}
		}
	}
}

func TestDay08ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day08p01, "A.A\n"},
		{"part2", day08p02, "A.A\n"},
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
