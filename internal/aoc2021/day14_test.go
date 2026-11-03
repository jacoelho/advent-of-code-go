package aoc2021

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day14Example = `NNCB

CH -> B
HH -> N
CB -> H
NH -> C
HB -> C
HC -> B
HN -> C
NN -> C
BH -> H
NC -> B
NB -> B
BN -> B
BB -> N
BC -> B
CC -> N
CN -> C`

func Test_day14p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day14Example), Want: "1588"},
		{Input: strings.NewReader(`AAB

AB -> C`), Want: "1"},
		{Input: strings.NewReader(`A`), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 14), Want: "3230"},
	}

	aoc.AOCTest(t, day14p01, tests)
}

func Test_day14p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day14Example), Want: "2188189693529"},
		{Input: strings.NewReader(`AAB

AB -> C`), Want: "1"},
		{Input: strings.NewReader(`A`), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 14), Want: "3542388214529"},
	}

	aoc.AOCTest(t, day14p02, tests)
}

func TestDay14InvalidRule(t *testing.T) {
	if _, err := day14p01(strings.NewReader(`AB

AB -> CD`)); err == nil {
		t.Fatal("expected invalid rule error")
	}
}
