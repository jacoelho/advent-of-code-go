package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestLookAndSayPhases(t *testing.T) {
	sequence := []byte("1")
	for _, want := range []string{"11", "21", "1211", "111221", "312211"} {
		sequence = lookAndSay(sequence)
		if string(sequence) != want {
			t.Fatalf("got %s, want %s", sequence, want)
		}
	}
}

func Test_day10p01(t *testing.T) {
	aoc.AOCTest(t, day10p01, []aoc.TestInput{
		{Input: strings.NewReader("1"), Want: "82350"},
		{Input: aoc.FileInput(t, 2015, 10), Want: "492982"},
	})
}

func Test_day10p02(t *testing.T) {
	aoc.AOCTest(t, day10p02, []aoc.TestInput{
		{Input: strings.NewReader("1"), Want: "1166642"},
		{Input: aoc.FileInput(t, 2015, 10), Want: "6989950"},
	})
}
