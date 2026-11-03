package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	aoc.AOCTest(t, day02p01, []aoc.TestInput{
		{Input: strings.NewReader("ULL\nRRDDD\nLURDL\nUUUUD"), Want: "1985"},
		{Input: aoc.FileInput(t, 2016, 2), Want: "74921"},
	})
}

func Test_day02p02(t *testing.T) {
	aoc.AOCTest(t, day02p02, []aoc.TestInput{
		{Input: strings.NewReader("ULL\nRRDDD\nLURDL\nUUUUD"), Want: "5DB3"},
		{Input: aoc.FileInput(t, 2016, 2), Want: "A6B35"},
	})
}
