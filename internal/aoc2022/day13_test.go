package aoc2022

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day13Example = `[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]`

func Test_day13p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day13Example), Want: "13"},
		{Input: aoc.FileInput(t, 2022, 13), Want: "5393"},
	}
	aoc.AOCTest(t, day13p01, tests)
}

func Test_day13p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day13Example), Want: "140"},
		{Input: aoc.FileInput(t, 2022, 13), Want: "26712"},
	}
	aoc.AOCTest(t, day13p02, tests)
}

func TestPacketEqualityAndDividers(t *testing.T) {
	aoc.AOCTest(t, day13p01, []aoc.TestInput{{Input: strings.NewReader("[]\n[]"), Want: "0"}, {Input: strings.NewReader("[[2]]\n[[6]]"), Want: "1"}})
	aoc.AOCTest(t, day13p02, []aoc.TestInput{{Input: strings.NewReader("[]\n[]"), Want: "12"}, {Input: strings.NewReader("[[2]]\n[[6]]"), Want: "8"}})
	if _, err := day13p01(strings.NewReader("[1,]\n[]")); err == nil {
		t.Fatal("expected malformed packet error")
	}
}
