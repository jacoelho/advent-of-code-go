package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const coordinatesExample = `1, 1
1, 6
8, 3
3, 4
5, 5
8, 9`

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(coordinatesExample), Want: "17"},
		{Input: strings.NewReader("1, 1\n0, 0\n0, 2\n2, 0\n2, 2"), Want: "1"},
		{Input: aoc.FileInput(t, 2018, 6), Want: "4754"},
	}

	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("0, 0"), Want: "199980001"},
		{Input: aoc.FileInput(t, 2018, 6), Want: "42344"},
	}

	aoc.AOCTest(t, day06p02, tests)
}

func TestSafeRegionSize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		limit int
		want  int
	}{
		{name: "example", input: coordinatesExample, limit: 32, want: 16},
		{name: "outside site bounds", input: "0, 0", limit: 3, want: 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := safeRegionSize(strings.NewReader(tt.input), tt.limit)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
