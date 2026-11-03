package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`light red bags contain 1 bright white bag, 2 muted yellow bags.
dark orange bags contain 3 bright white bags, 4 muted yellow bags.
bright white bags contain 1 shiny gold bag.
muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.
shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
dark olive bags contain 3 faded blue bags, 4 dotted black bags.
vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.
faded blue bags contain no other bags.
dotted black bags contain no other bags.`), Want: "4"},
		{Input: strings.NewReader(`shiny gold bags contain 2 dark red bags.
dark red bags contain 2 dark orange bags.
dark orange bags contain 2 dark yellow bags.
dark yellow bags contain 2 dark green bags.
dark green bags contain 2 dark blue bags.
dark blue bags contain 2 dark violet bags.
dark violet bags contain no other bags.`), Want: "0"},
		{Input: aoc.FileInput(t, 2020, 7), Want: "246"},
	}
	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`light red bags contain 1 bright white bag, 2 muted yellow bags.
dark orange bags contain 3 bright white bags, 4 muted yellow bags.
bright white bags contain 1 shiny gold bag.
muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.
shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
dark olive bags contain 3 faded blue bags, 4 dotted black bags.
vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.
faded blue bags contain no other bags.
dotted black bags contain no other bags.`), Want: "32"},
		{Input: strings.NewReader(`shiny gold bags contain 2 dark red bags.
dark red bags contain 2 dark orange bags.
dark orange bags contain 2 dark yellow bags.
dark yellow bags contain 2 dark green bags.
dark green bags contain 2 dark blue bags.
dark blue bags contain 2 dark violet bags.
dark violet bags contain no other bags.`), Want: "126"},
		{Input: aoc.FileInput(t, 2020, 7), Want: "2976"},
	}
	aoc.AOCTest(t, day07p02, tests)
}
func TestBagRulesCycle(t *testing.T) {
	if _, err := day07p02(strings.NewReader("shiny gold bags contain 1 dark red bag.\ndark red bags contain 1 shiny gold bag.")); err == nil {
		t.Fatal("expected cycle error")
	}
}

func TestDays07To12ReadErrors(t *testing.T) {
	failure := errors.New("read failure")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"07/1", day07p01}, {"07/2", day07p02}, {"08/1", day08p01}, {"08/2", day08p02},
		{"09/1", day09p01}, {"09/2", day09p02}, {"10/1", day10p01}, {"10/2", day10p02},
		{"11/1", day11p01}, {"11/2", day11p02}, {"12/1", day12p01}, {"12/2", day12p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestDays07To12MalformedInputs(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"bags", "shiny gold bags contain many dark red bags.", day07p01},
		{"boot", "mul +1", day08p01},
		{"XMAS", "not a number", day09p01},
		{"adapters", "-1", day10p01},
		{"seats", "LL\nL", day11p01},
		{"navigation", "R45", day12p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected malformed-input error")
			}
		})
	}
}
