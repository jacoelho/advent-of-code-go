package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const licenseExample = "2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2"

func Test_day08p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(licenseExample), Want: "138"},
		{Input: aoc.FileInput(t, 2018, 8), Want: "38567"},
	}

	aoc.AOCTest(t, day08p01, tests)
}

func Test_day08p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(licenseExample), Want: "66"},
		{Input: strings.NewReader("1 4 0 1 5 0 2 1 1"), Want: "10"},
		{Input: aoc.FileInput(t, 2018, 8), Want: "24453"},
	}

	aoc.AOCTest(t, day08p02, tests)
}

func TestDay08ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"08/1", day08p01},
		{"08/2", day08p02},
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

func TestDay08InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
		input string
	}{
		{name: "incomplete license", solve: day08p01, input: "0 2 1"},
		{name: "trailing license data", solve: day08p02, input: "0 0 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an input error")
			}
		})
	}
}
