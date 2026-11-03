package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	aoc.AOCTest(t, day20p01, []aoc.TestInput{
		{Input: strings.NewReader("11"), Want: "2"},
		{Input: strings.NewReader("150"), Want: "8"},
		{Input: aoc.FileInput(t, 2015, 20), Want: "786240"},
	})
}

func Test_day20p02(t *testing.T) {
	aoc.AOCTest(t, day20p02, []aoc.TestInput{
		{Input: strings.NewReader("12"), Want: "2"},
		{Input: aoc.FileInput(t, 2015, 20), Want: "831600"},
	})
}

func TestElfVisitLimit(t *testing.T) {
	got, err := firstHouse(strings.NewReader("70"), 1, 2)
	if err != nil || got != "48" {
		t.Fatalf("house = %q, %v; want48", got, err)
	}
}

func TestDay20ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day20p01,
		day20p02,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error = %v; want wrapped input failure", err)
		}
	}
}

func TestInvalidPresentTargets(t *testing.T) {
	for _, input := range []string{"", "0", "-1", "-", "abc", "9223372036854775808"} {
		if _, err := day20p01(strings.NewReader(input)); err == nil {
			t.Errorf("accepted target %q", input)
		}
	}
}
