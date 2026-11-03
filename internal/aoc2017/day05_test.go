package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`0
3
0
1
-3`), Want: "5"},
		{Input: strings.NewReader("-1"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 5), Want: "394829"},
	}

	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`0
3
0
1
-3`), Want: "10"},
		{Input: strings.NewReader("-1"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 5), Want: "31150702"},
	}

	aoc.AOCTest(t, day05p02, tests)
}

func TestDay05ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"05/1", day05p01},
		{"05/2", day05p02},
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

func TestDay05MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"jump offset", "not a number", day05p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
