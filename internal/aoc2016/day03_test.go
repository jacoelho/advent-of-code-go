package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	aoc.AOCTest(t, day03p01, []aoc.TestInput{
		{Input: strings.NewReader("5 10 25"), Want: "0"},
		{Input: strings.NewReader("1 1 2\n3 4 5"), Want: "1"},
		{Input: aoc.FileInput(t, 2016, 3), Want: "1050"},
	})
}

func Test_day03p02(t *testing.T) {
	aoc.AOCTest(t, day03p02, []aoc.TestInput{
		{Input: strings.NewReader("101 301 501\n102 302 502\n103 303 503\n201 401 601\n202 402 602\n203 403 603"), Want: "6"},
		{Input: aoc.FileInput(t, 2016, 3), Want: "1921"},
	})
}

func TestDay03ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day03p01", "3 4 5\n3 4 5\n3 4 5\n", day03p01},
		{"day03p02", "3 4 5\n3 4 5\n3 4 5\n", day03p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				reader := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(reader); !errors.Is(err, failure) {
					t.Fatalf("prefix %q: got %v, want input failure", prefix, err)
				}
			}
		})
	}
}

func TestDay03ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day03p01", day03p01},
		{"day03p02", day03p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay03MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"triangle sides", "1 2", day03p01},
		{"triangle row group", "3 4 5", day03p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
