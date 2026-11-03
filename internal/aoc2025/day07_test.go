package aoc2025

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`.......S.......
...............
.......^.......
...............
......^.^......
...............
.....^.^.^.....
...............
....^.^...^....
...............
...^.^...^.^...
...............
..^...^.....^..
...............
.^.^.^.^.^...^.
...............`),
			Want: "21",
		},
		{
			Input: aoc.FileInput(t, 2025, 07),
			Want:  "1533",
		},
	}
	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`.......S.......
...............
.......^.......
...............
......^.^......
...............
.....^.^.^.....
...............
....^.^...^....
...............
...^.^...^.^...
...............
..^...^.....^..
...............
.^.^.^.^.^...^.
...............`),
			Want: "40",
		},
		{
			Input: aoc.FileInput(t, 2025, 07),
			Want:  "10733529153890",
		},
	}
	aoc.AOCTest(t, day07p02, tests)
}

func TestManifoldExits(t *testing.T) {
	for _, tt := range []struct{ input, splits, timelines string }{
		{"S", "0", "1"},
		{"S\n^", "1", "2"},
		{"S.\n^.", "1", "2"},
		{"S.\n^.\n.^", "2", "3"},
		{"..S..\n..^..\n.^.^.\n..^..\n.....", "4", "6"},
		{".S.\n" + strings.Repeat(".^.\n^.^\n", 64), "192", "55340232221128654846"},
	} {
		for i, solve := range []func(io.Reader) (string, error){day07p01, day07p02} {
			want := [2]string{tt.splits, tt.timelines}[i]
			got, err := solve(strings.NewReader(tt.input))
			if err != nil || got != want {
				t.Fatalf("input %q, part %d = %q, %v; want %s", tt.input, i+1, got, err, want)
			}
		}
	}
}

func TestManifoldInputErrors(t *testing.T) {
	for _, input := range []string{"", ".", "SS", "S.\n.", "S\nx"} {
		for _, solve := range []func(io.Reader) (string, error){day07p01, day07p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected error", input)
			}
		}
	}
}
