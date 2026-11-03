package aoc2024

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47`),
			Want: "143",
		},
		{
			Input: aoc.FileInput(t, 2024, 5),
			Want:  "4996",
		},
	}

	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47`),
			Want: "123",
		},
		{
			Input: aoc.FileInput(t, 2024, 5),
			Want:  "6311",
		},
	}

	aoc.AOCTest(t, day05p02, tests)
}

func TestPageOrdering(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  [2]string
	}{
		{"chain", "1|2\n2|3\n\n3,1,2", [2]string{"0", "2"}},
		{"valid", "1|2\n2|3\n\n1,2,3", [2]string{"2", "0"}},
		{"original order breaks ties", "1|3\n\n3,2,1", [2]string{"0", "1"}},
		{"absent intermediate page", "1|2\n2|3\n\n3,1,4", [2]string{"1", "0"}},
		{"unrelated cycle", "8|9\n9|8\n1|2\n\n2,1,3", [2]string{"0", "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for part, solve := range []func(io.Reader) (string, error){day05p01, day05p02} {
				got, err := solve(strings.NewReader(tt.input))
				if err != nil || got != tt.want[part] {
					t.Errorf("part %d: got %q, %v; want %q", part+1, got, err, tt.want[part])
				}
			}
		})
	}
}

func TestPageOrderingCycle(t *testing.T) {
	input := "1|2\n2|1\n\n1,2,3"
	if got, err := day05p01(strings.NewReader(input)); err != nil || got != "0" {
		t.Errorf("part 1: got %q, %v; want 0", got, err)
	}
	if got, err := day05p02(strings.NewReader(input)); err == nil {
		t.Errorf("part 2: got %q, want cycle error", got)
	}
}

func TestPageOrderingInvalidInput(t *testing.T) {
	for _, input := range []string{"1|2\n\n1,2,1", "1|2|3\n\n1,2,3"} {
		for _, solve := range []func(io.Reader) (string, error){day05p01, day05p02} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("input %q: got %q, want error", input, got)
			}
		}
	}
	for _, input := range []string{"x|2\n\n1,2,3", "1|2\n\n1,x,3"} {
		for _, solve := range []func(io.Reader) (string, error){day05p01, day05p02} {
			if _, err := solve(strings.NewReader(input)); !errors.Is(err, strconv.ErrSyntax) {
				t.Errorf("input %q: got %v, want numeric syntax error", input, err)
			}
		}
	}
}

func TestDay05ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day05p01, "1|2\n\n1,2,3\n"},
		{"part2", day05p02, "1|2\n\n1,2,3\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				r := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(want))
				if _, err := tt.solve(r); !errors.Is(err, want) {
					t.Errorf("prefix %q: got %v, want reader error", prefix, err)
				}
			}
		})
	}
}
