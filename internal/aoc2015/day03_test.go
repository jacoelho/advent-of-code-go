package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	aoc.AOCTest(t, day03p01, []aoc.TestInput{
		{Input: strings.NewReader(">"), Want: "2"},
		{Input: strings.NewReader("^>v<"), Want: "4"},
		{Input: strings.NewReader("^v^v^v^v^v"), Want: "2"},
		{Input: aoc.FileInput(t, 2015, 3), Want: "2572"},
	})
}

func Test_day03p02(t *testing.T) {
	aoc.AOCTest(t, day03p02, []aoc.TestInput{
		{Input: strings.NewReader("^v"), Want: "3"},
		{Input: strings.NewReader("^>v<"), Want: "3"},
		{Input: strings.NewReader("^v^v^v^v^v"), Want: "11"},
		{Input: strings.NewReader("^^^"), Want: "3"},
		{Input: aoc.FileInput(t, 2015, 3), Want: "2631"},
	})
}
