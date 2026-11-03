package aoc2024

import (
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
)

func guardAnswersByCoordinates(input string) (int, int) {
	type point struct{ x, y int }
	type state struct{ x, y, dx, dy int }
	cells := make(map[point]byte)
	var start point
	y := 0
	for row := range strings.SplitSeq(input, "\n") {
		row = strings.TrimSuffix(row, "\r")
		for x := range len(row) {
			cells[point{x, y}] = row[x]
			if row[x] == '^' {
				start = point{x, y}
			}
		}
		y++
	}
	walk := func(obstacle point) (map[point]bool, bool) {
		visited := make(map[point]bool)
		seen := make(map[state]bool)
		current := state{start.x, start.y, 0, -1}
		for cells[point{current.x, current.y}] != 0 {
			if seen[current] {
				return visited, true
			}
			seen[current] = true
			visited[point{current.x, current.y}] = true
			next := point{current.x + current.dx, current.y + current.dy}
			if cells[next] == '#' || next == obstacle {
				current.dx, current.dy = -current.dy, current.dx
			} else {
				current.x, current.y = next.x, next.y
			}
		}
		return visited, false
	}
	visited, _ := walk(point{-1, -1})
	loops := 0
	for obstacle := range visited {
		if obstacle != start {
			if _, loopsHere := walk(obstacle); loopsHere {
				loops++
			}
		}
	}
	return len(visited), loops
}

func checkGuardAnswers(t *testing.T, input string, wantVisited, wantLoops int) {
	t.Helper()
	visited, err := day06p01(strings.NewReader(input))
	if err != nil || visited != strconv.Itoa(wantVisited) {
		t.Fatalf("part 1 for %q: got %q, %v; want %d", input, visited, err, wantVisited)
	}
	loops, err := day06p02(strings.NewReader(input))
	if err != nil || loops != strconv.Itoa(wantLoops) {
		t.Fatalf("part 2 for %q: got %q, %v; want %d", input, loops, err, wantLoops)
	}
}

func TestGuardBoundaries(t *testing.T) {
	for _, tt := range []struct {
		name, input    string
		visited, loops int
	}{
		{"turn and revisit", ".#..\n.^.#\n....", 3, 0},
		{"cyclic baseline", "######\n#....#\n#....#\n#.^..#\n#....#\n######", 14, 13},
		{"ragged exit", ".\n.^", 1, 0},
		{"empty row", ".\n\n^", 1, 0},
		{"unknown tile", "x\n^", 2, 0},
		{"NUL exit", "\x00\n^", 1, 0},
		{"UTF-8 bytes", "é\n.^", 2, 0},
		{"turns in place", "###\n#^#\n###", 1, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			visited, loops := guardAnswersByCoordinates(tt.input)
			if visited != tt.visited || loops != tt.loops {
				t.Fatalf("reference for %q: got %d/%d; want %d/%d", tt.input, visited, loops, tt.visited, tt.loops)
			}
			checkGuardAnswers(t, tt.input, tt.visited, tt.loops)
		})
	}
}

func TestGuardAgainstCoordinates(t *testing.T) {
	const seed = 202406
	random := rand.New(rand.NewPCG(seed, 0))
	for trial := range 400 {
		height, width := 1+random.IntN(6), 1+random.IntN(6)
		rows := make([][]byte, height)
		guardY := random.IntN(height)
		for y := range rows {
			rowWidth := width
			if trial%2 != 0 {
				rowWidth = random.IntN(7)
			}
			if y == guardY && rowWidth == 0 {
				rowWidth = 1
			}
			rows[y] = make([]byte, rowWidth)
			for x := range rows[y] {
				switch n := random.IntN(20); {
				case n < 6:
					rows[y][x] = '#'
				case n == 6:
					rows[y][x] = 0
				case n == 7:
					rows[y][x] = 'x'
				case n == 8:
					rows[y][x] = 0xff
				default:
					rows[y][x] = '.'
				}
			}
		}
		rows[guardY][random.IntN(len(rows[guardY]))] = '^'
		var input strings.Builder
		for y, row := range rows {
			if y > 0 {
				input.WriteByte('\n')
			}
			input.Write(row)
		}
		t.Run(strconv.Itoa(trial), func(t *testing.T) {
			t.Logf("seed %d, input %q", seed, input.String())
			visited, loops := guardAnswersByCoordinates(input.String())
			checkGuardAnswers(t, input.String(), visited, loops)
		})
	}
}
