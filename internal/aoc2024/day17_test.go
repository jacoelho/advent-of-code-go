package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0`),
			Want: "4,6,3,5,6,3,5,2,1,0",
		},
		{
			Input: aoc.FileInput(t, 2024, 17),
			Want:  "2,7,4,7,2,1,7,5,1",
		},
	}
	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Register A: 2024
		Register B: 0
		Register C: 0
		
		Program: 0,3,5,4,3,0`),
			Want: "117440",
		},
		{
			Input: aoc.FileInput(t, 2024, 17),
			Want:  "37221274271220",
		},
	}
	aoc.AOCTest(t, day17p02, tests)
}
