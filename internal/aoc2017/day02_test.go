package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`5 1 9 5
7 5 3
2 4 6 8`), Want: "18"},
		{Input: aoc.FileInput(t, 2017, 2), Want: "30994"},
	}

	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`5 9 2 8
9 4 7 3
3 8 6 5`), Want: "9"},
		{Input: strings.NewReader("2 2"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 2), Want: "233"},
	}

	aoc.AOCTest(t, day02p02, tests)
}

func TestDay02ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"02/1", day02p01},
		{"02/2", day02p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}

func TestDay02MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"spreadsheet zero", "5 0", day02p02},
		{"missing pair", "5 7", day02p02},
		{"ambiguous pairs", "2 4 8", day02p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
