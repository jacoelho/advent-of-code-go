package aoc2023

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day23p01(t *testing.T) {
	example := `#.#####################
#.......#########...###
#######.#########.#.###
###.....#.>.>.###.#.###
###v#####.#v#.###.#.###
###.>...#.#.#.....#...#
###v###.#.#.#########.#
###...#.#.#.......#...#
#####.#.#.#######.#.###
#.....#.#.#.......#...#
#.#####.#.#.#########v#
#.#...#...#...###...>.#
#.#.#v#######v###.###v#
#...#.>.#...>.>.#.###.#
#####v#.#.###v#.#.###.#
#.....#...#...#.#.#...#
#.#########.###.#.#.###
#...###...#...#...#.###
###.###.#.###v#####v###
#...#...#.#.>.>.#.>.###
#.###.###.#.###.#.#v###
#.....###...###...#...#
#####################.#`

	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(example),
			Want:  "94",
		},
		{
			Input: aoc.FileInput(t, 2023, 23),
			Want:  "2306",
		},
	}
	aoc.AOCTest(t, day23p01, tests)
}

func Test_day23p02(t *testing.T) {
	example := `#.#####################
#.......#########...###
#######.#########.#.###
###.....#.>.>.###.#.###
###v#####.#v#.###.#.###
###.>...#.#.#.....#...#
###v###.#.#.#########.#
###...#.#.#.......#...#
#####.#.#.#######.#.###
#.....#.#.#.......#...#
#.#####.#.#.#########v#
#.#...#...#...###...>.#
#.#.#v#######v###.###v#
#...#.>.#...>.>.#.###.#
#####v#.#.###v#.#.###.#
#.....#...#...#.#.#...#
#.#########.###.#.#.###
#...###...#...#...#.###
###.###.#.###v#####v###
#...#...#.#.>.>.#.>.###
#.###.###.#.###.#.#v###
#.....###...###...#...#
#####################.#`

	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(example),
			Want:  "154",
		},
		{
			Input: aoc.FileInput(t, 2023, 23),
			Want:  "6718",
		},
	}
	aoc.AOCTest(t, day23p02, tests)
}

func TestHikingCorridors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  [2]int
	}{
		{
			name:  "adjacent junctions",
			input: "#.###\n#...#\n#...#\n###.#",
			want:  [2]int{7, 7},
		},
		{
			name:  "degree-two end",
			input: "#.##\n#..#\n##..",
			want:  [2]int{3, 3},
		},
		{
			name:  "parallel corridors",
			input: "##.####\n#.....#\n#.###.#\n#.###.#\n#.....#\n###.###",
			want:  [2]int{10, 10},
		},
		{
			name:  "slope at junction",
			input: "#.###\n#>..#\n#...#\n###.#",
			want:  [2]int{5, 7},
		},
		{
			name:  "backward slope",
			input: "#.\n#^\n#.",
			want:  [2]int{-1, 2},
		},
		{
			name:  "start is end",
			input: ".",
			want:  [2]int{0, 0},
		},
	}
	for _, tt := range tests {
		for part, solve := range []func(io.Reader) (string, error){day23p01, day23p02} {
			t.Run(fmt.Sprintf("%s/part%d", tt.name, part+1), func(t *testing.T) {
				got, err := solve(strings.NewReader(tt.input))
				if tt.want[part] < 0 {
					if err == nil || !strings.Contains(err.Error(), "no hiking route") {
						t.Fatalf("got %q, %v; want no-route error", got, err)
					}
					return
				}
				if err != nil || got != strconv.Itoa(tt.want[part]) {
					t.Fatalf("got %q, %v; want %d", got, err, tt.want[part])
				}
			})
		}
	}
}

