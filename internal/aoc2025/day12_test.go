package aoc2025

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(giftExampleShapes + "\n\n4x4: 0 0 0 0 2 0\n12x5: 1 0 1 0 2 2\n12x5: 1 0 1 0 3 2"), Want: "2"},
		{
			Input: aoc.FileInput(t, 2025, 12),
			Want:  "505",
		},
	}
	aoc.AOCTest(t, day12p01, tests)
}

const giftExampleShapes = `0:
###
##.
##.

1:
###
##.
.##

2:
.##
###
##.

3:
##.
###
##.

4:
###
#..
###

5:
###
.#.
###`

func TestGiftRegions(t *testing.T) {
	for _, tt := range []struct{ region, want string }{
		{"4x4: 0 0 0 0 2 0", "1"},
		{"12x5: 1 0 1 0 2 2", "1"},
		{"12x5: 1 0 1 0 3 2", "0"},
	} {
		got, err := day12p01(strings.NewReader(giftExampleShapes + "\n\n" + tt.region))
		if err != nil || got != tt.want {
			t.Fatalf("region %q = %q, %v; want %s", tt.region, got, err, tt.want)
		}
	}
}

func TestGiftGeometry(t *testing.T) {
	for _, tt := range []struct{ name, input, want string }{
		{"dimensions", "0:\n##\n##\n\n1x4: 1", "0"},
		{"rotation", "0:\n##\n\n1x2: 1", "1"},
		{"interlocking", "0:\n##\n#.\n\n3x2: 2", "1"},
		{"reflection", "0:\n##.\n.##\n\n1:\n##..\n#..#\n####\n\n4x3: 1 1", "1"},
		{"empty demand", "0:\n#\n\n1x1: 0", "1"},
		{"padding", "0:\n...\n.#.\n...\n\n1x1: 1", "1"},
		{"large certified", "0:\n#\n\n1000x1000: 1000000", "1"},
		{"area rejection", "0:\n#\n\n1x1: 2", "0"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := day12p01(strings.NewReader(tt.input))
			if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %s", got, err, tt.want)
			}
		})
	}
}

func TestGiftSmallPackings(t *testing.T) {
	const shapes = "0:\n##\n\n1:\n##\n#.\n\n2:\n##\n##\n\n"
	for width := 1; width <= 3; width++ {
		for height := 1; height <= 3; height++ {
			for dominoes := range 3 {
				for corners := range 3 {
					for squares := range 3 {
						if dominoes+corners+squares > 3 {
							continue
						}
						counts := [3]int{dominoes, corners, squares}
						want := "0"
						if giftsFitByCells(width, height, counts) {
							want = "1"
						}
						input := shapes + fmt.Sprintf("%dx%d: %d %d %d", width, height, dominoes, corners, squares)
						got, err := day12p01(strings.NewReader(input))
						if err != nil || got != want {
							t.Fatalf("%dx%d counts %v = %q, %v; want %s", width, height, counts, got, err, want)
						}
					}
				}
			}
		}
	}
}

func giftsFitByCells(width, height int, counts [3]int) bool {
	orientations := [][][][2]int{
		{{{0, 0}, {1, 0}}, {{0, 0}, {0, 1}}},
		{
			{{0, 0}, {1, 0}, {0, 1}},
			{{0, 0}, {1, 0}, {1, 1}},
			{{0, 0}, {0, 1}, {1, 1}},
			{{1, 0}, {0, 1}, {1, 1}},
		},
		{{{0, 0}, {1, 0}, {0, 1}, {1, 1}}},
	}
	occupied := make(map[[2]int]bool)
	var place func([3]int) bool
	place = func(remaining [3]int) bool {
		kind := 0
		for kind < len(remaining) && remaining[kind] == 0 {
			kind++
		}
		if kind == len(remaining) {
			return true
		}
		remaining[kind]--
		for _, orientation := range orientations[kind] {
			for y := range height {
				for x := range width {
					var cells [][2]int
					valid := true
					for _, cell := range orientation {
						point := [2]int{x + cell[0], y + cell[1]}
						if point[0] >= width || point[1] >= height || occupied[point] {
							valid = false
							break
						}
						cells = append(cells, point)
					}
					if !valid {
						continue
					}
					for _, cell := range cells {
						occupied[cell] = true
					}
					if place(remaining) {
						return true
					}
					for _, cell := range cells {
						delete(occupied, cell)
					}
				}
			}
		}
		return false
	}
	return place(counts)
}

func TestGiftInputErrors(t *testing.T) {
	for _, input := range []string{
		"",
		"0:\nx\n\n1x1: 1",
		"0:\n##\n#\n\n2x2: 1",
		"0:\n.\n\n1x1: 1",
		"1:\n#\n\n1x1: 1",
		"0:\n#\n\n0:\n##\n\n2x2: 1 1",
		"0:\n#\n\n1x1: 1 0",
		"0:\n#\n\n1x1: -1",
		"0:\n#\n\n-1x1: 1",
		"0:\n#\n\n1x1: 9223372036854775808",
		"0:\n#\n\n1,1: 1",
	} {
		if _, err := day12p01(strings.NewReader(input)); err == nil {
			t.Fatalf("input %q: expected error", input)
		}
	}
	const input = "0:\n#\n\n1x1: 1"
	for _, formatted := range []string{input, input + "\n", strings.ReplaceAll(input+"\n", "\n", "\r\n")} {
		got, err := day12p01(strings.NewReader(formatted))
		if err != nil || got != "1" {
			t.Fatalf("line endings = %q, %v; want 1", got, err)
		}
	}
}

func TestGiftSearchLimits(t *testing.T) {
	shapes, regions, err := parseGiftSummary(strings.NewReader("0:\n##\n#.\n\n3x2: 2"))
	if err != nil {
		t.Fatal(err)
	}
	work := giftPackingWork{states: maxGiftStates}
	if _, err := packGiftRegion(shapes, regions[0], &work); err == nil {
		t.Fatal("exhausted exact search must return an error")
	}
	_, err = day12p01(strings.NewReader("0:\n##\n#.\n\n9x9: 27"))
	if err == nil || !strings.Contains(err.Error(), "unsupported gift region") {
		t.Fatalf("got %v; want unsupported exact-search region", err)
	}
}

func TestDay12ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day12p01, "0:\n#\n\n1x1: 1\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				input := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(input); !errors.Is(err, failure) {
					t.Fatalf("got %v; want reader failure", err)
				}
			}
		})
	}
}
