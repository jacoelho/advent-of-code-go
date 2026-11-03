package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day03Example = `00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010`

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day03Example), Want: "198"},
		{Input: strings.NewReader(`00
00`), Want: "0"},
		{Input: strings.NewReader(`10
11`), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 3), Want: "3959450"},
	}

	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day03Example), Want: "230"},
		{Input: strings.NewReader(`00
00`), Want: "0"},
		{Input: strings.NewReader(`10
11`), Want: "6"},
		{Input: aoc.FileInput(t, 2021, 3), Want: "7440311"},
	}

	aoc.AOCTest(t, day03p02, tests)
}

func TestDay03RaggedReport(t *testing.T) {
	if _, err := day03p01(strings.NewReader(`0
01`)); err == nil {
		t.Fatal("expected ragged report error")
	}
}

func TestDay03InvalidBit(t *testing.T) {
	if _, err := day03p02(strings.NewReader(`+0`)); err == nil {
		t.Fatal("expected invalid bit error")
	}
}

func TestDay03ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day03p01},
		{"part2", day03p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
