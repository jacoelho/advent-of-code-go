package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestRouteMustVisitEveryLocation(t *testing.T) {
	input := "A to B = 1\nC to D = 2"
	if _, err := day09p01(strings.NewReader(input)); err == nil {
		t.Fatal("disconnected locations cannot form a complete route")
	}
}

func Test_day09p01(t *testing.T) {
	aoc.AOCTest(t, day09p01, []aoc.TestInput{
		{Input: strings.NewReader("London to Dublin = 464\nLondon to Belfast = 518\nDublin to Belfast = 141"), Want: "605"},
		{Input: aoc.FileInput(t, 2015, 9), Want: "207"},
	})
}

func Test_day09p02(t *testing.T) {
	aoc.AOCTest(t, day09p02, []aoc.TestInput{
		{Input: strings.NewReader("London to Dublin = 464\nLondon to Belfast = 518\nDublin to Belfast = 141"), Want: "982"},
		{Input: aoc.FileInput(t, 2015, 9), Want: "804"},
	})
}

func TestDay09ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day09p01", day09p01},
		{"day09p02", day09p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay09MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"route distance", "A to B = nope", day09p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
