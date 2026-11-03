package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("ne,ne,ne"), Want: "3"},
		{Input: strings.NewReader("ne,ne,sw,sw"), Want: "0"},
		{Input: strings.NewReader("ne,ne,s,s"), Want: "2"},
		{Input: strings.NewReader("se,sw,se,sw,sw"), Want: "3"},
		{Input: strings.NewReader("ne,n"), Want: "2"},
		{Input: aoc.FileInput(t, 2017, 11), Want: "705"},
	}

	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("ne,ne,sw,sw"), Want: "2"},
		{Input: strings.NewReader("ne,n"), Want: "2"},
		{Input: aoc.FileInput(t, 2017, 11), Want: "1469"},
	}

	aoc.AOCTest(t, day11p02, tests)
}
