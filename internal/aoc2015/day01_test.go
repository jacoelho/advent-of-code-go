package aoc2015

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
		{Input: strings.NewReader("()()"), Want: "0"},
		{Input: strings.NewReader("))((((("), Want: "3"},
		{Input: aoc.FileInput(t, 2015, 1), Want: "280"},
	})
}

func Test_day01p02(t *testing.T) {
	aoc.AOCTest(t, day01p02, []aoc.TestInput{
		{Input: strings.NewReader("()())"), Want: "5"},
		{Input: strings.NewReader(")"), Want: "1"},
		{Input: aoc.FileInput(t, 2015, 1), Want: "1797"},
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
		{"floor direction", "(x)", day01p01},
		{"no basement", "()", day01p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
