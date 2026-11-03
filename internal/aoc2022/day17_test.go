package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(rockExample), Want: "3068"},
		{Input: aoc.FileInput(t, 2022, 17), Want: "3059"},
	}

	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(rockExample), Want: "1514285714288"},
		{Input: aoc.FileInput(t, 2022, 17), Want: "1500874635587"},
	}

	aoc.AOCTest(t, day17p02, tests)
}

func TestRockHeights(t *testing.T) {
	for _, test := range []struct {
		jets  string
		count int64
		want  string
	}{
		{rockExample, 10, "17"},
		{"<", 100, "220"},
		{">", 200, "520"},
		{"><><<<>>", 200, "283"},
		{rockExample, 63, "102"},
		{rockExample, 64, "104"},
		{rockExample, 98, "155"},
		{rockExample, 99, "157"},
	} {
		got, err := fallingRockHeight(strings.NewReader(test.jets), test.count)
		if err != nil {
			t.Fatal(err)
		}
		if got != test.want {
			t.Errorf("%d rocks with %q = %s, want %s", test.count, test.jets, got, test.want)
		}
	}
}

func TestRockErrors(t *testing.T) {
	for _, input := range []string{"", "<x>", "<>\n<>"} {
		if _, err := day17p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected invalid jet error for %q", input)
		}
	}
}

const rockExample = `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`

func TestDay17ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day17p01},
		{"part2", day17p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
