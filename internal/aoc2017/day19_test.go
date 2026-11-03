package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const tubeExample = `     |
     |  +--+
     A  |  C
 F---|----E|--+
     |  |  |  D
     +B-+  +--+`

func Test_day19p01(t *testing.T) {
	aoc.AOCTest(t, day19p01, []aoc.TestInput{
		{Input: strings.NewReader(tubeExample), Want: "ABCDEF"},
		{Input: aoc.FileInput(t, 2017, 19), Want: "LOHMDQATP"},
	})
}

func Test_day19p02(t *testing.T) {
	aoc.AOCTest(t, day19p02, []aoc.TestInput{
		{Input: strings.NewReader(tubeExample), Want: "38"},
		{Input: strings.NewReader(" |\n A"), Want: "2"},
		{Input: aoc.FileInput(t, 2017, 19), Want: "16492"},
	})
}

func TestFollowTubesRejectsAmbiguousTurn(t *testing.T) {
	if _, _, err := followTubes(strings.NewReader(" |\n-+-")); err == nil {
		t.Fatal("expected ambiguous turn error")
	}
}
