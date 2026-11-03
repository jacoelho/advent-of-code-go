package aoc2023

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`broadcaster -> a, b, c
%a -> b
%b -> c
%c -> inv
&inv -> a`),
			Want: "32000000",
		},
		{
			Input: strings.NewReader(`broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output`),
			Want: "11687500",
		},
		{
			Input: aoc.FileInput(t, 2023, 20),
			Want:  "787056720",
		},
	}
	aoc.AOCTest(t, day20p01, tests)
}

func Test_day20p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> rx`),
			Want: "1",
		},
		{
			Input: aoc.FileInput(t, 2023, 20),
			Want:  "212986464842911",
		},
	}
	aoc.AOCTest(t, day20p02, tests)
}

const pulseRememberedHigh = `broadcaster -> a, b
%a -> con
%b -> c
%c -> con
&con -> rx`

const pulseBeyondPrefix = `broadcaster -> a, c
%a -> b
%b -> a, inv
&inv -> g
%c -> d
%d -> g
&g -> rx`

func TestPulseNetworks(t *testing.T) {
	for _, tt := range []struct{ name, input, first, second string }{
		{"remembered high", pulseRememberedHigh, "11250000", "3"},
		{"beyond observed prefix", pulseBeyondPrefix, "18878112", "7"},
		{"duplicate root first", "broadcaster -> a, a, b\n%a -> g\n%b -> g\n&g -> rx", "24000000", "2"},
		{"duplicate root last", "broadcaster -> b, a, a\n%a -> g\n%b -> g\n&g -> rx", "24000000", "1"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for i, solve := range []func(io.Reader) (string, error){day20p01, day20p02} {
				want := [2]string{tt.first, tt.second}[i]
				got, err := solve(strings.NewReader(tt.input))
				if err != nil || got != want {
					t.Fatalf("part %d = %q, %v; want %s", i+1, got, err, want)
				}
			}
		})
	}
	for _, tt := range []struct{ input, want string }{
		{"broadcaster -> rx", "1"},
		{"broadcaster -> a, b\n%a -> rx\n%b -> rx", "2"},
	} {
		got, err := day20p02(strings.NewReader(tt.input))
		if err != nil || got != tt.want {
			t.Fatalf("input %q = %q, %v; want %s", tt.input, got, err, tt.want)
		}
	}
}

func TestPulseNoSolution(t *testing.T) {
	for _, input := range []string{
		"broadcaster -> a\n%a -> sink",
		"broadcaster -> g\n&g -> rx",
	} {
		_, err := day20p02(strings.NewReader(input))
		if err == nil || !strings.Contains(err.Error(), "rx") {
			t.Fatalf("input %q: got %v; want no rx solution", input, err)
		}
	}
}

func TestPulseInputErrors(t *testing.T) {
	for _, input := range []string{
		"%a -> rx",
		"broadcaster -> a\n%a -> rx\n&a -> rx",
		"broadcaster -> a\n!a -> rx",
		"broadcaster -> ",
		"broadcaster => rx",
	} {
		for _, solve := range []func(io.Reader) (string, error){day20p01, day20p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected error", input)
			}
		}
	}
}

func TestPulseNonDrainingPress(t *testing.T) {
	network, err := parsePulseNetwork(strings.NewReader("broadcaster -> a\n&a -> a, rx"))
	if err != nil {
		t.Fatal(err)
	}
	remaining := 32
	err = network.press(&remaining, func(pulseMessage) error { return nil })
	if err == nil || !strings.Contains(err.Error(), "event limit") {
		t.Fatalf("got %v; want pulse event limit", err)
	}
}

func TestPulseEqualDepthOrdering(t *testing.T) {
	a := pulseBranch{period: 1, traces: []pulseTrace{
		{events: []gatePulse{{depth: 2, high: true}, {depth: 2, high: false}}},
	}}
	b := pulseBranch{period: 2, traces: []pulseTrace{
		{events: []gatePulse{{depth: 2, high: true}}},
		{startHigh: true, events: []gatePulse{{depth: 2, high: false}}},
	}}
	for _, tt := range []struct {
		branches []pulseBranch
		want     string
	}{
		{[]pulseBranch{a, b}, "4"},
		{[]pulseBranch{b, a}, "3"},
	} {
		got, err := alignPulseBranches(tt.branches, 2)
		if err != nil || got.String() != tt.want {
			t.Fatalf("next low = %v, %v; want %s", got, err, tt.want)
		}
	}
}
