package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1`),
			Want: "21",
		},
		{
			Input: aoc.FileInput(t, 2023, 12),
			Want:  "7694",
		},
	}
	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1`),
			Want: "525152",
		},
		{
			Input: aoc.FileInput(t, 2023, 12),
			Want:  "5071883216318",
		},
	}
	aoc.AOCTest(t, day12p02, tests)
}

func TestDay12RejectsMalformedRecords(t *testing.T) {
	requireRejected(t, []string{
		"???.###",
		"???.### 1,x",
		"???.### 1,0",
		"???.### 1,-1",
		"??x.### 1,1",
		" 1,1",
	}, day12p01, day12p02)
}
