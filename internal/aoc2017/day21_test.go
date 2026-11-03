package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day21p01(t *testing.T) {
	aoc.AOCTest(t, day21p01, []aoc.TestInput{
		{Input: aoc.FileInput(t, 2017, 21), Want: "179"},
	})
}

func Test_day21p02(t *testing.T) {
	aoc.AOCTest(t, day21p02, []aoc.TestInput{
		{Input: aoc.FileInput(t, 2017, 21), Want: "2766750"},
	})
}

func TestEnhanceFractal(t *testing.T) {
	input := "../.# => ##./#../...\n.#./..#/### => #..#/..../..../#..#"
	got, err := enhanceFractal(strings.NewReader(input), 2)
	if err != nil || got != 12 {
		t.Fatalf("enhanceFractal() = %d, %v; want 12", got, err)
	}
}

func TestFractalRuleOrientations(t *testing.T) {
	patterns := []string{
		".#./..#/###", "#../#.#/##.", "###/#../.#.", ".##/#.#/..#",
		".#./#../###", "..#/#.#/.##", "###/..#/.#.", "##./#.#/#..",
	}
	for _, pattern := range patterns {
		t.Run(pattern, func(t *testing.T) {
			got, err := enhanceFractal(strings.NewReader(pattern+" => #.../..../..../...."), 1)
			if err != nil || got != 1 {
				t.Fatalf("enhanceFractal() = %d, %v; want 1", got, err)
			}
		})
	}
}

func TestFractalRulesRejectConflictingOrientations(t *testing.T) {
	input := "../.# => .../.../...\n#./.. => ###/###/###"
	if _, err := parseFractalRules(strings.NewReader(input)); err == nil {
		t.Fatal("expected conflicting rule error")
	}
}
