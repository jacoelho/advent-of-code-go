package aoc2015

import (
	"errors"
	"io"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day21p01(t *testing.T) {
	aoc.AOCTest(t, day21p01, []aoc.TestInput{{Input: aoc.FileInput(t, 2015, 21), Want: "111"}})
}

func Test_day21p02(t *testing.T) {
	aoc.AOCTest(t, day21p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2015, 21), Want: "188"}})
}

func TestEquipmentCombat(t *testing.T) {
	if !playerWins(combatStats{8, 5, 5}, combatStats{12, 7, 2}) {
		t.Fatal("player should win the four-turn example")
	}
	if !playerWins(combatStats{1, 1, 100}, combatStats{1, 1, 100}) {
		t.Fatal("player attacks first; damage is at least one")
	}
	if playerWins(combatStats{1, 1, 100}, combatStats{2, 1, 100}) {
		t.Fatal("boss survives the first hit")
	}
}

func TestDay21ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day21p01,
		day21p02,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error = %v; want wrapped input failure", err)
		}
	}
}
