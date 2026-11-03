package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const skyExample = `position=< 9,  1> velocity=< 0,  2>
position=< 7,  0> velocity=<-1,  0>
position=< 3, -2> velocity=<-1,  1>
position=< 6, 10> velocity=<-2, -1>
position=< 2, -4> velocity=< 2,  2>
position=<-6, 10> velocity=< 2, -2>
position=< 1,  8> velocity=< 1, -1>
position=< 1,  7> velocity=< 1,  0>
position=<-3, 11> velocity=< 1, -2>
position=< 7,  6> velocity=<-1, -1>
position=<-2,  3> velocity=< 1,  0>
position=<-4,  3> velocity=< 2,  0>
position=<10, -3> velocity=<-1,  1>
position=< 5, 11> velocity=< 1, -2>
position=< 4,  7> velocity=< 0, -1>
position=< 8, -2> velocity=< 0,  1>
position=<15,  0> velocity=<-2,  0>
position=< 1,  6> velocity=< 1,  0>
position=< 8,  9> velocity=< 0, -1>
position=< 3,  3> velocity=<-1,  1>
position=< 0,  5> velocity=< 0, -1>
position=<-2,  2> velocity=< 2,  0>
position=< 5, -2> velocity=< 1,  2>
position=< 1,  4> velocity=< 2,  1>
position=<-2,  7> velocity=< 2, -2>
position=< 3,  6> velocity=<-1, -1>
position=< 5,  0> velocity=< 1,  0>
position=<-6,  0> velocity=< 2,  0>
position=< 5,  9> velocity=< 1, -2>
position=<14,  7> velocity=<-2,  0>
position=<-3,  6> velocity=< 2, -1>`

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(skyExample),
			Want:  "HI",
		},
		{
			Input: aoc.FileInput(t, 2018, 10),
			Want:  "JJXZHKFP",
		},
	}

	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(skyExample),
			Want:  "3",
		},
		{
			Input: strings.NewReader("position=<0,0> velocity=<0,0>\nposition=<-1,-20> velocity=<1,2>"),
			Want:  "10",
		},
		{
			Input: strings.NewReader("position=<0,0> velocity=<1,1>\nposition=<5,5> velocity=<1,1>"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 10),
			Want:  "10036",
		},
	}

	aoc.AOCTest(t, day10p02, tests)
}

func TestSkyErrors(t *testing.T) {
	for _, input := range []string{"", "position 1,2", "position=<0,0> velocity=<0,0>\nposition=<3000000,0> velocity=<-1,0>"} {
		if _, err := day10p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}
