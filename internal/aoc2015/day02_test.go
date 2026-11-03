package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	aoc.AOCTest(t, day02p01, []aoc.TestInput{
		{Input: strings.NewReader("2x3x4"), Want: "58"},
		{Input: strings.NewReader("1x1x10"), Want: "43"},
		{Input: aoc.FileInput(t, 2015, 2), Want: "1588178"},
	})
}

func Test_day02p02(t *testing.T) {
	aoc.AOCTest(t, day02p02, []aoc.TestInput{
		{Input: strings.NewReader("2x3x4"), Want: "34"},
		{Input: strings.NewReader("1x1x10"), Want: "14"},
		{Input: aoc.FileInput(t, 2015, 2), Want: "3783758"},
	})
}

func TestDay02ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day02p01", day02p01},
		{"day02p02", day02p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay02MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"box dimensions", "2x3", day02p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
