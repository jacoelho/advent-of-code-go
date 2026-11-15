package aoc2022

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
		{Input: strings.NewReader("1=-0-2\n12111\n2=0=\n21\n2=01\n111\n20012\n112\n1=-1=\n1-12\n12\n1=\n122"), Want: "2=-1=0"},
		{Input: strings.NewReader("0"), Want: "0"},
		{Input: strings.NewReader("="), Want: "="},
		{Input: strings.NewReader("=\n-"), Want: "-2"},
		{Input: strings.NewReader("2\n="), Want: "0"},
		{Input: strings.NewReader("1" + strings.Repeat("0", 100) + "\n="), Want: "1" + strings.Repeat("0", 99) + "="},
		{Input: aoc.FileInput(t, 2022, 25), Want: "2011-=2=-1020-1===-1"},
	}
	aoc.AOCTest(t, day25p01, tests)
}

func TestSNAFUSumCarries(t *testing.T) {
	tests := []struct {
		digit string
		count int
		want  string
	}{
		{digit: "-", count: 3, want: "-2"},
		{digit: "-", count: 2, want: "="},
		{digit: "1", count: 3, want: "1="},
		{digit: "1", count: 4, want: "1-"},
		{digit: "1", count: 5, want: "10"},
		{digit: "1", count: 2022, want: "1=11-2"},
		{digit: "1", count: 12345, want: "1-0---0"},
	}
	for _, tt := range tests {
		input := strings.TrimSuffix(strings.Repeat(tt.digit+"\n", tt.count), "\n")
		got, err := day25p01(strings.NewReader(input))
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Fatalf("%d copies of %q: got %q; want %q", tt.count, tt.digit, got, tt.want)
		}
	}
}

func TestSNAFUInvalidInput(t *testing.T) {
	for _, input := range []string{"", "3", "1\n\n2"} {
		if _, err := day25p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected SNAFU error for %q", input)
		}
	}
}

func TestDay25ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day25p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
