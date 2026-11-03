package aoc2021

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const diracExample = "Player 1 starting position: 4\nPlayer 2 starting position: 8"

func Test_day21p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(diracExample), Want: "739785"},
		{Input: aoc.FileInput(t, 2021, 21), Want: "711480"},
	}

	aoc.AOCTest(t, day21p01, tests)
}

func Test_day21p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(diracExample), Want: "444356092776315"},
		{Input: aoc.FileInput(t, 2021, 21), Want: "265845890886828"},
	}

	aoc.AOCTest(t, day21p02, tests)
}

func TestDiracErrors(t *testing.T) {
	for _, input := range []string{
		"Player 1 starting position: 0\nPlayer 2 starting position: 8",
		"Player 1 starting position: 4\nPlayer 2 starting position: 11",
		"Player 1 starting position: 4",
		diracExample + "\nPlayer 3 starting position: 6",
	} {
		if _, err := day21p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}
