package aoc2021

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const simultaneousBingo = `1

1 1 1 1 1
2 2 2 2 2
2 2 2 2 2
2 2 2 2 2
2 2 2 2 2

1 1 1 1 1
3 3 3 3 3
3 3 3 3 3
3 3 3 3 3
3 3 3 3 3

9 9 9 9 9
9 9 9 9 9
9 9 9 9 9
9 9 9 9 9
9 9 9 9 9`

const day04Example = `7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
 8  2 23  4 24
21  9 14 16  7
 6 10  3 18  5
 1 12 20 15 19

 3 15  0  2 22
 9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
 2  0 12  3  7`

func Test_day04p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day04Example), Want: "4512"},
		{Input: strings.NewReader(simultaneousBingo), Want: "40"},
		{Input: aoc.FileInput(t, 2021, 4), Want: "21607"},
	}

	aoc.AOCTest(t, day04p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day04Example), Want: "1924"},
		{Input: strings.NewReader(simultaneousBingo), Want: "60"},
		{Input: aoc.FileInput(t, 2021, 4), Want: "19012"},
	}

	aoc.AOCTest(t, day04p02, tests)
}

func TestDay04IncompleteBoard(t *testing.T) {
	if _, err := day04p01(strings.NewReader(`1

1 2 3 4 5`)); err == nil {
		t.Fatal("expected incomplete board error")
	}
}