func TestHikingSmallMazes(t *testing.T) {
	for mask := range 64 {
		for variant := range 2 {
			rows := []string{"#.##", "####", "####", "####", "##.#"}
			for y := range 3 {
				row := []byte(rows[y+1])
				for x := range 2 {
					bit := 2*y + x
					if mask&(1<<bit) != 0 {
						row[x+1] = '.'
						if variant == 1 {
							row[x+1] = "^>v<."[(bit+mask)%5]
						}
					}
				}
				rows[y+1] = string(row)
			}
			input := strings.Join(rows, "\n")
			for part, solve := range []func(io.Reader) (string, error){day23p01, day23p02} {
				t.Run(fmt.Sprintf("mask%d/variant%d/part%d", mask, variant, part+1), func(t *testing.T) {
					want := longestCellHike(rows, part == 0)
					got, err := solve(strings.NewReader(input))
					if want < 0 {
						if err == nil || !strings.Contains(err.Error(), "no hiking route") {
							t.Fatalf("maze:\n%s\ngot %q, %v; want no-route error", input, got, err)
						}
						return
					}
					if err != nil || got != strconv.Itoa(want) {
						t.Fatalf("maze:\n%s\ngot %q, %v; want %d", input, got, err, want)
					}
				})
			}
		}
	}
}

func longestCellHike(rows []string, slippery bool) int {
	endX, endY := strings.IndexByte(rows[len(rows)-1], '.'), len(rows)-1
	visited := make(map[[2]int]bool)
	directions := [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	longest := -1

	var walk func(int, int, int)
	walk = func(x, y, distance int) {
		if x == endX && y == endY {
			longest = max(longest, distance)
			return
		}
		visited[[2]int{x, y}] = true
		for i, direction := range directions {
			if arrow := strings.IndexByte("^>v<", rows[y][x]); slippery && arrow >= 0 && i != arrow {
				continue
			}
			nx, ny := x+direction[0], y+direction[1]
			if ny < 0 || ny >= len(rows) || nx < 0 || nx >= len(rows[ny]) {
				continue
			}
			if rows[ny][nx] != '#' && !visited[[2]int{nx, ny}] {
				walk(nx, ny, distance+1)
			}
		}
		delete(visited, [2]int{x, y})
	}
	walk(strings.IndexByte(rows[0], '.'), 0, 0)
	return longest
}

func TestHikingLargeChain(t *testing.T) {
	network := trailNetwork{
		segments: make([][]trailSegment, 70),
		startIdx: 0,
		endIdx:   69,
	}
	for i := range 69 {
		network.segments[i] = append(network.segments[i], trailSegment{toIdx: i + 1, distance: 1})
		network.segments[i+1] = append(network.segments[i+1], trailSegment{toIdx: i, distance: 1})
	}
	got, err := findLongestHike(network)
	if err != nil || got != 69 {
		t.Fatalf("got %d, %v; want 69", got, err)
	}
}

func TestHikingJunctionLimit(t *testing.T) {
	input := strings.Repeat(strings.Repeat(".", 65)+"\n", 65)
	for part, solve := range []func(io.Reader) (string, error){day23p01, day23p02} {
		t.Run(fmt.Sprintf("part%d", part+1), func(t *testing.T) {
			got, err := solve(strings.NewReader(input))
			if err == nil || !strings.Contains(err.Error(), "exceeds") || !strings.Contains(err.Error(), "junctions") {
				t.Fatalf("got %q, %v; want junction-limit error", got, err)
			}
		})
	}
}

func TestDay23ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day23p01, "#.\n#.\n"},
		{"part2", day23p02, "#.\n#.\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range []struct{ name, prefix string }{{"immediate", ""}, {"after prefix", tt.prefix}} {
				t.Run(mode.name, func(t *testing.T) {
					input := io.MultiReader(strings.NewReader(mode.prefix), iotest.ErrReader(failure))
					_, err := tt.solve(input)
					if !errors.Is(err, failure) {
						t.Fatalf("got %v, want reader failure", err)
					}
				})
			}
		})
	}
}
