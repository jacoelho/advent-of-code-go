package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test2016Days01To08ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day01p01", day01p01}, {"day01p02", day01p02},
		{"day02p01", day02p01}, {"day02p02", day02p02},
		{"day03p01", day03p01}, {"day03p02", day03p02},
		{"day04p01", day04p01}, {"day04p02", day04p02},
		{"day05p01", day05p01}, {"day05p02", day05p02},
		{"day06p01", day06p01}, {"day06p02", day06p02},
		{"day07p01", day07p01}, {"day07p02", day07p02},
		{"day08p01", day08p01}, {"day08p02", day08p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func Test2016Days01To08MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"turn", "Q3", day01p01},
		{"no repeat", "R2, L3", day01p02},
		{"keypad direction", "UX", day02p01},
		{"triangle sides", "1 2", day03p01},
		{"triangle row group", "3 4 5", day03p02},
		{"room", "room-xyz[abcde]", day04p01},
		{"door ID", "", day05p01},
		{"message width", "abc\nab", day06p01},
		{"IPv7 brackets", "aba[bab", day07p02},
		{"screen bounds", "rect 51x1", day08p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

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
