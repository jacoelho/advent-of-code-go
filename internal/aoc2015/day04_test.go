package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day04p01(t *testing.T) {
	aoc.AOCTest(t, day04p01, []aoc.TestInput{
		{Input: strings.NewReader("abcdef"), Want: "609043"},
		{Input: strings.NewReader("pqrstuv"), Want: "1048970"},
		{Input: aoc.FileInput(t, 2015, 4), Want: "254575"},
	})
}

func Test_day04p02(t *testing.T) {
	aoc.AOCTest(t, day04p02, []aoc.TestInput{
		{Input: strings.NewReader("abcdef"), Want: "6742839"},
		{Input: aoc.FileInput(t, 2015, 4), Want: "1038736"},
	})
}
