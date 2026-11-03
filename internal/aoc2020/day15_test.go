package aoc2020

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day15p01(t *testing.T) {
	aoc.AOCTest(t, day15p01, []aoc.TestInput{
		{Input: strings.NewReader("0,3,6"), Want: "436"},
		{Input: aoc.FileInput(t, 2020, 15), Want: "289"},
	})
}

func Test_day15p02(t *testing.T) {
	aoc.AOCTest(t, day15p02, []aoc.TestInput{
		{Input: strings.NewReader("0,3,6"), Want: "175594"},
		{Input: aoc.FileInput(t, 2020, 15), Want: "1505722"},
	})
}

func TestMemoryGameStartingHistory(t *testing.T) {
	for _, tt := range []struct {
		starting   []int
		turn, want int
	}{
		{[]int{1, 2, 1}, 4, 2}, {[]int{100, 1, 100}, 4, 2}, {[]int{1, 2, 1}, 2, 2},
		{[]int{0, 3, 6}, 10, 0},
	} {
		if got := memoryGameNumber(tt.starting, tt.turn); got != tt.want {
			t.Errorf("%v turn %d: got %d, want %d", tt.starting, tt.turn, got, tt.want)
		}
	}
}

func TestMemoryGameErrors(t *testing.T) {
	for _, input := range []string{"", "1,-1", "1,x", "1,"} {
		if _, err := day15p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day15p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
	if _, err := day15p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
}
