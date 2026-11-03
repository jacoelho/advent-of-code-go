package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day14p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("flqrgnkx"),
			Want:  "8108",
		},
		{
			Input: aoc.FileInput(t, 2017, 14),
			Want:  "8190",
		},
	}

	aoc.AOCTest(t, day14p01, tests)
}

func Test_day14p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("flqrgnkx"),
			Want:  "1242",
		},
		{
			Input: aoc.FileInput(t, 2017, 14),
			Want:  "1134",
		},
	}

	aoc.AOCTest(t, day14p02, tests)
}

func TestDiskRegionsUseOrthogonalNeighbours(t *testing.T) {
	var grid diskGrid
	grid[0][0], grid[1][1] = true, true
	if got := diskRegions(grid); got != 2 {
		t.Errorf("diagonal regions=%d, want 2", got)
	}

	grid[0][1] = true
	if got := diskRegions(grid); got != 1 {
		t.Errorf("connected regions=%d, want 1", got)
	}
}

func TestDiskErrors(t *testing.T) {
	if _, err := day14p01(strings.NewReader("")); err == nil {
		t.Error("expected empty-key error")
	}
}
