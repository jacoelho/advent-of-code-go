package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

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

func TestDay24ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day24p01},
		{"part2", day24p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
