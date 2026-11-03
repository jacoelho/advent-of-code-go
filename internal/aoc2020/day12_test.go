package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`F10
N3
F7
R90
F11`), Want: "25"},
		{Input: strings.NewReader(`F2
R180
F1`), Want: "1"},
		{Input: strings.NewReader(`R270
F2
E3`), Want: "5"},
		{Input: strings.NewReader(`L270
F2
N3`), Want: "1"},
		{Input: aoc.FileInput(t, 2020, 12), Want: "1710"},
	}
	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`F10
N3
F7
R90
F11`), Want: "286"},
		{Input: strings.NewReader(`F2
R180
F1`), Want: "11"},
		{Input: strings.NewReader(`R270
F2
E3`), Want: "22"},
		{Input: strings.NewReader(`L270
F2
N3`), Want: "22"},
		{Input: aoc.FileInput(t, 2020, 12), Want: "62045"},
	}
	aoc.AOCTest(t, day12p02, tests)
}

func TestDay12ReadErrors(t *testing.T) {
	failure := errors.New("read failure")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"12/1", day12p01},
		{"12/2", day12p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestDay12MalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"navigation", "R45", day12p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected malformed-input error")
			}
		})
	}
}
