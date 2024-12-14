package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day14p01example(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3`),
			Want: "12",
		},
	}
	aoc.AOCTest(t, day14p01(11, 7), tests)
}

func Test_day14p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 14),
			Want:  "229868730",
		},
	}
	aoc.AOCTest(t, day14p01(101, 103), tests)
}

func Test_day14p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 14),
			Want:  "7861",
		},
	}
	aoc.AOCTest(t, day14p02, tests)
}
