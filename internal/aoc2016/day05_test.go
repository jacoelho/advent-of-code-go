package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	aoc.AOCTest(t, day05p01, []aoc.TestInput{
		{Input: strings.NewReader("abc"), Want: "18f47a30"},
		{Input: aoc.FileInput(t, 2016, 5), Want: "f77a0e6e"},
	})
}

func Test_day05p02(t *testing.T) {
	aoc.AOCTest(t, day05p02, []aoc.TestInput{
		{Input: strings.NewReader("abc"), Want: "05ace8e3"},
		{Input: aoc.FileInput(t, 2016, 5), Want: "999828ec"},
	})
}
