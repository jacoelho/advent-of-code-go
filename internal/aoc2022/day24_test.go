package aoc2022

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const valleyExample = `#.######
#>>.<^<#
#.<..<<#
#>v.><>#
#<^v^^>#
######.#`

func Test_day24p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(valleyExample), Want: "18"},
		{Input: strings.NewReader("#.#\n#.#\n#.#"), Want: "2"},
		{Input: strings.NewReader("#.##\n#.>#\n##.#"), Want: "4"},
		{Input: aoc.FileInput(t, 2022, 24), Want: "322"},
	}
	aoc.AOCTest(t, day24p01, tests)
}

func Test_day24p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(valleyExample), Want: "54"},
		{Input: strings.NewReader("#.#\n#.#\n#.#"), Want: "6"},
		{Input: strings.NewReader("#.##\n#.>#\n##.#"), Want: "10"},
		{Input: aoc.FileInput(t, 2022, 24), Want: "974"},
	}
	aoc.AOCTest(t, day24p02, tests)
}

func TestValleyFailures(t *testing.T) {
	for _, input := range []string{
		"",
		"#.#\n#^#\n#.#",
		"#..#\n#..#\n##.#",
		"#.#\n###\n#.#",
		"#.#\n#..#\n#.#",
		"#.#\n..#\n#.#",
	} {
		if _, err := day24p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected valley error for %q", input)
		}
	}
}
