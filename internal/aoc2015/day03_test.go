package aoc2015

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
		{Input: strings.NewReader(">"), Want: "2"},
		{Input: strings.NewReader("^>v<"), Want: "4"},
		{Input: strings.NewReader("^v^v^v^v^v"), Want: "2"},
		{Input: aoc.FileInput(t, 2015, 3), Want: "2572"},
	})
}

func Test_day03p02(t *testing.T) {
	aoc.AOCTest(t, day03p02, []aoc.TestInput{
		{Input: strings.NewReader("^v"), Want: "3"},
		{Input: strings.NewReader("^>v<"), Want: "3"},
		{Input: strings.NewReader("^v^v^v^v^v"), Want: "11"},
		{Input: strings.NewReader("^^^"), Want: "3"},
		{Input: aoc.FileInput(t, 2015, 3), Want: "2631"},
	})
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
		{"delivery direction", "north", day03p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
