package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const generatorExample = "Generator A starts with 65\nGenerator B starts with 8921"

func Test_day15p01(t *testing.T) {
	if testing.Short() {
		t.Skip("full 40-million-pair generator fixtures")
	}

	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(generatorExample),
			Want:  "588",
		},
		{
			Input: aoc.FileInput(t, 2017, 15),
			Want:  "594",
		},
	}

	aoc.AOCTest(t, day15p01, tests)
}

func Test_day15p02(t *testing.T) {
	if testing.Short() {
		t.Skip("full filtered generator fixtures")
	}

	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(generatorExample),
			Want:  "309",
		},
		{
			Input: aoc.FileInput(t, 2017, 15),
			Want:  "328",
		},
	}

	aoc.AOCTest(t, day15p02, tests)
}

func TestGeneratorExamples(t *testing.T) {
	tests := []struct {
		pairs    int
		divisors [2]uint64
		want     string
	}{
		{pairs: 5, divisors: [2]uint64{1, 1}, want: "1"},
		{pairs: 1055, divisors: [2]uint64{4, 8}, want: "0"},
		{pairs: 1056, divisors: [2]uint64{4, 8}, want: "1"},
	}

	for _, tt := range tests {
		got, err := judgeGenerators(strings.NewReader(generatorExample), tt.pairs, tt.divisors)
		if err != nil || got != tt.want {
			t.Errorf("pairs=%d divisors=%v: got %s, %v; want %s", tt.pairs, tt.divisors, got, err, tt.want)
		}
	}
}

func TestGeneratorErrors(t *testing.T) {
	for _, input := range []string{"", "Generator A starts with 1", "Generator A starts with -1\nGenerator B starts with 1", "Generator A starts with 2147483647\nGenerator B starts with 1"} {
		if _, err := judgeGenerators(strings.NewReader(input), 1, [2]uint64{1, 1}); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}
