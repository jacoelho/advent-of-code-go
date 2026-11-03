package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day15Example = `1163751742
1381373672
2136511328
3694931569
7463417111
1319128137
1359912421
3125421639
1293138521
2311944581`

func Test_day15p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day15Example), Want: "40"},
		{Input: strings.NewReader("12\n34\n56"), Want: "12"},
		{Input: strings.NewReader(`7`), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 15), Want: "386"},
	}

	aoc.AOCTest(t, day15p01, tests)
}

func Test_day15p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day15Example), Want: "315"},
		{Input: strings.NewReader("8"), Want: "37"},
		{Input: aoc.FileInput(t, 2021, 15), Want: "2806"},
	}

	aoc.AOCTest(t, day15p02, tests)
}

func TestDay15ZeroRisk(t *testing.T) {
	if _, err := day15p01(strings.NewReader(`0`)); err == nil {
		t.Fatal("expected zero risk error")
	}
}

func TestDay15ExpandedGridLimit(t *testing.T) {
	if _, err := day15p02(strings.NewReader(strings.Repeat("1", 40001))); err == nil {
		t.Fatal("expected expanded grid limit error")
	}
}

func TestDay15ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day15p01},
		{"part2", day15p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
