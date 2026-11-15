package aoc2019

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("COM)B\nB)C\nC)D\nD)E\nE)F\nB)G\nG)H\nD)I\nE)J\nJ)K\nK)L"),
			Want:  "42",
		},
		{
			Input: aoc.FileInput(t, 2019, 6),
			Want:  "204521",
		},
	}
	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("COM)B\nB)C\nC)D\nD)E\nE)F\nB)G\nG)H\nD)I\nE)J\nJ)K\nK)L\nK)YOU\nI)SAN"),
			Want:  "4",
		},
		{
			Input: aoc.FileInput(t, 2019, 6),
			Want:  "307",
		},
	}
	aoc.AOCTest(t, day06p02, tests)
}

func TestOrbitMapErrors(t *testing.T) {
	for _, input := range []string{"", "COM", "COM)", ")B", "A)B)C", "COM)B\nA)B", "A)B\nB)A"} {
		for _, solve := range []func(io.Reader) (string, error){day06p01, day06p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("expected orbit map error for %q", input)
			}
		}
	}
	for _, input := range []string{"COM)YOU", "COM)YOU\nX)SAN"} {
		if _, err := day06p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected transfer error for %q", input)
		}
	}
}
