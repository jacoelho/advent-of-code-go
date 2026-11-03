package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`abcdef
bababc
abbcde
abcccd
aabcdd
abcdee
ababab`), Want: "12"},
		{Input: aoc.FileInput(t, 2018, 2), Want: "6972"},
	}

	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`abcde
fghij
klmno
pqrst
fguij
axcye
wvxyz`), Want: "fgij"},
		{Input: aoc.FileInput(t, 2018, 2), Want: "aixwcbzrmdvpsjfgllthdyoqe"},
	}

	aoc.AOCTest(t, day02p02, tests)
}
