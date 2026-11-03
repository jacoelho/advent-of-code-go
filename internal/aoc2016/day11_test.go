package aoc2016

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const elevatorExample = `The first floor contains a hydrogen-compatible microchip and a lithium-compatible microchip.
The second floor contains a hydrogen generator.
The third floor contains a lithium generator.
The fourth floor contains nothing relevant.`

func Test_day11p01(t *testing.T) {
	aoc.AOCTest(t, day11p01, []aoc.TestInput{
		{Input: strings.NewReader(elevatorExample), Want: "11"},
		{Input: strings.NewReader(strings.NewReplacer("hydrogen", "cobalt", "lithium", "carbon").Replace(elevatorExample)), Want: "11"},
		{Input: strings.NewReader("The first floor contains a hydrogen generator and a hydrogen-compatible microchip.\nThe second floor contains nothing relevant.\nThe third floor contains nothing relevant.\nThe fourth floor contains nothing relevant."), Want: "3"},
		{Input: aoc.FileInput(t, 2016, 11), Want: "33"},
	})
}

func Test_day11p02(t *testing.T) {
	aoc.AOCTest(t, day11p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 11), Want: "57"}})
}

func TestElevatorErrors(t *testing.T) {
	unsafe := `The first floor contains a hydrogen-compatible microchip and a lithium generator.
The second floor contains a hydrogen generator and a lithium-compatible microchip.
The third floor contains nothing relevant.
The fourth floor contains nothing relevant.`
	for _, input := range []string{"", "floor 1", unsafe, strings.Replace(elevatorExample, "a lithium generator", "nothing relevant", 1)} {
		if _, err := day11p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day11p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day11p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
