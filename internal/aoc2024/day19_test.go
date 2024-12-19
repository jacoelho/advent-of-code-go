package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day19p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb`),
			Want: "6",
		},
		{
			Input: aoc.FileInput(t, 2024, 19),
			Want:  "308",
		},
	}
	aoc.AOCTest(t, day19p01, tests)
}

func Test_day19p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb`),
			Want: "16",
		},
		{
			Input: aoc.FileInput(t, 2024, 19),
			Want:  "662726441391898",
		},
	}
	aoc.AOCTest(t, day19p02, tests)
}
