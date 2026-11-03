package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`v...>>.vv>
.vv>>.vv..
>>.>v>...v
>>v>>.>.v.
v>v.vv.v..
>.>>..v...
.vv..>.>v.
v.v..>>v.v
....v..v.>`), Want: "58"},
		{Input: strings.NewReader("."), Want: "1"},
		{Input: strings.NewReader(">"), Want: "1"},
		{Input: strings.NewReader("v"), Want: "1"},
		{Input: aoc.FileInput(t, 2021, 25), Want: "305"},
	}
	aoc.AOCTest(t, day25p01, tests)
}

func TestSeaCucumberPhases(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{input: ">>.\n...", want: ">.>..."},
		{input: ".>\nv.", want: ">.v."},
		{input: ".v\n.>", want: "..>v"},
		{input: ">v\n..", want: ">..v"},
	}
	for _, tt := range tests {
		floor, err := parseSeaFloor(strings.NewReader(tt.input))
		if err != nil {
			t.Fatal(err)
		}
		if !floor.step(make([]byte, len(floor.cells))) || string(floor.cells) != tt.want {
			t.Fatalf("one step from %q: got %q; want %q", tt.input, floor.cells, tt.want)
		}
	}
}

func TestSeaCucumberFailures(t *testing.T) {
	for _, input := range []string{"", ">.", "v\n.", ".x", "..\n.", "..\n\n.."} {
		if _, err := day25p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected sea cucumber error for %q", input)
		}
	}
}

func TestDay25ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day25p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
