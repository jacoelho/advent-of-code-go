package aoc2015

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const seatingExample = `Alice would gain 54 happiness units by sitting next to Bob.
Alice would lose 79 happiness units by sitting next to Carol.
Alice would lose 2 happiness units by sitting next to David.
Bob would gain 83 happiness units by sitting next to Alice.
Bob would lose 7 happiness units by sitting next to Carol.
Bob would lose 63 happiness units by sitting next to David.
Carol would lose 62 happiness units by sitting next to Alice.
Carol would gain 60 happiness units by sitting next to Bob.
Carol would gain 55 happiness units by sitting next to David.
David would gain 46 happiness units by sitting next to Alice.
David would lose 7 happiness units by sitting next to Bob.
David would gain 41 happiness units by sitting next to Carol.`

const negativeSeating = `A would lose 1 happiness units by sitting next to B.
A would lose 1 happiness units by sitting next to C.
B would lose 1 happiness units by sitting next to A.
B would lose 1 happiness units by sitting next to C.
C would lose 1 happiness units by sitting next to A.
C would lose 1 happiness units by sitting next to B.`

func Test_day13p01(t *testing.T) {
	aoc.AOCTest(t, day13p01, []aoc.TestInput{
		{Input: strings.NewReader(seatingExample), Want: "330"},
		{Input: strings.NewReader(negativeSeating), Want: "-6"},
		{Input: aoc.FileInput(t, 2015, 13), Want: "664"},
	})
}

func Test_day13p02(t *testing.T) {
	aoc.AOCTest(t, day13p02, []aoc.TestInput{
		{Input: strings.NewReader(seatingExample), Want: "286"},
		{Input: strings.NewReader(negativeSeating), Want: "-4"},
		{Input: aoc.FileInput(t, 2015, 13), Want: "640"},
	})
}

func TestSeatingErrors(t *testing.T) {
	for _, input := range []string{"", "A likes B", "A would gain 1 happiness units by sitting next to B.", seatingExample + "\n" + strings.Split(seatingExample, "\n")[0]} {
		if _, err := day13p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day13p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day13p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
