package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day23p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("pos=<0,0,0>, r=4\npos=<1,0,0>, r=1\npos=<4,0,0>, r=3\npos=<0,2,0>, r=1\npos=<0,5,0>, r=3\npos=<0,0,3>, r=1\npos=<1,1,1>, r=1\npos=<1,1,2>, r=1\npos=<1,3,1>, r=1"),
			Want:  "7",
		},
		{
			Input: aoc.FileInput(t, 2018, 23),
			Want:  "491",
		},
	}
	aoc.AOCTest(t, day23p01, tests)
}

func Test_day23p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("pos=<10,12,12>, r=2\npos=<12,14,12>, r=2\npos=<16,12,12>, r=4\npos=<14,14,14>, r=6\npos=<50,50,50>, r=200\npos=<10,10,10>, r=5"),
			Want:  "36",
		},
		{
			Input: strings.NewReader("pos=<10,0,0>, r=0\npos=<-10,0,0>, r=0\npos=<1,0,0>, r=0"),
			Want:  "1",
		},
		{
			Input: strings.NewReader("pos=<10,0,0>, r=0\npos=<10,0,0>, r=0\npos=<0,0,0>, r=0"),
			Want:  "10",
		},
		{
			Input: aoc.FileInput(t, 2018, 23),
			Want:  "60474080",
		},
	}
	aoc.AOCTest(t, day23p02, tests)
}
