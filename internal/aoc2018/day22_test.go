package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestCaveRequiresTorchAtTarget(t *testing.T) {
	if _, err := day22p02(strings.NewReader("depth: 1\ntarget: 1,1")); err == nil {
		t.Fatal("expected invalid torch region error")
	}
}

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("depth: 510\ntarget: 10,10"),
			Want:  "114",
		},
		{
			Input: strings.NewReader("depth: 510\ntarget: 0,0"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 22),
			Want:  "5622",
		},
	}
	aoc.AOCTest(t, day22p01, tests)
}

func Test_day22p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("depth: 510\ntarget: 10,10"),
			Want:  "45",
		},
		{
			Input: strings.NewReader("depth: 510\ntarget: 0,0"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 22),
			Want:  "1089",
		},
	}
	aoc.AOCTest(t, day22p02, tests)
}
