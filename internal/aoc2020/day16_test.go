package aoc2020

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day16p01(t *testing.T) {
	aoc.AOCTest(t, day16p01, []aoc.TestInput{
		{Input: strings.NewReader(`class: 1-3 or 5-7
row: 6-11 or 33-44
seat: 13-40 or 45-50

your ticket:
7,1,14

nearby tickets:
7,3,47
40,4,50
55,2,20
38,6,12`), Want: "71"},
		{Input: aoc.FileInput(t, 2020, 16), Want: "26053"},
	})
}

const departureTicketExample = `departure class: 0-1 or 4-19
departure row: 0-5 or 8-19
seat: 0-13 or 16-19

your ticket:
11,12,13

nearby tickets:
3,9,18
15,1,5
5,14,9`

func Test_day16p02(t *testing.T) {
	aoc.AOCTest(t, day16p02, []aoc.TestInput{
		{Input: strings.NewReader(departureTicketExample), Want: "132"},
		{Input: strings.NewReader(departureTicketExample + "\n999,999,999"), Want: "132"},
		{Input: strings.NewReader(strings.Replace(departureTicketExample, "departure class:", "nondeparture class:", 1)), Want: "11"},
		{Input: aoc.FileInput(t, 2020, 16), Want: "1515506256421"},
	})
}

func TestTicketAssignmentErrors(t *testing.T) {
	for _, input := range []string{
		"a: 0-1 or 3-4\nb: 0-1 or 3-4\n\nyour ticket:\n0,1\n\nnearby tickets:\n0,1",
		"a: 0-0 or 0-0\nb: 0-0 or 0-0\nc: 1-2 or 1-2\n\nyour ticket:\n0,1,2\n\nnearby tickets:\n0,1,2",
		strings.Replace(departureTicketExample, "11,12,13", "999,12,13", 1),
		strings.Replace(departureTicketExample, "11,12,13", "11,12", 1),
		"invalid",
	} {
		if _, err := day16p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected invalid or ambiguous ticket error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day16p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
	if _, err := day16p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
}
