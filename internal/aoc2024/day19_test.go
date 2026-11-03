package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day19p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb`),
			Want: "6",
		},
		{
			Input: aoc.FileInput(t, 2024, 19),
			Want:  "308",
		},
	}
	aoc.AOCTest(t, day19p01, tests)
}

func Test_day19p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb`),
			Want: "16",
		},
		{
			Input: aoc.FileInput(t, 2024, 19),
			Want:  "662726441391898",
		},
	}
	aoc.AOCTest(t, day19p02, tests)
}

func TestDesignParseErrors(t *testing.T) {
	for _, input := range []string{"r, b", "r, \n\nr"} {
		for _, solve := range []func(io.Reader) (string, error){day19p01, day19p02} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: got %q; want towel parse error", input, got)
			}
		}
	}
}

func TestDesignLineEndings(t *testing.T) {
	for _, newline := range []string{"\n", "\r\n"} {
		input := strings.ReplaceAll("r, b\n\nrb\n", "\n", newline)
		for _, solve := range []func(io.Reader) (string, error){day19p01, day19p02} {
			got, err := solve(strings.NewReader(input))
			if err != nil || got != "1" {
				t.Fatalf("got %q, %v; want 1", got, err)
			}
		}
	}
}

func TestDay19ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day19p01, "r, b\n\nr\n"},
		{"part2", day19p02, "r, b\n\nr\n"},
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
