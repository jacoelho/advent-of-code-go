package aoc2015

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const reindeerExample = `Comet can fly 14 km/s for 10 seconds, but then must rest for 127 seconds.
Dancer can fly 16 km/s for 11 seconds, but then must rest for 162 seconds.`

func Test_day14p01(t *testing.T) {
	aoc.AOCTest(t, day14p01, []aoc.TestInput{{Input: aoc.FileInput(t, 2015, 14), Want: "2660"}})
}

func Test_day14p02(t *testing.T) {
	aoc.AOCTest(t, day14p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2015, 14), Want: "1256"}})
}

func TestReindeerExamples(t *testing.T) {
	if got, err := reindeerDistance(strings.NewReader(reindeerExample), 1000); err != nil || got != "1120" {
		t.Errorf("distance = %s, %v; want 1120", got, err)
	}
	if got, err := reindeerPoints(strings.NewReader(reindeerExample), 1000); err != nil || got != "689" {
		t.Errorf("points = %s, %v; want 689", got, err)
	}
}

func TestReindeerPointsIncludeTiesAndFinalSecond(t *testing.T) {
	input := `A can fly 2 km/s for 1 seconds, but then must rest for 1 seconds.
B can fly 1 km/s for 10 seconds, but then must rest for 1 seconds.`
	first, second, _ := strings.Cut(input, "\n")
	for _, order := range []string{input, second + "\n" + first} {
		for _, tt := range []struct {
			seconds int
			want    string
		}{{1, "1"}, {3, "3"}} {
			if got, err := reindeerPoints(strings.NewReader(order), tt.seconds); err != nil || got != tt.want {
				t.Errorf("points after %d seconds = %s, %v; want %s", tt.seconds, got, err, tt.want)
			}
		}
	}
}

func TestReindeerErrors(t *testing.T) {
	for _, input := range []string{"", "Comet flies fast", "Comet can fly 14 km/s for 0 seconds, but then must rest for 1 seconds."} {
		if _, err := day14p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day14p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day14p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
