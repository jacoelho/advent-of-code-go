package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const fabricExample = `#1 @ 1,3: 4x4
#2 @ 3,1: 4x4
#3 @ 5,5: 2x2`

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(fabricExample), Want: "4"},
		{Input: strings.NewReader("#1 @ 0,0: 1x1\n#2 @ 1,0: 1x1"), Want: "0"},
		{Input: aoc.FileInput(t, 2018, 3), Want: "107820"},
	}

	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(fabricExample), Want: "3"},
		{Input: aoc.FileInput(t, 2018, 3), Want: "661"},
	}

	aoc.AOCTest(t, day03p02, tests)
}
