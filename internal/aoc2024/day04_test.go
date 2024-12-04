package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day04p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`),
			Want: "18",
		},
		{
			Input: aoc.FileInput(t, 2024, 4),
			Want:  "2468",
		},
	}

	aoc.AOCTest(t, day04p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`),
			Want: "9",
		},
		{
			Input: aoc.FileInput(t, 2024, 4),
			Want:  "1864",
		},
	}

	aoc.AOCTest(t, day04p02, tests)
}
