package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day07Example = `16,1,2,0,4,2,7,1,2,14`

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day07Example), Want: "37"},
		{Input: aoc.FileInput(t, 2021, 7), Want: "344297"},
	}

	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day07Example), Want: "168"},
		{Input: strings.NewReader(`0,2`), Want: "2"},
		{Input: aoc.FileInput(t, 2021, 7), Want: "97164301"},
	}

	aoc.AOCTest(t, day07p02, tests)
}

func TestDay07NegativePosition(t *testing.T) {
	if _, err := day07p02(strings.NewReader(`-1,2`)); err == nil {
		t.Fatal("expected negative position error")
	}
}

func TestDay07ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day07p01},
		{"part2", day07p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
