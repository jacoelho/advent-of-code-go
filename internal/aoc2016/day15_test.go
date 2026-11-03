package aoc2016

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const discExample = `Disc #1 has 5 positions; at time=0, it is at position 4.
Disc #2 has 2 positions; at time=0, it is at position 1.`

func Test_day15p01(t *testing.T) {
	aoc.AOCTest(t, day15p01, []aoc.TestInput{
		{Input: strings.NewReader(discExample), Want: "5"},
		{Input: strings.NewReader("Disc #5 has 7 positions; at time=0, it is at position 0."), Want: "2"},
		{Input: strings.NewReader("Disc #1 has 6 positions; at time=0, it is at position 5.\nDisc #2 has 8 positions; at time=0, it is at position 0."), Want: "6"},
		{Input: aoc.FileInput(t, 2016, 15), Want: "376777"},
	})
}

func Test_day15p02(t *testing.T) {
	aoc.AOCTest(t, day15p02, []aoc.TestInput{
		{Input: strings.NewReader(discExample), Want: "85"},
		{Input: strings.NewReader("Disc #5 has 7 positions; at time=0, it is at position 0."), Want: "16"},
		{Input: aoc.FileInput(t, 2016, 15), Want: "3903937"},
	})
}

func TestCapsuleDiscErrors(t *testing.T) {
	for _, input := range []string{"", "Disc #1 has 0 positions; at time=0, it is at position 0.", "Disc #1 has 6 positions; at time=0, it is at position 5.\nDisc #2 has 8 positions; at time=0, it is at position 1."} {
		if _, err := day15p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day15p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day15p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
