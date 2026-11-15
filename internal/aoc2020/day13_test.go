package aoc2020

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day13p01(t *testing.T) {
	aoc.AOCTest(t, day13p01, []aoc.TestInput{
		{Input: strings.NewReader("939\n7,13,x,x,59,x,31,19"), Want: "295"},
		{Input: strings.NewReader("10\n5"), Want: "0"},
		{Input: aoc.FileInput(t, 2020, 13), Want: "136"},
	})
}

func Test_day13p02(t *testing.T) {
	aoc.AOCTest(t, day13p02, []aoc.TestInput{
		{Input: strings.NewReader("939\n7,13,x,x,59,x,31,19"), Want: "1068781"},
		{Input: strings.NewReader("0\nx,3,5"), Want: "8"},
		{Input: strings.NewReader("0\n6,x,8"), Want: "6"},
		{Input: strings.NewReader("0\n3,x,x,3"), Want: "0"},
		{Input: strings.NewReader("0\n17,x,13,19"), Want: "3417"},
		{Input: aoc.FileInput(t, 2020, 13), Want: "305068317272992"},
	})
}

func TestBusScheduleErrors(t *testing.T) {
	for _, input := range []string{"", "0\nx,x", "0\n0", "0\n6,8", "0\n2,3,5,7,11,13,17,19,23,29,31,37,41,43,47,53"} {
		if _, err := day13p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	if _, err := day13p01(strings.NewReader("1\n4611686018427387904")); err == nil {
		t.Error("expected overflow error for a wait times bus beyond int64")
	}
	failure := errors.New("read failed")
	if _, err := day13p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
	if _, err := day13p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Fatalf("got %v, want read error", err)
	}
}
