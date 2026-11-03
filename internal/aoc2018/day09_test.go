package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("9 players; last marble is worth 25 points"), Want: "32"},
		{Input: strings.NewReader("9 players; last marble is worth 23 points"), Want: "32"},
		{Input: strings.NewReader("10 players; last marble is worth 1618 points"), Want: "8317"},
		{Input: aoc.FileInput(t, 2018, 9), Want: "386018"},
	}

	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("9 players; last marble is worth 25 points"), Want: "22563"},
		{Input: aoc.FileInput(t, 2018, 9), Want: "3085518618"},
	}

	aoc.AOCTest(t, day09p02, tests)
}
