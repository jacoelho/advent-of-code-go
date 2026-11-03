package aoc2021

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day01Example = `199
200
208
210
200
207
240
269
260
263`

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day01Example), Want: "7"},
		{Input: strings.NewReader(``), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 1), Want: "1527"},
	}

	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day01Example), Want: "5"},
		{Input: strings.NewReader(`1
2
3`), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 1), Want: "1575"},
	}

	aoc.AOCTest(t, day01p02, tests)
}

func TestDay01Malformed(t *testing.T) {
	if _, err := day01p01(strings.NewReader(`not a depth`)); err == nil {
		t.Fatal("expected malformed error")
	}
}
