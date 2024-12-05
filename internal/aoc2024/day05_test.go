package aoc2024

import (
	"strings"
	"testing"

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
