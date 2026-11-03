package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("dabAcCaCBAcCcaDA"), Want: "10"},
		{Input: strings.NewReader("abBA"), Want: "0"},
		{Input: strings.NewReader("abAB"), Want: "4"},
		{Input: strings.NewReader(""), Want: "0"},
		{Input: aoc.FileInput(t, 2018, 5), Want: "11636"},
	}

	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("dabAcCaCBAcCcaDA"), Want: "4"},
		{Input: strings.NewReader(""), Want: "0"},
		{Input: aoc.FileInput(t, 2018, 5), Want: "5302"},
	}

	aoc.AOCTest(t, day05p02, tests)
}
