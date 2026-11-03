package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestArmyStalemate(t *testing.T) {
	input := "Immune System:\n1 units each with 10 hit points with an attack that does 0 fire damage at initiative 2\n\nInfection:\n1 units each with 10 hit points with an attack that does 0 cold damage at initiative 1"
	if _, err := day24p01(strings.NewReader(input)); err == nil {
		t.Fatal("expected stalemate error")
	}
	got, err := day24p02(strings.NewReader(input))
	if err != nil || got != "1" {
		t.Fatalf("day24p02() = %q, %v; want 1", got, err)
	}
}

func TestArmyImmunityOverridesWeakness(t *testing.T) {
	input := "Immune System:\n1 units each with 10 hit points (weak to cold; immune to cold) with an attack that does 10 fire damage at initiative 2\n\nInfection:\n1 units each with 10 hit points (weak to fire; immune to fire) with an attack that does 10 cold damage at initiative 1"
	if _, err := day24p01(strings.NewReader(input)); err == nil {
		t.Fatal("expected mutually immune armies to reach a stalemate")
	}
}

func Test_day24p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("Immune System:\n17 units each with 5390 hit points (weak to radiation, bludgeoning) with an attack that does 4507 fire damage at initiative 2\n989 units each with 1274 hit points (immune to fire; weak to bludgeoning, slashing) with an attack that does 25 slashing damage at initiative 3\n\nInfection:\n801 units each with 4706 hit points (weak to radiation) with an attack that does 116 bludgeoning damage at initiative 1\n4485 units each with 2961 hit points (immune to radiation; weak to fire, cold) with an attack that does 12 slashing damage at initiative 4"),
			Want:  "5216",
		},
		{
			Input: aoc.FileInput(t, 2018, 24),
			Want:  "22859",
		},
	}
	aoc.AOCTest(t, day24p01, tests)
}

func Test_day24p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("Immune System:\n17 units each with 5390 hit points (weak to radiation, bludgeoning) with an attack that does 4507 fire damage at initiative 2\n989 units each with 1274 hit points (immune to fire; weak to bludgeoning, slashing) with an attack that does 25 slashing damage at initiative 3\n\nInfection:\n801 units each with 4706 hit points (weak to radiation) with an attack that does 116 bludgeoning damage at initiative 1\n4485 units each with 2961 hit points (immune to radiation; weak to fire, cold) with an attack that does 12 slashing damage at initiative 4"),
			Want:  "51",
		},
		{
			Input: aoc.FileInput(t, 2018, 24),
			Want:  "2834",
		},
	}
	aoc.AOCTest(t, day24p02, tests)
}
