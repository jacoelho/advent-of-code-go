package aoc2015

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day22p01(t *testing.T) {
	aoc.AOCTest(t, day22p01, []aoc.TestInput{{Input: aoc.FileInput(t, 2015, 22), Want: "1269"}})
}

func Test_day22p02(t *testing.T) {
	aoc.AOCTest(t, day22p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2015, 22), Want: "1309"}})
}

func TestWizardExamples(t *testing.T) {
	for _, tc := range []struct{ boss, want int }{{13, 226}, {14, 641}} {
		got, err := leastMana(wizardState{hp: 10, mana: 250, boss: tc.boss}, 8, wizardNormal)
		if err != nil || got != tc.want {
			t.Errorf("boss %d: mana = %d, %v; want %d", tc.boss, got, err, tc.want)
		}
	}
}

func TestWizardHardModeBeforeEffects(t *testing.T) {
	state := wizardState{hp: 1, boss: 3, poison: 1}
	if _, err := leastMana(state, 8, wizardHard); err == nil {
		t.Fatal("hard mode must kill the player before poison takes effect")
	}
	if got, err := leastMana(state, 8, wizardNormal); err != nil || got != 0 {
		t.Fatalf("poison win = %d, %v; want0", got, err)
	}
}

func TestWizardEffectExpiresBeforeCasting(t *testing.T) {
	got, err := leastMana(wizardState{hp: 1, mana: 72, boss: 4, recharge: 1, shield: 1}, 8, wizardNormal)
	if err != nil || got != 53 {
		t.Fatalf("recharge then missile = %d, %v; want53", got, err)
	}
}
