package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	aoc.AOCTest(t, day05p01, []aoc.TestInput{
		{Input: strings.NewReader("ugknbfddgicrmopn\naaa\njchzalrnumimnmhp\nhaegwjzuvuyypxyu\ndvszwmarrgswjxmb"), Want: "2"},
		{Input: aoc.FileInput(t, 2015, 5), Want: "258"},
	})
}

func Test_day05p02(t *testing.T) {
	aoc.AOCTest(t, day05p02, []aoc.TestInput{
		{Input: strings.NewReader("qjhvhtzxzqqjkmpb\nxxyxx\nuurcxstgmygtbstg\nieodomkazucvgmuy"), Want: "2"},
		{Input: strings.NewReader("aaa"), Want: "0"},
		{Input: strings.NewReader("aaaa"), Want: "1"},
		{Input: aoc.FileInput(t, 2015, 5), Want: "53"},
	})
}
