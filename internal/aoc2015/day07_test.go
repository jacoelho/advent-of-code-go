package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestCircuitGates(t *testing.T) {
	input := `123 -> x
456 -> y
x AND y -> d
x OR y -> e
x LSHIFT 2 -> f
y RSHIFT 2 -> g
NOT x -> h
NOT y -> i
32768 LSHIFT 1 -> j`
	for wire, want := range map[string]string{
		"d": "72", "e": "507", "f": "492", "g": "114",
		"h": "65412", "i": "65079", "j": "0",
	} {
		t.Run(wire, func(t *testing.T) {
			aoc.AOCTest(t, day07p01, []aoc.TestInput{{Input: strings.NewReader(input + "\n" + wire + " -> a"), Want: want}})
		})
	}
}

func TestCircuitUnresolvableSignal(t *testing.T) {
	for _, input := range []string{"b -> a\na -> b", "missing -> a"} {
		if _, err := day07p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected error for circuit %q", input)
		}
	}
}

func Test_day07p01(t *testing.T) {
	aoc.AOCTest(t, day07p01, []aoc.TestInput{
		{Input: strings.NewReader("123 -> x\n456 -> y\nx AND y -> d\nx OR y -> e\nx LSHIFT 2 -> f\ny RSHIFT 2 -> g\nNOT x -> h\nNOT y -> i\nd -> a"), Want: "72"},
		{Input: aoc.FileInput(t, 2015, 7), Want: "3176"},
	})
}

func Test_day07p02(t *testing.T) {
	aoc.AOCTest(t, day07p02, []aoc.TestInput{
		{Input: strings.NewReader("1 -> b\nb LSHIFT 1 -> a"), Want: "4"},
		{Input: aoc.FileInput(t, 2015, 7), Want: "14710"},
	})
}
