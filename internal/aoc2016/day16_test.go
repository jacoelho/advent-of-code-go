package aoc2016

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day16p01(t *testing.T) {
	aoc.AOCTest(t, day16p01, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 16), Want: "00000100100001100"}})
}

func Test_day16p02(t *testing.T) {
	aoc.AOCTest(t, day16p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 16), Want: "00011010100010010"}})
}

func TestDragonChecksumExamples(t *testing.T) {
	for _, tt := range []struct {
		seed string
		size int
		want string
	}{{"10000", 20, "01100"}, {"110010110100", 12, "100"}, {"101101", 3, "101"}, {"1", 3, "100"}, {"1", 1, "1"}} {
		if got, err := dragonChecksum(strings.NewReader(tt.seed), tt.size); err != nil || got != tt.want {
			t.Errorf("seed=%s size=%d: got %s, %v; want %s", tt.seed, tt.size, got, err, tt.want)
		}
	}
}

func TestDragonErrors(t *testing.T) {
	for _, input := range []string{"", "1012", "10 01"} {
		if _, err := day16p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day16p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day16p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
