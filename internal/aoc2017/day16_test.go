package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day16p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2017, 16),
			Want:  "bijankplfgmeodhc",
		},
	}

	aoc.AOCTest(t, day16p01, tests)
}

func Test_day16p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2017, 16),
			Want:  "bpjahknliomefdgc",
		},
	}

	aoc.AOCTest(t, day16p02, tests)
}

func TestDanceExamples(t *testing.T) {
	tests := []struct {
		rounds uint64
		want   string
	}{
		{rounds: 0, want: "abcde"},
		{rounds: 1, want: "baedc"},
		{rounds: 2, want: "ceadb"},
		{rounds: 4, want: "abcde"},
		{rounds: 1_000_000_000, want: "abcde"},
	}

	for _, tt := range tests {
		got, err := dancePrograms(strings.NewReader("s1,x3/4,pe/b"), 5, tt.rounds)
		if err != nil || got != tt.want {
			t.Errorf("rounds=%d: got %s, %v; want %s", tt.rounds, got, err, tt.want)
		}
	}
}

func TestDanceErrors(t *testing.T) {
	for _, input := range []string{"", "s-1", "x0/16", "pa/z", "s1,", "x1"} {
		if _, err := day16p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}
