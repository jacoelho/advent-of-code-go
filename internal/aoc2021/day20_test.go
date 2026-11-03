package aoc2021

import (
	"io"
	"strings"
	"testing"

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
