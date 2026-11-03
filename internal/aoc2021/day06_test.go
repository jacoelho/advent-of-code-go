package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day06Example = `3,4,3,1,2`

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day06Example), Want: "5934"},
		{Input: aoc.FileInput(t, 2021, 6), Want: "361169"},
	}

	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day06Example), Want: "26984457539"},
		{Input: aoc.FileInput(t, 2021, 6), Want: "1634946868992"},
	}

	aoc.AOCTest(t, day06p02, tests)
}

func TestDay06InvalidTimer(t *testing.T) {
	if _, err := day06p01(strings.NewReader(`9`)); err == nil {
		t.Fatal("expected invalid timer error")
	}
}

func TestDay06ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day06p01},
		{"part2", day06p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
