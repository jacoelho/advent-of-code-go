package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day01Example = `199
200
208
210
200
207
240
269
260
263`

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day01Example), Want: "7"},
		{Input: strings.NewReader(``), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 1), Want: "1527"},
	}

	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day01Example), Want: "5"},
		{Input: strings.NewReader(`1
2
3`), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 1), Want: "1575"},
	}

	aoc.AOCTest(t, day01p02, tests)
}

func TestDay01Malformed(t *testing.T) {
	if _, err := day01p01(strings.NewReader(`not a depth`)); err == nil {
		t.Fatal("expected malformed error")
	}
}

func TestDay01ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day01p01},
		{"part2", day01p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
