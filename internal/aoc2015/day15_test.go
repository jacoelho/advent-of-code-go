package aoc2015

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const cookieExample = `Butterscotch: capacity -1, durability -2, flavor 6, texture 3, calories 8
Cinnamon: capacity 2, durability 3, flavor -2, texture -1, calories 3`

const zeroQuantityCookie = `Good: capacity 1, durability 1, flavor 1, texture 1, calories 5
Bad: capacity -1, durability -1, flavor -1, texture -1, calories 5`

func Test_day15p01(t *testing.T) {
	aoc.AOCTest(t, day15p01, []aoc.TestInput{
		{Input: strings.NewReader(cookieExample), Want: "62842880"},
		{Input: strings.NewReader(zeroQuantityCookie), Want: "100000000"},
		{Input: strings.NewReader("Bad: capacity -1, durability -1, flavor -1, texture -1, calories 5"), Want: "0"},
		{Input: aoc.FileInput(t, 2015, 15), Want: "13882464"},
	})
}

func Test_day15p02(t *testing.T) {
	aoc.AOCTest(t, day15p02, []aoc.TestInput{
		{Input: strings.NewReader(cookieExample), Want: "57600000"},
		{Input: strings.NewReader(zeroQuantityCookie), Want: "100000000"},
		{Input: strings.NewReader("Good: capacity 1, durability 1, flavor 1, texture 1, calories 6"), Want: "0"},
		{Input: aoc.FileInput(t, 2015, 15), Want: "11171160"},
	})
}

func TestCookieErrors(t *testing.T) {
	for _, input := range []string{"", "Sugar: 1, 2, 3, 4, 5", cookieExample + "\n" + strings.Split(cookieExample, "\n")[0]} {
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
