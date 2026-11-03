package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	aoc.AOCTest(t, day07p01, []aoc.TestInput{
		{Input: strings.NewReader("abba[mnop]qrst\nabcd[bddb]xyyx\naaaa[qwer]tyui\nioxxoj[asdfgh]zxcvbn"), Want: "2"},
		{Input: aoc.FileInput(t, 2016, 7), Want: "115"},
	})
}

func Test_day07p02(t *testing.T) {
	aoc.AOCTest(t, day07p02, []aoc.TestInput{
		{Input: strings.NewReader("aba[bab]xyz\nxyx[xyx]xyx\naaa[kek]eke\nzazbz[bzb]cdb"), Want: "3"},
		{Input: aoc.FileInput(t, 2016, 7), Want: "231"},
	})
}
