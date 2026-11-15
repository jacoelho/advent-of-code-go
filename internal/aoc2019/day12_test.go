package aoc2019

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const moonExample = `<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>`

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 12),
			Want:  "8362",
		},
	}
	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(moonExample), Want: "4686774924"},
		{
			Input: aoc.FileInput(t, 2019, 12),
			Want:  "478373365921244",
		},
	}
	aoc.AOCTest(t, day12p02, tests)
}

func TestMoonErrors(t *testing.T) {
	for _, input := range []string{
		"",
		"<x=1, y=2, z=3",
		"x=1, y=2, z=3>",
		"<x=1, y=2>",
		"<x=1, y=2, w=3>",
		"<x=1, y=a, z=3>",
		"<x=1,y=2,z=3>",
	} {
		for _, solve := range []func(io.Reader) (string, error){day12p01, day12p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("expected moon error for %q", input)
			}
		}
	}
}
