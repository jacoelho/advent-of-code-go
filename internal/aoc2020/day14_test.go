package aoc2020

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day14p01(t *testing.T) {
	aoc.AOCTest(t, day14p01, []aoc.TestInput{
		{Input: strings.NewReader("mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X\nmem[8] = 11\nmem[7] = 101\nmem[8] = 0"), Want: "165"},
		{Input: aoc.FileInput(t, 2020, 14), Want: "11612740949946"},
	})
}

func Test_day14p02(t *testing.T) {
	aoc.AOCTest(t, day14p02, []aoc.TestInput{
		{Input: strings.NewReader("mask = 000000000000000000000000000000X1001X\nmem[42] = 100\nmask = 00000000000000000000000000000000X0XX\nmem[26] = 1"), Want: "208"},
		{Input: strings.NewReader("mask = 000000000000000000000000000000000000\nmem[0] = 1\nmem[0] = 2\nmem[1] = 3"), Want: "5"},
		{Input: strings.NewReader("mask = 0000000000000000000000000000000000XX\nmem[0] = 3\nmask = 0000000000000000000000000000000000X1\nmem[0] = 5"), Want: "16"},
		{Input: strings.NewReader("mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\nmem[0] = 1"), Want: "68719476736"},
		{Input: strings.NewReader("mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\nmem[0] = 1\nmem[1] = 2"), Want: "137438953472"},
		{Input: strings.NewReader("mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\nmem[0] = 1\nmem[1] = 0"), Want: "0"},
		{Input: aoc.FileInput(t, 2020, 14), Want: "3394509207186"},
	})
}

func TestDockingProgramErrors(t *testing.T) {
	for _, input := range []string{"mem[0] = 1", "mask = XX", "mask = 00000000000000000000000000000000000Z", "mask = 000000000000000000000000000000000000\nmem[x] = 1"} {
		if _, err := day14p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day14p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
	if _, err := day14p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
}
