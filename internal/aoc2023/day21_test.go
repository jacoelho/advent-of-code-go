package aoc2023

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day21p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("S"), Want: "0"},
		{Input: strings.NewReader("###\n#S#\n###"), Want: "0"},
		{
			Input: aoc.FileInput(t, 2023, 21),
			Want:  "3598",
		},
	}
	aoc.AOCTest(t, day21p01, tests)
}

func Test_day21p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("S"), Want: "702322399865956"},
		{Input: strings.NewReader("###\n.S.\n###"), Want: "26501366"},
		{Input: strings.NewReader("###\n#S.\n###"), Want: "1"},
		{Input: strings.NewReader("###\n#S#\n###"), Want: "0"},
		{Input: strings.NewReader("S##\n###\n#.."), Want: "0"},
		{
			Input: aoc.FileInput(t, 2023, 21),
			Want:  "601441063166538",
		},
	}
	aoc.AOCTest(t, day21p02, tests)
}

func TestGardenExample(t *testing.T) {
	const input = `...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........`
	got, err := day21(strings.NewReader(input), 6)
	if err != nil || got != "16" {
		t.Fatalf("six steps = %q, %v; want 16", got, err)
	}
}

func TestGardenSmallWalks(t *testing.T) {
	for _, input := range []string{
		"S",
		"S..\n...\n...",
		"###\n.S.\n###",
		"###\n#S.\n###",
		"###\n#S#\n###",
		"S##\n###\n#..",
	} {
		t.Run(input, func(t *testing.T) {
			garden, start, err := parseGarden(strings.NewReader(input))
			if err != nil {
				t.Fatal(err)
			}
			finite := gardenFrontiers(input, 12, false)
			repeated := gardenFrontiers(input, 12, true)
			for steps := range 13 {
				got, err := day21(strings.NewReader(input), steps)
				if err != nil || got != strconv.Itoa(finite[steps]) {
					t.Fatalf("finite steps %d = %q, %v; want %d", steps, got, err, finite[steps])
				}
				count, err := repeatedGardenPlots(garden, start, steps)
				if err != nil {
					t.Fatalf("repeated steps %d: %v", steps, err)
				}
				if count.String() != strconv.Itoa(repeated[steps]) {
					t.Fatalf("repeated steps %d = %s; want %d", steps, count, repeated[steps])
				}
			}
		})
	}
}

func TestGardenWallPatterns(t *testing.T) {
	for mask := range 256 {
		var input strings.Builder
		bit := 0
		for cell := range 9 {
			switch {
			case cell == 4:
				input.WriteByte('S')
			default:
				if mask&(1<<bit) != 0 {
					input.WriteByte('#')
				} else {
					input.WriteByte('.')
				}
				bit++
			}
			if cell%3 == 2 && cell != 8 {
				input.WriteByte('\n')
			}
		}
		garden, start, err := parseGarden(strings.NewReader(input.String()))
		if err != nil {
			t.Fatal(err)
		}
		_, err = repeatedGardenPlots(garden, start, 12)
		if err != nil {
			if !strings.Contains(err.Error(), "unsupported garden distances") {
				t.Fatalf("mask %d: %v", mask, err)
			}
			continue
		}
		want := gardenFrontiers(input.String(), 12, true)
		for steps := range 13 {
			got, err := repeatedGardenPlots(garden, start, steps)
			if err != nil || got.String() != strconv.Itoa(want[steps]) {
				t.Fatalf("mask %d, steps %d = %v, %v; want %d", mask, steps, got, err, want[steps])
			}
		}
	}
}

func gardenFrontiers(input string, steps int, repeated bool) []int {
	rows := strings.Split(input, "\n")
	width, height := len(rows[0]), len(rows)
	frontier := make(map[[2]int]bool)
	for y, row := range rows {
		for x := range row {
			if row[x] == 'S' {
				frontier[[2]int{x, y}] = true
			}
		}
	}
	counts := []int{len(frontier)}
	for range steps {
		next := make(map[[2]int]bool)
		for point := range frontier {
			for _, offset := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
				p := [2]int{point[0] + offset[0], point[1] + offset[1]}
				x, y := p[0], p[1]
				if repeated {
					x = (x%width + width) % width
					y = (y%height + height) % height
				} else if x < 0 || y < 0 || x >= width || y >= height {
					continue
				}
				if rows[y][x] != '#' {
					next[p] = true
				}
			}
		}
		frontier = next
		counts = append(counts, len(frontier))
	}
	return counts
}

func TestGardenUnsupported(t *testing.T) {
	for _, input := range []string{
		"..#\n#S.\n.#.",
		"S.\n..",
		"S..\n...",
		"S" + strings.Repeat(".", 202) + strings.Repeat("\n"+strings.Repeat(".", 203), 202),
	} {
		_, err := day21p02(strings.NewReader(input))
		if err == nil || !strings.Contains(err.Error(), "unsupported garden") {
			t.Fatalf("got %v; want unsupported garden", err)
		}
	}
}

func TestGardenCountBeyondInt64(t *testing.T) {
	garden, start, err := parseGarden(strings.NewReader("S"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := repeatedGardenPlots(garden, start, 4000000000)
	if err != nil || got.String() != "16000000008000000001" {
		t.Fatalf("open plane count = %v, %v; want 16000000008000000001", got, err)
	}
}

func TestGardenInputErrors(t *testing.T) {
	for _, input := range []string{"", ".", "SS", "S.\n.", "Sx", "S\n\n"} {
		for _, solve := range []func(io.Reader) (string, error){day21p01, day21p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected error", input)
			}
		}
	}
	if _, err := day21(strings.NewReader("S"), -1); err == nil {
		t.Fatal("negative finite steps: expected error")
	}
	garden, start, err := parseGarden(strings.NewReader("S"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repeatedGardenPlots(garden, start, -1); err == nil {
		t.Fatal("negative repeated steps: expected error")
	}
}

func TestDay21ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day21p01, "S\n"},
		{"part2", day21p02, "S\n"},
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
