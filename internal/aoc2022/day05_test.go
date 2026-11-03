package aoc2022

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day05Example = `    [D]
[N] [C]
[Z] [M] [P]
 1   2   3

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2`

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day05Example), Want: "CMZ"},
		{Input: aoc.FileInput(t, 2022, 5), Want: "FZCMJCRHZ"},
	}
	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day05Example), Want: "MCD"},
		{Input: aoc.FileInput(t, 2022, 5), Want: "JSDHQMZGF"},
	}
	aoc.AOCTest(t, day05p02, tests)
}

func TestCratesInvalidInput(t *testing.T) {
	for _, input := range []string{
		"[A] [B]\n 1   1 \n\nmove 0 from 1 to 1",
		"[A] [B]\n 1   2 \n\nmove 2 from 1 to 2",
		"[A] [B]\n 1   2 \n\nmove 1 from 1 to 2",
	} {
		if _, err := day05p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected crate error for %q", input)
		}
	}
}
