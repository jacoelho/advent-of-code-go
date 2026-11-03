package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day23p01(t *testing.T) {
	aoc.AOCTest(t, day23p01, []aoc.TestInput{
		{Input: strings.NewReader("inc b\njio b, +2\ntpl b\ninc b\n"), Want: "2"},
		{Input: strings.NewReader("inc b\njmp -2\n"), Want: "1"},
		{Input: aoc.FileInput(t, 2015, 23), Want: "255"},
	})
}

func Test_day23p02(t *testing.T) {
	aoc.AOCTest(t, day23p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2015, 23), Want: "334"}})
}

func TestTinyComputerHalvesAndBranches(t *testing.T) {
	got, err := tinyComputer(strings.NewReader("inc a\ntpl a\ninc a\nhlf a\njie a, +2\ninc b\ninc b"), 0)
	if err != nil || got != [2]uint64{2, 1} {
		t.Fatalf("registers = %v, %v; want [2 1]", got, err)
	}
	for _, input := range []string{"jmp +0", "inc c", "wat a", "jie a, x"} {
		if _, err := day23p01(strings.NewReader(input)); err == nil {
			t.Errorf("accepted %q", input)
		}
	}
}

func TestDay23ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day23p01,
		day23p02,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error = %v; want wrapped input failure", err)
		}
	}
}
