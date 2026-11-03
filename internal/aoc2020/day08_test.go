package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day08p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`nop +0
acc +1
jmp +4
acc +3
jmp -3
acc -99
acc +1
jmp -4
acc +6`), Want: "5"},
		{Input: aoc.FileInput(t, 2020, 8), Want: "1930"},
	}
	aoc.AOCTest(t, day08p01, tests)
}

func Test_day08p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`nop +0
acc +1
jmp +4
acc +3
jmp -3
acc -99
acc +1
jmp -4
acc +6`), Want: "8"},
		{Input: aoc.FileInput(t, 2020, 8), Want: "1688"},
	}
	aoc.AOCTest(t, day08p02, tests)
}
func TestBootCodeTermination(t *testing.T) {
	aoc.AOCTest(t, day08p01, []aoc.TestInput{{Input: strings.NewReader("acc +7"), Want: "7"}})
	for _, input := range []string{"jmp +2", "jmp -1"} {
		if _, err := day08p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected invalid jump error for %q", input)
		}
	}
	aoc.AOCTest(t, day08p02, []aoc.TestInput{{Input: strings.NewReader("nop +99\njmp -1"), Want: "0"}})
	if _, err := day08p02(strings.NewReader("acc +7")); err == nil {
		t.Fatal("cannot repair by changing acc")
	}
}
