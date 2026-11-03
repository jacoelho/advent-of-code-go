package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	aoc.AOCTest(t, day03p01, []aoc.TestInput{
		{Input: strings.NewReader("5 10 25"), Want: "0"},
		{Input: strings.NewReader("1 1 2\n3 4 5"), Want: "1"},
		{Input: aoc.FileInput(t, 2016, 3), Want: "1050"},
	})
}

func Test_day03p02(t *testing.T) {
	aoc.AOCTest(t, day03p02, []aoc.TestInput{
		{Input: strings.NewReader("101 301 501\n102 302 502\n103 303 503\n201 401 601\n202 402 602\n203 403 603"), Want: "6"},
		{Input: aoc.FileInput(t, 2016, 3), Want: "1921"},
	})
}
