package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day19p01(t *testing.T) {
	aoc.AOCTest(t, day19p01, []aoc.TestInput{{Input: strings.NewReader("5"), Want: "3"}, {Input: strings.NewReader("8"), Want: "1"}, {Input: strings.NewReader("9"), Want: "3"}, {Input: strings.NewReader("1"), Want: "1"}, {Input: aoc.FileInput(t, 2016, 19), Want: "1815603"}})
}
func Test_day19p02(t *testing.T) {
	aoc.AOCTest(t, day19p02, []aoc.TestInput{{Input: strings.NewReader("5"), Want: "2"}, {Input: strings.NewReader("8"), Want: "7"}, {Input: strings.NewReader("9"), Want: "9"}, {Input: strings.NewReader("10"), Want: "1"}, {Input: strings.NewReader("1"), Want: "1"}, {Input: aoc.FileInput(t, 2016, 19), Want: "1410630"}})
}
