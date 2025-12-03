package aoc2025

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`),
			Want:  "1227775554",
		},
		{
			Input: aoc.FileInput(t, 2025, 02),
			Want:  "24157613387",
		},
	}
	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`),
			Want:  "4174379265",
		},
		{
			Input: aoc.FileInput(t, 2025, 02),
			Want:  "33832678380",
		},
	}
	aoc.AOCTest(t, day02p02, tests)
}
