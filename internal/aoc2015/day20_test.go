package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	aoc.AOCTest(t, day20p01, []aoc.TestInput{
		{Input: strings.NewReader("11"), Want: "2"},
		{Input: strings.NewReader("150"), Want: "8"},
		{Input: aoc.FileInput(t, 2015, 20), Want: "786240"},
	})
}

func Test_day20p02(t *testing.T) {
	aoc.AOCTest(t, day20p02, []aoc.TestInput{
		{Input: strings.NewReader("12"), Want: "2"},
		{Input: aoc.FileInput(t, 2015, 20), Want: "831600"},
	})
}

func TestElfVisitLimit(t *testing.T) {
	got, err := firstHouse(strings.NewReader("70"), 1, 2)
	if err != nil || got != "48" {
		t.Fatalf("house = %q, %v; want48", got, err)
	}
}
