package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day08Example = `30373
25512
65332
33549
35390`

func Test_day08p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day08Example), Want: "21"},
		{Input: aoc.FileInput(t, 2022, 8), Want: "1801"},
	}
	aoc.AOCTest(t, day08p01, tests)
}

func Test_day08p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day08Example), Want: "8"},
		{Input: aoc.FileInput(t, 2022, 8), Want: "209880"},
	}
	aoc.AOCTest(t, day08p02, tests)
}

func TestForestShape(t *testing.T) {
	if _, err := day08p01(strings.NewReader("12\n3")); err == nil {
		t.Fatal("expected ragged forest error")
	}
}

func TestDay08ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day08p01},
		{"part2", day08p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
