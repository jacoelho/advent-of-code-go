package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("18"),
			Want:  "33,45",
		},
		{
			Input: strings.NewReader("42"),
			Want:  "21,61",
		},
		{
			Input: aoc.FileInput(t, 2018, 11),
			Want:  "20,46",
		},
	}

	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("18"),
			Want:  "90,269,16",
		},
		{
			Input: strings.NewReader("42"),
			Want:  "232,251,12",
		},
		{
			Input: aoc.FileInput(t, 2018, 11),
			Want:  "231,65,14",
		},
	}

	aoc.AOCTest(t, day11p02, tests)
}

func TestFuelGridIncludesBoundarySquares(t *testing.T) {
	var totals [301][301]int
	totals[300][300] = 100
	if x, y, size := bestFuelSquare(&totals, 1, 1); x != 300 || y != 300 || size != 1 {
		t.Errorf("got %d,%d,%d, want 300,300,1", x, y, size)
	}
	if x, y, size := bestFuelSquare(&totals, 300, 300); x != 1 || y != 1 || size != 300 {
		t.Errorf("got %d,%d,%d, want 1,1,300", x, y, size)
	}
}

func TestFuelGridErrors(t *testing.T) {
	for _, input := range []string{"", "-1", "serial"} {
		if _, err := day11p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}
