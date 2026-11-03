package aoc2016

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
		{Input: strings.NewReader("ULL\nRRDDD\nLURDL\nUUUUD"), Want: "1985"},
		{Input: aoc.FileInput(t, 2016, 2), Want: "74921"},
	})
}

func Test_day02p02(t *testing.T) {
	aoc.AOCTest(t, day02p02, []aoc.TestInput{
		{Input: strings.NewReader("ULL\nRRDDD\nLURDL\nUUUUD"), Want: "5DB3"},
		{Input: aoc.FileInput(t, 2016, 2), Want: "A6B35"},
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
		{"keypad direction", "UX", day02p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
