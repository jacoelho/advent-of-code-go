package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day08p01(t *testing.T) {
	aoc.AOCTest(t, day08p01, []aoc.TestInput{
		{Input: strings.NewReader("\"\"\n\"abc\"\n\"aaa\\\"aaa\"\n\"\\x27\""), Want: "12"},
		{Input: aoc.FileInput(t, 2015, 8), Want: "1333"},
	})
}

func Test_day08p02(t *testing.T) {
	aoc.AOCTest(t, day08p02, []aoc.TestInput{
		{Input: strings.NewReader("\"\"\n\"abc\"\n\"aaa\\\"aaa\"\n\"\\x27\""), Want: "19"},
		{Input: aoc.FileInput(t, 2015, 8), Want: "2046"},
	})
}

func TestDay08ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day08p01", day08p01},
		{"day08p02", day08p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay08MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"string escape", `"\xgg"`, day08p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
