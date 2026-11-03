package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`2333133121414131402`),
			Want:  "1928",
		},
		{
			Input: aoc.FileInput(t, 2024, 9),
			Want:  "6332189866718",
		},
	}

	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`2333133121414131402`),
			Want:  "2858",
		},
		{
			Input: aoc.FileInput(t, 2024, 9),
			Want:  "6353648390778",
		},
	}

	aoc.AOCTest(t, day09p02, tests)
}

func TestDiskZeroLengthFiles(t *testing.T) {
	tests := []struct {
		input string
		want  [2]string
	}{
		{"12023", [2]string{"12", "12"}},
		{"120203", [2]string{"0", "0"}},
		{"1202023", [2]string{"18", "18"}},
		{"1200203", [2]string{"38", "60"}},
		{"00000", [2]string{"0", "0"}},
	}
	for _, tt := range tests {
		for part, solve := range []func(io.Reader) (string, error){day09p01, day09p02} {
			got, err := solve(strings.NewReader(tt.input))
			if err != nil || got != tt.want[part] {
				t.Errorf("part %d input %q: got %q, %v; want %q", part+1, tt.input, got, err, tt.want[part])
			}
		}
	}
}

func TestDiskMapInput(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day09p01, day09p02} {
		for _, input := range []string{"12023\n", "12023\r\n"} {
			if got, err := solve(strings.NewReader(input)); err != nil || got != "12" {
				t.Errorf("input %q: got %q, %v; want 12", input, got, err)
			}
		}
		for _, input := range []string{"12x3", "12\n3"} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("input %q: got %q, want error", input, got)
			}
		}
	}
}

func TestDay09ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day09p01, "12345\n"},
		{"part2", day09p02, "12345\n"},
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
