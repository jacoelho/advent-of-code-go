package aoc2022

import (
	"strings"
	"testing"

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
