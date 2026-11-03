package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day11Example = `5483143223
2745854711
5264556173
6141336146
6357385478
4167524645
2176841721
6882881134
4846848554
5283751526`

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day11Example), Want: "1656"},
		{Input: strings.NewReader(`0`), Want: "10"},
		{Input: aoc.FileInput(t, 2021, 11), Want: "1739"},
	}

	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day11Example), Want: "195"},
		{Input: strings.NewReader(`0`), Want: "10"},
		{Input: aoc.FileInput(t, 2021, 11), Want: "324"},
	}

	aoc.AOCTest(t, day11p02, tests)
}

func TestDay11NoSynchronization(t *testing.T) {
	if _, err := day11p02(strings.NewReader(`05`)); err == nil {
		t.Fatal("expected no synchronization error")
	}
}

func TestDay11ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day11p01},
		{"part2", day11p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
