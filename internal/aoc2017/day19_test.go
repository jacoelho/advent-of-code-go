package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

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

func TestDay19ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"19/1", day19p01},
		{"19/2", day19p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
