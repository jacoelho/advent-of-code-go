package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	aoc.AOCTest(t, day05p01, []aoc.TestInput{
		{Input: strings.NewReader("abc"), Want: "18f47a30"},
		{Input: aoc.FileInput(t, 2016, 5), Want: "f77a0e6e"},
	})
}

func Test_day05p02(t *testing.T) {
	aoc.AOCTest(t, day05p02, []aoc.TestInput{
		{Input: strings.NewReader("abc"), Want: "05ace8e3"},
		{Input: aoc.FileInput(t, 2016, 5), Want: "999828ec"},
	})
}

func TestDay05ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day05p01", day05p01},
		{"day05p02", day05p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay05MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"door ID", "", day05p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
