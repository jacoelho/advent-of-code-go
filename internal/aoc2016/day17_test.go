package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

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

func TestVaultBudgetAfterGoal(t *testing.T) {
	// The first DFS goal is reached after 78 generated paths; other branches remain.
	got, err := longestVault("ihgpwlah", vaultBudget{paths: 78, bytes: 1 << 20})
	if err == nil || got != 0 {
		t.Fatalf("longest returned %d, %v after exhaustion; want no partial answer", got, err)
	}
}

func TestDay17ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day17p01,
		day17p02,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error%v lost input cause", err)
		}
	}
}
