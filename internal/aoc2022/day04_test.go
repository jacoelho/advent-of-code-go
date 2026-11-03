package aoc2022

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day04Example = `2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8`

func Test_day04p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day04Example), Want: "2"},
		{Input: aoc.FileInput(t, 2022, 4), Want: "569"},
	}
	aoc.AOCTest(t, day04p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day04Example), Want: "4"},
		{Input: aoc.FileInput(t, 2022, 4), Want: "936"},
	}
	aoc.AOCTest(t, day04p02, tests)
}

func TestAssignmentEndpoints(t *testing.T) {
	aoc.AOCTest(t, day04p02, []aoc.TestInput{{Input: strings.NewReader("1-2,2-3"), Want: "1"}})
	if _, err := day04p01(strings.NewReader("3-1,2-4")); err == nil {
		t.Fatal("expected reversed range error")
	}
}
