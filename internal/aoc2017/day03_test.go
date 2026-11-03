package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("1"), Want: "0"},
		{Input: strings.NewReader("12"), Want: "3"},
		{Input: strings.NewReader("23"), Want: "2"},
		{Input: strings.NewReader("1024"), Want: "31"},
		{Input: aoc.FileInput(t, 2017, 3), Want: "438"},
	}

	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("1"), Want: "2"},
		{Input: strings.NewReader("747"), Want: "806"},
		{Input: aoc.FileInput(t, 2017, 3), Want: "266330"},
	}

	aoc.AOCTest(t, day03p02, tests)
}

func TestDay03ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"03/1", day03p01},
		{"03/2", day03p02},
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

func TestDay03MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"spiral target", "0", day03p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
