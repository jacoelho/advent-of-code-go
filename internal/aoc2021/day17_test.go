package aoc2021

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
		{Input: strings.NewReader("target area: x=20..30, y=-10..-5"), Want: "45"},
		{Input: strings.NewReader("target area: x=-30..-20, y=-10..-5"), Want: "45"},
		{Input: strings.NewReader("target area: x=-1..1, y=-1..-1"), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 17), Want: "10878"},
	}

	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("target area: x=20..30, y=-10..-5"), Want: "112"},
		{Input: strings.NewReader("target area: x=-30..-20, y=-10..-5"), Want: "112"},
		{Input: strings.NewReader("target area: x=-1..1, y=-1..-1"), Want: "6"},
		{Input: aoc.FileInput(t, 2021, 17), Want: "4716"},
	}

	aoc.AOCTest(t, day17p02, tests)
}

func TestProbeErrors(t *testing.T) {
	for _, input := range []string{
		"target area: x=1..2, y=0..1",
		"target area: x=2..1, y=-2..-1",
		"target area: x=1..2, y=-1..-2",
		"target area: x=-2147483648..2147483647, y=-2147483648..-1",
		"target area: x=1..2, y=invalid",
	} {
		if _, err := day17p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

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
