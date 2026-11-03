package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

const monkeyMapExample = `        ...#
        .#..
        #...
        ....
...#.......#
........#...
..#....#....
..........#.
        ...#....
        .....#..
        .#......
        ......#.

10R5L5R10L4R5L5`

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(monkeyMapExample), Want: "6032"},
		{Input: strings.NewReader(monkeyMapExample + "\n"), Want: "6032"},
		{Input: strings.NewReader(". .\n\n1"), Want: "1012"},
		{Input: strings.NewReader(". .\n\n2"), Want: "1004"},
		{Input: aoc.FileInput(t, 2022, 22), Want: "88226"},
	}
	aoc.AOCTest(t, day22p01, tests)
}

func Test_day22p02(t *testing.T) {
	drawing, path, _ := strings.Cut(monkeyMapExample, "\n\n")
	shifted := " " + strings.ReplaceAll(drawing, "\n", "\n ") + "\n\n" + path
	tests := []aoc.TestInput{
		{Input: strings.NewReader(monkeyMapExample), Want: "5031"},
		{Input: strings.NewReader(shifted), Want: "5035"},
		{Input: strings.NewReader(" .\n...\n .\n .\n\n1"), Want: "2013"},
		{Input: strings.NewReader(" .\n..#\n .\n .\n\n1"), Want: "1008"},
		{Input: strings.NewReader("  .\n ...\n  .\n  .\n\n1"), Want: "2017"},
		{Input: strings.NewReader("  .\n ..#\n  .\n  .\n\n1"), Want: "1012"},
		{Input: aoc.FileInput(t, 2022, 22), Want: "57305"},
	}
	aoc.AOCTest(t, day22p02, tests)
}

func TestCubeEdgesReverse(t *testing.T) {
	for _, input := range []string{monkeyMapExample, " .\n...\n .\n .\n\n0"} {
		board, err := parseMonkeyMap(strings.NewReader(input))
		if err != nil {
			t.Fatal(err)
		}
		cube, err := foldCube(board)
		if err != nil {
			t.Fatal(err)
		}
		for _, face := range cube.faces {
			for facing := range 4 {
				for offset := range cube.side {
					local := [...]grid.Position2D[int]{
						{X: cube.side - 1, Y: offset},
						{X: offset, Y: cube.side - 1},
						{X: 0, Y: offset},
						{X: offset, Y: 0},
					}[facing]
					start := face.origin.Add(local)
					next, direction, err := cube.wrap(start, facing, nil)
					if err != nil {
						t.Fatal(err)
					}
					back, reverse, err := cube.wrap(next, (direction+2)%4, nil)
					if err != nil || back != start || reverse != (facing+2)%4 {
						t.Fatalf("edge %v/%d: reverse returned %v/%d, %v", start, facing, back, reverse, err)
					}
				}
			}
		}
	}
}

func TestCubeInvalidNets(t *testing.T) {
	for _, input := range []string{
		"..\n..\n..\n\n0",
		". . . . . .\n\n0",
		" .\n..\n .\n .\n\n0",
		"......\n\n0",
	} {
		if _, err := day22p02(strings.NewReader(input)); err == nil {
			t.Fatalf("expected invalid cube net error for %q", input)
		}
	}
}

func TestMapInvalidInput(t *testing.T) {
	for _, input := range []string{"", ".\n\n", "#\n\n1", ".\n\n1R", ".\n\n1X2", ".\n\nR1", "x\n\n1"} {
		if _, err := day22p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected map input error for %q", input)
		}
	}
}

func TestDay22ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day22p01},
		{"part2", day22p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
