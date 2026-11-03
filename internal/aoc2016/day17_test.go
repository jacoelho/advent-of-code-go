package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	aoc.AOCTest(t, day17p01, []aoc.TestInput{
		{Input: strings.NewReader("ihgpwlah"), Want: "DDRRRD"},
		{Input: aoc.FileInput(t, 2016, 17), Want: "RLDUDRDDRR"},
	})
}

func Test_day17p02(t *testing.T) {
	aoc.AOCTest(t, day17p02, []aoc.TestInput{
		{Input: strings.NewReader("ihgpwlah"), Want: "370"},
		{Input: aoc.FileInput(t, 2016, 17), Want: "590"},
	})
}

func TestUnreachableVault(t *testing.T) {
	if _, err := day17p01(strings.NewReader("hijkl")); err == nil {
		t.Fatal("accepted unreachable vault")
	}
	if _, err := day17p02(strings.NewReader("hijkl")); err == nil {
		t.Fatal("accepted unreachable vault")
	}
}

func TestVaultSearchBudget(t *testing.T) {
	for _, budget := range []vaultBudget{{1, 1 << 20}, {1000, 1}} {
		if got, err := shortestVault("ihgpwlah", budget); err == nil || got != "" {
			t.Errorf("shortest returned%q,%v after exhaustion", got, err)
		}
		if got, err := longestVault("ihgpwlah", budget); err == nil || got != 0 {
			t.Errorf("longest returned%d,%v after exhaustion", got, err)
		}
	}
}
