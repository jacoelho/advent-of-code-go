package aoc2019

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("R1,L1,U2,R2\nL1,R1,D2,R2,U4"), Want: "4"},
		{
			Input: strings.NewReader("R8,U5,L5,D3\nU7,R6,D4,L4"),
			Want:  "6",
		},
		{
			Input: aoc.FileInput(t, 2019, 3),
			Want:  "386",
		},
	}
	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("R1,L1,U2,R2\nL1,R1,D2,R2,U4"), Want: "16"},
		{
			Input: strings.NewReader("R8,U5,L5,D3\nU7,R6,D4,L4"),
			Want:  "30",
		},
		{
			Input: aoc.FileInput(t, 2019, 3),
			Want:  "6484",
		},
	}
	aoc.AOCTest(t, day03p02, tests)
}

func TestWiresWithoutIntersection(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day03p01, day03p02} {
		if _, err := solve(strings.NewReader("R1\nL1")); err == nil {
			t.Fatal("expected no-intersection error")
		}
	}
}
