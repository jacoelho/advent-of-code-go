package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	aoc.AOCTest(t, day06p01, []aoc.TestInput{
		{Input: strings.NewReader("eedadn\ndrvtee\neandsr\nraavrd\natevrs\ntsrnev\nsdttsa\nrasrtv\nnssdts\nntnada\nsvetve\ntesnvt\nvntsnd\nvrdear\ndvrsen\nenarar"), Want: "easter"},
		{Input: aoc.FileInput(t, 2016, 6), Want: "qqqluigu"},
	})
}

func Test_day06p02(t *testing.T) {
	aoc.AOCTest(t, day06p02, []aoc.TestInput{
		{Input: strings.NewReader("eedadn\ndrvtee\neandsr\nraavrd\natevrs\ntsrnev\nsdttsa\nrasrtv\nnssdts\nntnada\nsvetve\ntesnvt\nvntsnd\nvrdear\ndvrsen\nenarar"), Want: "advent"},
		{Input: aoc.FileInput(t, 2016, 6), Want: "lsoypmia"},
	})
}
