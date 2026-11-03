package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const licenseExample = "2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2"

func Test_day08p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(licenseExample), Want: "138"},
		{Input: aoc.FileInput(t, 2018, 8), Want: "38567"},
	}

	aoc.AOCTest(t, day08p01, tests)
}

func Test_day08p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(licenseExample), Want: "66"},
		{Input: strings.NewReader("1 4 0 1 5 0 2 1 1"), Want: "10"},
		{Input: aoc.FileInput(t, 2018, 8), Want: "24453"},
	}

	aoc.AOCTest(t, day08p02, tests)
}
