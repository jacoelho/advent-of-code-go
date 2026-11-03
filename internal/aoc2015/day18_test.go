package aoc2015

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const lightsExample = `.#.#.#
...##.
#....#
..#...
#.#..#
####..`

func Test_day18p01(t *testing.T) {
	aoc.AOCTest(t, day18p01, []aoc.TestInput{
		{Input: strings.NewReader("#"), Want: "0"},
		{Input: aoc.FileInput(t, 2015, 18), Want: "814"},
	})
}

func Test_day18p02(t *testing.T) {
	aoc.AOCTest(t, day18p02, []aoc.TestInput{
		{Input: strings.NewReader("."), Want: "1"},
		{Input: strings.NewReader("...\n..."), Want: "4"},
		{Input: aoc.FileInput(t, 2015, 18), Want: "924"},
	})
}

func TestLightExamples(t *testing.T) {
	for _, tt := range []struct {
		steps   int
		corners lightCornerMode
		want    string
	}{{4, normalLightCorners, "4"}, {5, stuckLightCorners, "17"}, {0, stuckLightCorners, "17"}} {
		if got, err := animateLights(strings.NewReader(lightsExample), tt.steps, tt.corners); err != nil || got != tt.want {
			t.Errorf("steps=%d corners=%d: got %s, %v; want %s", tt.steps, tt.corners, got, err, tt.want)
		}
	}
}

func TestLightErrors(t *testing.T) {
	for _, input := range []string{"", ".\n..", ".?"} {
		if _, err := day18p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day18p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day18p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
