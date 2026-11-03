package aoc2022

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day16p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(valveExample), Want: "1651"},
		{Input: strings.NewReader("Valve AA has flow rate=0; tunnel leads to valve BB\nValve BB has flow rate=5; tunnel leads to valve AA"), Want: "140"},
		{Input: strings.NewReader("Valve AA has flow rate=0; tunnel leads to valve AA"), Want: "0"},
		{Input: strings.NewReader("Valve AA has flow rate=5; tunnel leads to valve AA"), Want: "145"},
		{Input: aoc.FileInput(t, 2022, 16), Want: "1751"},
	}

	aoc.AOCTest(t, day16p01, tests)
}

func Test_day16p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(valveExample), Want: "1707"},
		{Input: strings.NewReader("Valve AA has flow rate=0; tunnel leads to valve BB\nValve BB has flow rate=5; tunnel leads to valve AA"), Want: "120"},
		{Input: strings.NewReader("Valve AA has flow rate=0; tunnel leads to valve AA"), Want: "0"},
		{Input: strings.NewReader("Valve AA has flow rate=5; tunnel leads to valve AA"), Want: "125"},
		{Input: aoc.FileInput(t, 2022, 16), Want: "2207"},
	}

	aoc.AOCTest(t, day16p02, tests)
}

func TestValveErrors(t *testing.T) {
	for _, input := range []string{
		"Valve BB has flow rate=5; tunnel leads to valve BB",
		"Valve AA has flow rate=0; tunnel leads to valve BB",
		"Valve AA has flow rate=-1; tunnel leads to valve AA",
		"Valve AA has flow rate=0; tunnel leads to valve AA\nValve AA has flow rate=0; tunnel leads to valve AA",
	} {
		if _, err := day16p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

const valveExample = `Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II`
