package aoc2019

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day24p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`....#
#..#.
#..##
..#..
#....`),
			Want: "2129920",
		},
		{
			Input: aoc.FileInput(t, 2019, 24),
			Want:  "18401265",
		},
	}
	aoc.AOCTest(t, day24p01, tests)
}

func Test_day24p02_example(t *testing.T) {
	result, err := simulateRecursive(strings.NewReader(`....#
#..#.
#.?##
..#..
#....`), 10)
	if err != nil {
		t.Fatal(err)
	}
	if result != "99" {
		t.Errorf("got %s, want 99", result)
	}
}

func Test_day24p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 24),
			Want:  "2078",
		},
	}
	aoc.AOCTest(t, day24p02, tests)
}
