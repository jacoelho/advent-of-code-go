package aoc2016

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	aoc.AOCTest(t, day12p01, []aoc.TestInput{
		{Input: strings.NewReader("cpy 41 a\ninc a\ninc a\ndec a\njnz a 2\ndec a"), Want: "42"},
		{Input: strings.NewReader("cpy c a"), Want: "0"},
		{Input: aoc.FileInput(t, 2016, 12), Want: "318077"},
	})
}

func Test_day12p02(t *testing.T) {
	aoc.AOCTest(t, day12p02, []aoc.TestInput{
		{Input: strings.NewReader("cpy c a"), Want: "1"},
		{Input: aoc.FileInput(t, 2016, 12), Want: "9227731"},
	})
}

func TestDay12ReaderErrors(t *testing.T) {
	failure := errors.New("read failed")
	if _, err := day12p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day12p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
