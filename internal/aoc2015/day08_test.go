package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day08p01(t *testing.T) {
	aoc.AOCTest(t, day08p01, []aoc.TestInput{
		{Input: strings.NewReader("\"\"\n\"abc\"\n\"aaa\\\"aaa\"\n\"\\x27\""), Want: "12"},
		{Input: aoc.FileInput(t, 2015, 8), Want: "1333"},
	})
}

func Test_day08p02(t *testing.T) {
	aoc.AOCTest(t, day08p02, []aoc.TestInput{
		{Input: strings.NewReader("\"\"\n\"abc\"\n\"aaa\\\"aaa\"\n\"\\x27\""), Want: "19"},
		{Input: aoc.FileInput(t, 2015, 8), Want: "2046"},
	})
}
