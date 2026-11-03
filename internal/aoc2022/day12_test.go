package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day12Example = `Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi`

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day12Example), Want: "31"},
		{Input: aoc.FileInput(t, 2022, 12), Want: "447"},
	}
	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day12Example), Want: "29"},
		{Input: aoc.FileInput(t, 2022, 12), Want: "446"},
	}
	aoc.AOCTest(t, day12p02, tests)
}

func TestUnreachableHill(t *testing.T) {
	if _, err := day12p01(strings.NewReader("SE")); err == nil {
		t.Fatal("expected unreachable hill error")
	}
}

func TestDay12ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day12p01},
		{"part2", day12p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
