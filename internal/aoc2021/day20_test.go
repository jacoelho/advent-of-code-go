package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(imageExample), Want: "35"},
		{Input: aoc.FileInput(t, 2021, 20), Want: "5573"},
	}

	aoc.AOCTest(t, day20p01, tests)
}

func Test_day20p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(imageExample), Want: "3351"},
		{Input: aoc.FileInput(t, 2021, 20), Want: "20097"},
	}

	aoc.AOCTest(t, day20p02, tests)
}

func TestImageBackground(t *testing.T) {
	for _, algorithm := range []string{strings.Repeat(".", 512), "#" + strings.Repeat(".", 511)} {
		for _, solve := range []func(io.Reader) (string, error){day20p01, day20p02} {
			got, err := solve(strings.NewReader(algorithm + "\n\n.\n"))
			if err != nil {
				t.Fatal(err)
			}
			if got != "0" {
				t.Fatalf("uniform dark image after even steps = %s, want 0", got)
			}
		}
	}
	for _, solve := range []func(io.Reader) (string, error){day20p01, day20p02} {
		if _, err := solve(strings.NewReader(strings.Repeat("#", 512) + "\n\n.\n")); err == nil {
			t.Fatal("expected infinite lit image error")
		}
	}
}

func TestImageErrors(t *testing.T) {
	for _, input := range []string{
		".#\n\n#",
		strings.Repeat(".", 512) + "\n#",
		strings.Repeat(".", 512) + "\n\n#x",
		strings.Repeat(".", 512) + "\n\n##\n#",
	} {
		if _, err := day20p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected invalid image error for %q", input)
		}
	}
}

const imageExample = `..#.#..#####.#.#.#.###.##.....###.##.#..###.####..#####..#....#..#..##..###..######.###...####..#..#####..##..#.#####...##.#.#..#.##..#.#......#.###.######.###.####...#.##.##..#..#..#####.....#.#....###..#.##......#.....#..#..#..##..#...##.######.####.####.#.#...#.......#..#.#.#...####.##.#......#..#...##.#.##..#...##.#.##..###.#......#.#.......#.#.#.####.###.##...#.....####.#..#..#.##.#....##..#.####....##...##..#...#......#.#.......#.......##..####..#...#.#.#...##..#.#..###..#####........#..####......#..#

#..#.
#....
##..#
..#..
..###`

func TestDay20ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day20p01},
		{"part2", day20p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
