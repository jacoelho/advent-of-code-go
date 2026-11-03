package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day14Example = `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`

func Test_day14p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day14Example), Want: "24"},
		{Input: aoc.FileInput(t, 2022, 14), Want: "1513"},
	}
	aoc.AOCTest(t, day14p01, tests)
}

func Test_day14p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day14Example), Want: "93"},
		{Input: aoc.FileInput(t, 2022, 14), Want: "22646"},
	}
	aoc.AOCTest(t, day14p02, tests)
}

func TestInfiniteSandFloor(t *testing.T) {
	aoc.AOCTest(t, day14p01, []aoc.TestInput{{Input: strings.NewReader("0,1 -> 0,1"), Want: "0"}})
	aoc.AOCTest(t, day14p02, []aoc.TestInput{{Input: strings.NewReader("0,1 -> 0,1"), Want: "9"}})
	if _, err := day14p02(strings.NewReader("0,1 -> 1,2")); err == nil {
		t.Fatal("expected diagonal rock error")
	}
}

func TestDay14ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day14p01},
		{"part2", day14p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
