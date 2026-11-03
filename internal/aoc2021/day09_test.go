package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day09Example = `2199943210
3987894921
9856789892
8767896789
9899965678`

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day09Example), Want: "15"},
		{Input: strings.NewReader(`5`), Want: "6"},
		{Input: aoc.FileInput(t, 2021, 9), Want: "550"},
	}

	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day09Example), Want: "1134"},
		{Input: aoc.FileInput(t, 2021, 9), Want: "1100682"},
	}

	aoc.AOCTest(t, day09p02, tests)
}

func TestDay09FewerThanThreeBasins(t *testing.T) {
	if _, err := day09p02(strings.NewReader(`5`)); err == nil {
		t.Fatal("expected fewer than three basins error")
	}
}

func TestDay09RaggedHeightmap(t *testing.T) {
	if _, err := day09p01(strings.NewReader(`1
12`)); err == nil {
		t.Fatal("expected ragged heightmap error")
	}
}

func TestDay09ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day09p01},
		{"part2", day09p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
