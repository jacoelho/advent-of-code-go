package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestRouteMustVisitEveryLocation(t *testing.T) {
	input := "A to B = 1\nC to D = 2"
	if _, err := day09p01(strings.NewReader(input)); err == nil {
		t.Fatal("disconnected locations cannot form a complete route")
	}
}

func Test_day09p01(t *testing.T) {
	aoc.AOCTest(t, day09p01, []aoc.TestInput{
		{Input: strings.NewReader("London to Dublin = 464\nLondon to Belfast = 518\nDublin to Belfast = 141"), Want: "605"},
		{Input: aoc.FileInput(t, 2015, 9), Want: "207"},
	})
}

func Test_day09p02(t *testing.T) {
	aoc.AOCTest(t, day09p02, []aoc.TestInput{
		{Input: strings.NewReader("London to Dublin = 464\nLondon to Belfast = 518\nDublin to Belfast = 141"), Want: "982"},
		{Input: aoc.FileInput(t, 2015, 9), Want: "804"},
	})
}
