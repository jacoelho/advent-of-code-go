package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const medicineExampleRules = "H => HO\nH => OH\nO => HH\n"

func Test_day19p01(t *testing.T) {
	aoc.AOCTest(t, day19p01, []aoc.TestInput{
		{Input: strings.NewReader(medicineExampleRules + "\nHOH"), Want: "4"},
		{Input: strings.NewReader(medicineExampleRules + "\nHOHOHO"), Want: "7"},
		{Input: strings.NewReader("AA => B\n\nAAA"), Want: "2"},
		{Input: strings.NewReader("H => OO\n\nH2O"), Want: "1"},
		{Input: strings.NewReader("H => X\n\nHe"), Want: "1"},
		{Input: aoc.FileInput(t, 2015, 19), Want: "576"},
	})
}

func Test_day19p02(t *testing.T) {
	rules := "e => H\ne => O\n" + medicineExampleRules
	aoc.AOCTest(t, day19p02, []aoc.TestInput{
		{Input: strings.NewReader(rules + "\nHOH"), Want: "3"},
		{Input: strings.NewReader(rules + "\nHOHOHO"), Want: "6"},
		{Input: strings.NewReader(rules + "\ne"), Want: "0"},
		{Input: aoc.FileInput(t, 2015, 19), Want: "207"},
	})
}

func TestMedicineMinimumDerivation(t *testing.T) {
	aoc.AOCTest(t, day19p02, []aoc.TestInput{
		{Input: strings.NewReader("e => A\ne => BC\nA => D\nD => BC\n\nBC"), Want: "1"},
		{Input: strings.NewReader("e => A\nA => B\nB => A\nB => CC\n\nCC"), Want: "3"},
		{Input: strings.NewReader("e => ABCD\n\nABCD"), Want: "1"},
		{Input: strings.NewReader("e => He\nH => X\n\nXe"), Want: "2"},
		{Input: strings.NewReader("e => He\nH => X\nXe => BC\n\nBC"), Want: "3"},
		{Input: strings.NewReader("e => Ar\nAr => B\n\nB"), Want: "2"},
	})
	for _, input := range []string{
		"e => A\nA => A\n\nB",
		"e => He\nH => H\n\nXe",
		"e => AB\nAB => A\n\nA",
	} {
		if _, err := day19p02(strings.NewReader(input)); err == nil {
			t.Errorf("accepted unreachable or unsupported grammar %q", input)
		}
	}
}
