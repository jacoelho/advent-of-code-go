package aoc2023

import (
	"errors"
	"io"
	"maps"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533`),
			Want: "102",
		},
		{
			Input: aoc.FileInput(t, 2023, 17),
			Want:  "755",
		},
	}
	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533`),
			Want: "94",
		},
		{
			Input: strings.NewReader(`111111111111
999999999991
999999999991
999999999991
999999999991`),
			Want: "71",
		},
		{
			Input: aoc.FileInput(t, 2023, 17),
			Want:  "881",
		},
	}
	aoc.AOCTest(t, day17p02, tests)
}

func TestCrucibleStrides(t *testing.T) {
	tests := []struct {
		name  string
		input string
		part1 string
		part2 string
	}{
		{"singleton", "7", "0", ""},
		{"three steps", "1234", "9", ""},
		{"four steps", "12345", "", "14"},
		{"vertical start", "1\n2\n3\n4\n5", "", "14"},
		{"ten steps", "11111111111", "", "10"},
		{"eleven steps", "111111111111", "", ""},
		{"zero heat and final stride", "000009\n000009\n000009\n000009\n999900", "0", "36"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for part, solve := range []func(io.Reader) (string, error){day17p01, day17p02} {
				want := []string{tt.part1, tt.part2}[part]
				got, err := solve(strings.NewReader(tt.input))
				if want == "" {
					if err == nil {
						t.Errorf("part %d returned %q for an unreachable destination", part+1, got)
					}
				} else if err != nil || got != want {
					t.Errorf("part %d = %q, %v; want %q", part+1, got, err, want)
				}
			}
		})
	}
}

func TestHeatGridErrors(t *testing.T) {
	for _, input := range []string{"", "1\n\n2", "12\n3", "1x", strings.Repeat(strings.Repeat("1", 500)+"\n", 501)} {
		for _, solve := range []func(io.Reader) (string, error){day17p01, day17p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("expected an error for invalid heat grid of %d bytes", len(input))
			}
		}
	}
	failure := errors.New("heat reader failed")
	for _, solve := range []func(io.Reader) (string, error){day17p01, day17p02} {
		input := io.MultiReader(strings.NewReader("12\n34\n"), iotest.ErrReader(failure))
		if _, err := solve(input); !errors.Is(err, failure) {
			t.Fatalf("error = %v, want reader failure", err)
		}
	}
}

func TestCrucibleTinyGrids(t *testing.T) {
	random := rand.New(rand.NewPCG(17, 2023))
	for height := 1; height <= 7; height++ {
		for width := 1; width <= 7; width++ {
			rows := make([]string, height)
			for y := range rows {
				row := make([]byte, width)
				for x := range row {
					row[x] = byte('0' + random.IntN(10))
				}
				rows[y] = string(row)
			}
			for _, limits := range [][2]int{{1, 3}, {4, 10}} {
				want, reachable := crucibleStepOracle(rows, limits[0], limits[1])
				got, err := crucibleHeat(strings.NewReader(strings.Join(rows, "\n")), limits[0], limits[1])
				if !reachable {
					if err == nil {
						t.Fatalf("%v strides %v returned %q for an unreachable destination", rows, limits, got)
					}
				} else if err != nil || got != strconv.Itoa(want) {
					t.Fatalf("%v strides %v = %q, %v; want %d", rows, limits, got, err, want)
				}
			}
		}
	}
}

func crucibleStepOracle(rows []string, minSteps, maxSteps int) (int, bool) {
	type state struct {
		x, y, direction, steps int
	}
	width, height := len(rows[0]), len(rows)
	if width == 1 && height == 1 {
		return 0, minSteps == 1
	}
	directions := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	distance := map[state]int{{direction: -1}: 0}
	// Bellman–Ford relaxes one-cell moves, retaining the full direction/run state.
	for pass := 0; pass < 4*width*height*maxSteps; pass++ {
		next := maps.Clone(distance)
		changed := false
		for current, cost := range distance {
			for direction, delta := range directions {
				steps := 1
				if current.direction >= 0 {
					if direction == (current.direction+2)%4 {
						continue
					}
					if direction == current.direction {
						steps = current.steps + 1
					} else if current.steps < minSteps {
						continue
					}
				}
				x, y := current.x+delta[0], current.y+delta[1]
				if steps > maxSteps || x < 0 || y < 0 || x >= width || y >= height {
					continue
				}
				to := state{x: x, y: y, direction: direction, steps: steps}
				candidate := cost + int(rows[y][x]-'0')
				if previous, ok := next[to]; !ok || candidate < previous {
					next[to] = candidate
					changed = true
				}
			}
		}
		distance = next
		if !changed {
			break
		}
	}
	best := math.MaxInt
	for at, cost := range distance {
		if at.x == width-1 && at.y == height-1 && at.steps >= minSteps {
			best = min(best, cost)
		}
	}
	return best, best != math.MaxInt
}
