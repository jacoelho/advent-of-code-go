package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day18p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("set a 1\nadd a 2\nmul a a\nmod a 5\nsnd a\nset a 0\nrcv a\njgz a -1\nset a 1\njgz a -2"),
			Want:  "4",
		},
		{
			Input: strings.NewReader("snd +7\nrcv 1"),
			Want:  "7",
		},
		{
			Input: strings.NewReader("snd 0\nrcv 1"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("snd 5\nrcv 0\nsnd 7\nrcv 1"),
			Want:  "7",
		},
		{
			Input: strings.NewReader("snd 5\njgz -1 2\nsnd 7\nrcv 1"),
			Want:  "7",
		},
		{
			Input: aoc.FileInput(t, 2017, 18),
			Want:  "3188",
		},
	}

	aoc.AOCTest(t, day18p01, tests)
}

func Test_day18p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("snd 1\nsnd 2\nsnd p\nrcv a\nrcv b\nrcv c\nrcv d"),
			Want:  "3",
		},
		{
			Input: strings.NewReader("rcv a"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("jgz p 5\nrcv a\nsnd a\nrcv b\njgz 1 5\nsnd 7\nrcv c\njgz c 2\nsnd 9"),
			Want:  "1",
		},
		{
			Input: aoc.FileInput(t, 2017, 18),
			Want:  "7112",
		},
	}

	aoc.AOCTest(t, day18p02, tests)
}

func TestDuetErrors(t *testing.T) {
	inputs := []string{
		"",
		"set 1 2",
		"set aa 2",
		"set a -",
		"snd 1 2",
		"jnz 1 2",
		"rcv 1",
		"set a 9223372036854775807\nadd a 1",
		"set a -9223372036854775808\nmul a -1",
		"mod a 0",
	}
	for _, input := range inputs {
		if _, err := day18p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}

	if _, err := day18p02(strings.NewReader("rcv 1")); err == nil {
		t.Error("expected receive-destination error")
	}
}
