package aoc2021

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day07Example = `16,1,2,0,4,2,7,1,2,14`

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day07Example), Want: "37"},
		{Input: aoc.FileInput(t, 2021, 7), Want: "344297"},
	}

	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day07Example), Want: "168"},
		{Input: strings.NewReader(`0,2`), Want: "2"},
		{Input: aoc.FileInput(t, 2021, 7), Want: "97164301"},
	}

	aoc.AOCTest(t, day07p02, tests)
}

func TestDay07NegativePosition(t *testing.T) {
	if _, err := day07p02(strings.NewReader(`-1,2`)); err == nil {
		t.Fatal("expected negative position error")
	}
}
