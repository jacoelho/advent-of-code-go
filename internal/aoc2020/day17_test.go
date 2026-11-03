package aoc2020

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	aoc.AOCTest(t, day17p01, []aoc.TestInput{
		{Input: strings.NewReader(".#.\n..#\n###"), Want: "112"},
		{Input: strings.NewReader("#"), Want: "0"},
		{Input: aoc.FileInput(t, 2020, 17), Want: "315"},
	})
}

func Test_day17p02(t *testing.T) {
	aoc.AOCTest(t, day17p02, []aoc.TestInput{
		{Input: strings.NewReader(".#.\n..#\n###"), Want: "848"},
		{Input: strings.NewReader("#"), Want: "0"},
		{Input: aoc.FileInput(t, 2020, 17), Want: "1520"},
	})
}

func TestPocketDimensionErrors(t *testing.T) {
	for _, input := range []string{"", "..\n#", "x"} {
		if _, err := day17p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day17p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
	if _, err := day17p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
}
