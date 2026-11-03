package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day04p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`aa bb cc dd ee
aa bb cc dd aa
aa bb cc dd aaa`), Want: "2"},
		{Input: aoc.FileInput(t, 2017, 4), Want: "337"},
	}

	aoc.AOCTest(t, day04p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`abcde fghij
abcde xyz ecdab
a ab abc abd abf abj
iiii oiii ooii oooi oooo
oiii ioii iioi iiio`), Want: "3"},
		{Input: strings.NewReader("ab abb"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 4), Want: "231"},
	}

	aoc.AOCTest(t, day04p02, tests)
}
