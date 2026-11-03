package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`F10
N3
F7
R90
F11`), Want: "25"},
		{Input: strings.NewReader(`F2
R180
F1`), Want: "1"},
		{Input: strings.NewReader(`R270
F2
E3`), Want: "5"},
		{Input: strings.NewReader(`L270
F2
N3`), Want: "1"},
		{Input: aoc.FileInput(t, 2020, 12), Want: "1710"},
	}
	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`F10
N3
F7
R90
F11`), Want: "286"},
		{Input: strings.NewReader(`F2
R180
F1`), Want: "11"},
		{Input: strings.NewReader(`R270
F2
E3`), Want: "22"},
		{Input: strings.NewReader(`L270
F2
N3`), Want: "22"},
		{Input: aoc.FileInput(t, 2020, 12), Want: "62045"},
	}
	aoc.AOCTest(t, day12p02, tests)
}
