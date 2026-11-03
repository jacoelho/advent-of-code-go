package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day04p01(t *testing.T) {
	aoc.AOCTest(t, day04p01, []aoc.TestInput{
		{Input: strings.NewReader("abcdef"), Want: "609043"},
		{Input: strings.NewReader("pqrstuv"), Want: "1048970"},
		{Input: aoc.FileInput(t, 2015, 4), Want: "254575"},
	})
}

func Test_day04p02(t *testing.T) {
	aoc.AOCTest(t, day04p02, []aoc.TestInput{
		{Input: strings.NewReader("abcdef"), Want: "6742839"},
		{Input: aoc.FileInput(t, 2015, 4), Want: "1038736"},
	})
}

func TestDay04ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day04p01", day04p01},
		{"day04p02", day04p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay04MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"empty secret", "", day04p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
