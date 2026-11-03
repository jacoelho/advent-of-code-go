package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	aoc.AOCTest(t, day01p01, []aoc.TestInput{
		{Input: strings.NewReader("R2, L3"), Want: "5"},
		{Input: strings.NewReader("R2, R2, R2"), Want: "2"},
		{Input: strings.NewReader("R5, L5, R5, R3"), Want: "12"},
		{Input: aoc.FileInput(t, 2016, 1), Want: "242"},
	})
}

func Test_day01p02(t *testing.T) {
	aoc.AOCTest(t, day01p02, []aoc.TestInput{
		{Input: strings.NewReader("R8, R4, R4, R8"), Want: "4"},
		{Input: strings.NewReader("R1, R1, R1, R1"), Want: "0"},
		{Input: aoc.FileInput(t, 2016, 1), Want: "150"},
	})
}

func TestDay01ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day01p01", day01p01},
		{"day01p02", day01p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay01MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"turn", "Q3", day01p01},
		{"no repeat", "R2, L3", day01p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
