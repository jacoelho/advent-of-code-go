package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`pbga (66)
xhth (57)
ebii (61)
havc (66)
ktlj (57)
fwft (72) -> ktlj, cntj, xhth
qoyq (66)
padx (45) -> pbga, havc, qoyq
tknk (41) -> ugml, padx, fwft
jptl (61)
ugml (68) -> gyxo, ebii, jptl
gyxo (61)
cntj (57)`), Want: "tknk"},
		{Input: aoc.FileInput(t, 2017, 7), Want: "vtzay"},
	}

	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`pbga (66)
xhth (57)
ebii (61)
havc (66)
ktlj (57)
fwft (72) -> ktlj, cntj, xhth
qoyq (66)
padx (45) -> pbga, havc, qoyq
tknk (41) -> ugml, padx, fwft
jptl (61)
ugml (68) -> gyxo, ebii, jptl
gyxo (61)
cntj (57)`), Want: "60"},
		{Input: strings.NewReader(`root (0) -> branch, left, right
left (9)
right (9)
branch (0) -> a, b, c
a (3)
b (3)
c (2)`), Want: "3"},
		{Input: strings.NewReader(`root (0) -> branch, left, right
left (10)
right (10)
branch (0) -> a, b
a (4)
b (5)`), Want: "5"},
		{Input: strings.NewReader(`root (0) -> holder
holder (0) -> a, b, c
a (5)
b (5)
c (4)`), Want: "5"},
		{Input: strings.NewReader(`r (0) -> a, b
a (0) -> x, y, z
b (15)
x (5)
y (5)
z (7)`), Want: "5"},
		{Input: strings.NewReader(`r (0) -> a, b
a (0) -> x, y
b (10)
x (5)
y (7)`), Want: "5"},
		{Input: strings.NewReader(`r (0) -> a, b
a (0) -> x, y
b (5)
x (5)
y (5)`), Want: "10"},
		{Input: strings.NewReader(`r (0) -> a, b, c
a (10)
b (10)
c (0) -> leaf
leaf (12)`), Want: "10"},
		{Input: aoc.FileInput(t, 2017, 7), Want: "910"},
	}

	aoc.AOCTest(t, day07p02, tests)
}

func Test_day07p02InvalidCorrection(t *testing.T) {
	tests := map[string]string{
		"ambiguous children": `r (0) -> a, b
a (5)
b (7)`,
		"inconsistent ancestor": `r (0) -> a, b
a (0) -> x, y
b (12)
x (5)
y (7)`,
		"unbalanced sibling subtree": `r (0) -> a, b, c
a (11)
b (0) -> x, y
c (10)
x (4)
y (6)`,
		"negative corrected weight": `r (0) -> a, b, c
a (0) -> x, y
b (1)
c (1)
x (1)
y (1)`,
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := day07p02(strings.NewReader(input)); err == nil {
				t.Fatal("expected an invalid correction error")
			}
		})
	}
}
