package aoc2016

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day13p01(t *testing.T) {
	aoc.AOCTest(t, day13p01, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 13), Want: "82"}})
}

func Test_day13p02(t *testing.T) {
	aoc.AOCTest(t, day13p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 13), Want: "138"}})
}

func TestOfficeMazeExample(t *testing.T) {
	if got, _, err := officeMaze(strings.NewReader("10"), officePosition{7, 4}, -1); err != nil || got != 11 {
		t.Errorf("distance=%d, %v; want 11", got, err)
	}
	for _, tt := range []struct{ distance, want int }{{0, 1}, {1, 3}, {2, 5}} {
		if _, got, err := officeMaze(strings.NewReader("10"), officePosition{}, tt.distance); err != nil || got != tt.want {
			t.Errorf("reachable within %d=%d, %v; want %d", tt.distance, got, err, tt.want)
		}
	}
}

func TestOfficeMazeErrors(t *testing.T) {
	for _, input := range []string{"", "-1", "x", "0"} {
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
