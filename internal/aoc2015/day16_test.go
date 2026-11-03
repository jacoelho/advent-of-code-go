package aoc2015

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day16p01(t *testing.T) {
	aoc.AOCTest(t, day16p01, []aoc.TestInput{
		{Input: strings.NewReader("Sue 1: cats: 0\nSue 2: cats: 7"), Want: "2"},
		{Input: strings.NewReader("Sue 3: akitas: 0, cars: 2, perfumes: 1"), Want: "3"},
		{Input: aoc.FileInput(t, 2015, 16), Want: "103"},
	})
}

func Test_day16p02(t *testing.T) {
	aoc.AOCTest(t, day16p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2015, 16), Want: "405"}})
}

func TestAuntStrictRanges(t *testing.T) {
	for _, tt := range []struct {
		property        string
		equal, matching int
	}{{"cats", 7, 8}, {"trees", 3, 4}, {"pomeranians", 3, 2}, {"goldfish", 5, 4}} {
		t.Run(tt.property, func(t *testing.T) {
			input := fmt.Sprintf("Sue 1: %s: %d\nSue 2: %s: %d", tt.property, tt.equal, tt.property, tt.matching)
			if got, err := day16p02(strings.NewReader(input)); err != nil || got != "2" {
				t.Errorf("got %s, %v; want 2", got, err)
			}
		})
	}
}

func TestAuntErrors(t *testing.T) {
	for _, input := range []string{"", "Sue 1: cats: 0", "Sue 1: cats: 7\nSue 2: cats: 7", "Sue 1: cats: 7, cats: 7", "Sue 1: unknown: 2", "Sue 1: cats: -1"} {
		if _, err := day16p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day16p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day16p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
