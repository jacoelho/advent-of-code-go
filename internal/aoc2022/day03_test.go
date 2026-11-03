package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day03Example = `vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw`

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day03Example), Want: "157"},
		{Input: aoc.FileInput(t, 2022, 3), Want: "8233"},
	}
	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day03Example), Want: "70"},
		{Input: aoc.FileInput(t, 2022, 3), Want: "2821"},
	}
	aoc.AOCTest(t, day03p02, tests)
}

func TestRucksackInvalidInput(t *testing.T) {
	for _, input := range []string{"abc", "abAB", "abab", "11"} {
		if _, err := day03p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected item error for %q", input)
		}
	}
	if _, err := day03p02(strings.NewReader("a\na\na")); err == nil {
		t.Fatal("expected odd compartment length error")
	}
	if _, err := day03p02(strings.NewReader("aa\naa")); err == nil {
		t.Fatal("expected incomplete group error")
	}
}

func TestDay03ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day03p01},
		{"part2", day03p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
