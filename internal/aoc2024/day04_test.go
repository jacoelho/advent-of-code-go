package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day04p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`),
			Want: "18",
		},
		{
			Input: aoc.FileInput(t, 2024, 4),
			Want:  "2468",
		},
	}

	aoc.AOCTest(t, day04p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`),
			Want: "9",
		},
		{
			Input: aoc.FileInput(t, 2024, 4),
			Want:  "1864",
		},
	}

	aoc.AOCTest(t, day04p02, tests)
}

func TestWordSearchCountsAtEdges(t *testing.T) {
	aoc.AOCTest(t, day04p01, []aoc.TestInput{{Input: strings.NewReader("XMAS\nMMMM\nAAAA\nSAMX"), Want: "3"}})
	aoc.AOCTest(t, day04p02, []aoc.TestInput{{Input: strings.NewReader("MXS\nXAX\nMXS"), Want: "1"}})
}

func TestWordSearchInputErrors(t *testing.T) {
	for _, input := range []string{"XMAS\nXMA", "XMAS\n\nXMAS", "XMAX\nXMAx", "XM.S"} {
		for _, solve := range []func(io.Reader) (string, error){day04p01, day04p02} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("input %q: got %q; want error", input, got)
			}
		}
	}
}

func TestDay04ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day04p01, "XMAS\n"},
		{"part2", day04p02, "XMAS\n"},
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
