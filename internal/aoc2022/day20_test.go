package aoc2022

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(mixingExample), Want: "3"},
		{Input: strings.NewReader("0\n1\n1\n2\n-2\n3\n-3"), Want: "4"},
		{Input: aoc.FileInput(t, 2022, 20), Want: "23321"},
	}

	aoc.AOCTest(t, day20p01, tests)
}

func Test_day20p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(mixingExample), Want: "1623178306"},
		{Input: strings.NewReader("0\n1\n1\n2\n-2\n3\n-3"), Want: "-3246356612"},
		{Input: aoc.FileInput(t, 2022, 20), Want: "1428396909280"},
	}

	aoc.AOCTest(t, day20p02, tests)
}

func TestMixingErrors(t *testing.T) {
	for _, input := range []string{"0", "1\n2", "0\n0", "0\n9223372036854775807", "0\nx"} {
		if _, err := day20p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected mixing input error for %q", input)
		}
	}
}

const mixingExample = `1
2
-3
3
-2
0
4`
