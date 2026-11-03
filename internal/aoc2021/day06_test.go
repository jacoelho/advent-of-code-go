package aoc2021

import (
	"strings"
	"testing"

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
