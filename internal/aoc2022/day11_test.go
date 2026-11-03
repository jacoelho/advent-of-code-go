package aoc2022

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day11Example = `Monkey 0:
  Starting items: 79, 98
  Operation: new = old * 19
  Test: divisible by 23
    If true: throw to monkey 2
    If false: throw to monkey 3

Monkey 1:
  Starting items: 54, 65, 75, 74
  Operation: new = old + 6
  Test: divisible by 19
    If true: throw to monkey 2
    If false: throw to monkey 0

Monkey 2:
  Starting items: 79, 60, 97
  Operation: new = old * old
  Test: divisible by 13
    If true: throw to monkey 1
    If false: throw to monkey 3

Monkey 3:
  Starting items: 74
  Operation: new = old + 3
  Test: divisible by 17
    If true: throw to monkey 0
    If false: throw to monkey 1`

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day11Example), Want: "10605"},
		{Input: aoc.FileInput(t, 2022, 11), Want: "99840"},
	}
	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day11Example), Want: "2713310158"},
		{Input: aoc.FileInput(t, 2022, 11), Want: "20683044837"},
	}
	aoc.AOCTest(t, day11p02, tests)
}

const reversedMonkeys = `Monkey 1:
  Starting items:
  Operation: new = old + 0
  Test: divisible by 1
    If true: throw to monkey 0
    If false: throw to monkey 0

Monkey 0:
  Starting items: 9
  Operation: new = old + 0
  Test: divisible by 1
    If true: throw to monkey 1
    If false: throw to monkey 1`

func TestMonkeyIDsAndArithmetic(t *testing.T) {
	aoc.AOCTest(t, day11p01, []aoc.TestInput{{Input: strings.NewReader(reversedMonkeys), Want: "400"}})
	large := strings.Replace(reversedMonkeys,
		"Starting items: 9\n  Operation: new = old + 0",
		"Starting items: 9223372036854775807\n  Operation: new = old * old", 1)
	aoc.AOCTest(t, day11p02, []aoc.TestInput{{Input: strings.NewReader(reversedMonkeys), Want: "100000000"}, {Input: strings.NewReader(large), Want: "100000000"}})
}

func TestMonkeyInvalidInput(t *testing.T) {
	for _, input := range []string{
		strings.Replace(reversedMonkeys, "throw to monkey 0", "throw to monkey 1", 1),
		strings.Replace(reversedMonkeys, "Monkey 0:", "Monkey 1:", 1),
		strings.Replace(reversedMonkeys, "divisible by 1", "divisible by 0", 1),
	} {
		if _, err := day11p01(strings.NewReader(input)); err == nil {
			t.Fatal("expected invalid monkey error")
		}
	}
}
