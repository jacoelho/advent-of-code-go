package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const smallCoprocessor = `set b 3
set c b
jnz a 2
jnz 1 5
mul b 1
sub b 0
set c b
sub c -4
set f 1
set d 2
set e 2
set g d
mul g e
sub g b
jnz g 2
set f 0
sub e -1
set g e
sub g b
jnz g -8
sub d -1
set g d
sub g b
jnz g -13
jnz f 2
sub h -1
set g b
sub g c
jnz g 2
jnz 1 3
sub b -1
jnz 1 -23`

func Test_day23p01(t *testing.T) {
	aoc.AOCTest(t, day23p01, []aoc.TestInput{
		{Input: strings.NewReader("set a 2\nmul a a\nsub a 1\njnz a 2\nmul a a"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 23), Want: "4225"},
	})
}

func Test_day23p02(t *testing.T) {
	aoc.AOCTest(t, day23p02, []aoc.TestInput{
		{Input: strings.NewReader(smallCoprocessor), Want: "2"},
		{Input: strings.NewReader(strings.Replace(smallCoprocessor, "sub h -1", "sub h -2", 1)), Want: "4"},
		{Input: strings.NewReader("set h 7"), Want: "7"},
		{Input: aoc.FileInput(t, 2017, 23), Want: "905"},
	})
}

func TestCoprocessorRejectsOverflow(t *testing.T) {
	for _, input := range []string{
		"set h 9223372036854775807\nsub h -1",
		"set h -9223372036854775808\nsub h 1",
		"set h 9223372036854775807\nmul h 2",
	} {
		if _, err := day23p02(strings.NewReader(input)); err == nil {
			t.Fatalf("expected overflow for %q", input)
		}
	}
}

func TestCoprocessorInstructionBudget(t *testing.T) {
	program, err := parseCoprocessor(strings.NewReader("jnz 1 0"))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := runCoprocessor(program, [26]int64{}, 10); err == nil {
		t.Fatal("expected instruction limit error")
	}
}
