package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

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

func TestDay21ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"21/1", day21p01},
		{"21/2", day21p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
