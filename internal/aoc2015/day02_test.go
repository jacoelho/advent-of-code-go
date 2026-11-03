package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	aoc.AOCTest(t, day02p01, []aoc.TestInput{
		{Input: strings.NewReader("2x3x4"), Want: "58"},
		{Input: strings.NewReader("1x1x10"), Want: "43"},
		{Input: aoc.FileInput(t, 2015, 2), Want: "1588178"},
	})
}

func Test_day02p02(t *testing.T) {
	aoc.AOCTest(t, day02p02, []aoc.TestInput{
		{Input: strings.NewReader("2x3x4"), Want: "34"},
		{Input: strings.NewReader("1x1x10"), Want: "14"},
		{Input: aoc.FileInput(t, 2015, 2), Want: "3783758"},
	})
}
