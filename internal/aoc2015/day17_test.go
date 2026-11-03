package aoc2015

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	aoc.AOCTest(t, day17p01, []aoc.TestInput{
		{Input: strings.NewReader("150\n150"), Want: "2"},
		{Input: strings.NewReader("1"), Want: "0"},
		{Input: aoc.FileInput(t, 2015, 17), Want: "4372"},
	})
}

func Test_day17p02(t *testing.T) {
	aoc.AOCTest(t, day17p02, []aoc.TestInput{
		{Input: strings.NewReader("150\n150\n100\n50"), Want: "2"},
		{Input: strings.NewReader("1"), Want: "0"},
		{Input: aoc.FileInput(t, 2015, 17), Want: "4"},
	})
}

func TestContainerExample(t *testing.T) {
	all, minimal, err := containerCombinations(strings.NewReader("20\n15\n10\n5\n5"), 25)
	if err != nil {
		t.Fatal(err)
	}
	if all.String() != "4" || minimal.String() != "3" {
		t.Errorf("got all=%s, minimal=%s; want 4, 3", all, minimal)
	}
}

func TestContainerErrors(t *testing.T) {
	for _, input := range []string{"", "0", "-1", "20\nfive"} {
		if _, err := day17p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day17p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day17p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
