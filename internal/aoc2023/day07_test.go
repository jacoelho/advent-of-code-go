package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483`),
			Want: "6440",
		},
		{
			Input: aoc.FileInput(t, 2023, 7),
			Want:  "253954294",
		},
	}
	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483`),
			Want: "5905",
		},
		{
			Input: aoc.FileInput(t, 2023, 7),
			Want:  "254837398",
		},
	}
	aoc.AOCTest(t, day07p02, tests)
}

func TestDay07RejectsMalformedHands(t *testing.T) {
	requireRejected(t, []string{
		"32T3K",
		"32T3K x",
		"32T3K -1",
		"32T3 765",
		"32T3KK 765",
		"32T3X 765",
	}, day07p01, day07p02)
}
