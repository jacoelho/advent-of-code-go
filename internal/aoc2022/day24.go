package aoc2022

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type blizzardValley struct {
	cells         []byte
	width, height int
	period        int
	start, end    grid.Position2D[int]
}

func parseValley(r io.Reader) (blizzardValley, error) {
	var valley blizzardValley
	var rows []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		row := s.Text()
		if len(row) < 3 || len(rows) > 0 && len(row) != len(rows[0]) {
			return valley, fmt.Errorf("valley must be a rectangle at least three cells wide")
		}
		if len(rows) >= 1_000_000/len(row) {
			return valley, fmt.Errorf("valley exceeds one million cells")
		}
		rows = append(rows, row)
	}
	if err := s.Err(); err != nil {
		return valley, err
	}
	if len(rows) < 3 {
		return valley, fmt.Errorf("valley requires top and bottom walls with an interior")
	}
	valley.width, valley.height = len(rows[0])-2, len(rows)-2
	openings := [2]int{}
	for y, row := range rows {
		for x := range row {
			c := row[x]
			if x == 0 || x == len(row)-1 || y == 0 || y == len(rows)-1 {
				if c == '#' {
					continue
				}
				if c != '.' || x == 0 || x == len(row)-1 || y != 0 && y != len(rows)-1 {
					return valley, fmt.Errorf("invalid valley boundary at %d,%d", x, y)
				}
				if y == 0 {
					valley.start = grid.Position2D[int]{X: x, Y: y}
					openings[0]++
				} else {
					valley.end = grid.Position2D[int]{X: x, Y: y}
					openings[1]++
				}
				continue
			}
			if c != '.' && c != '>' && c != '<' && c != '^' && c != 'v' {
				return valley, fmt.Errorf("invalid valley interior cell %q", c)
			}
			valley.cells = append(valley.cells, c)
		}
	}
	if openings != [2]int{1, 1} {
		return valley, fmt.Errorf("valley requires exactly one top and one bottom opening")
	}
	valley.period = xmath.LCM(valley.width, valley.height)
	if valley.period > 4_000_000/(len(valley.cells)+2) {
		return valley, fmt.Errorf("valley exceeds four million position-phase states")
	}
	return valley, nil
}

func (valley blizzardValley) positionID(position grid.Position2D[int]) (int, bool) {
	if position == valley.start {
		return len(valley.cells), true
	}
	if position == valley.end {
		return len(valley.cells) + 1, true
	}
	if position.X < 1 || position.X > valley.width || position.Y < 1 || position.Y > valley.height {
		return 0, false
	}
	return (position.Y-1)*valley.width + position.X - 1, true
}

func (valley blizzardValley) blocked(position grid.Position2D[int], minute int) bool {
	if position == valley.start || position == valley.end {
		return false
	}
	x, y := position.X-1, position.Y-1
	w, h := valley.width, valley.height
	// Inverse motion checks the original map without storing every hazard phase.
	return valley.cells[y*w+xmath.Modulo(x-minute, w)] == '>' ||
		valley.cells[y*w+xmath.Modulo(x+minute, w)] == '<' ||
		valley.cells[xmath.Modulo(y-minute, h)*w+x] == 'v' ||
		valley.cells[xmath.Modulo(y+minute, h)*w+x] == '^'
}

func (valley blizzardValley) travel(from, to grid.Position2D[int], minute int) (int, error) {
	type state struct {
		position grid.Position2D[int]
		minute   int
	}
	positions := len(valley.cells) + 2
	seen := make([]bool, positions*valley.period)
	id, _ := valley.positionID(from)
	seen[(minute%valley.period)*positions+id] = true
	queue := collections.NewDeque[state](64)
	queue.PushBack(state{position: from, minute: minute})
	offsets := append(grid.OffsetsNeighbours4[int](), grid.Position2D[int]{})
	for queue.Size() > 0 {
		current, _ := queue.PopFront()
		if current.position == to {
			return current.minute, nil
		}
		nextMinute := current.minute + 1
		for _, offset := range offsets {
			next := current.position.Add(offset)
			id, valid := valley.positionID(next)
			if !valid || valley.blocked(next, nextMinute) {
				continue
			}
			key := (nextMinute%valley.period)*positions + id
			if !seen[key] {
				seen[key] = true
				queue.PushBack(state{position: next, minute: nextMinute})
			}
		}
	}
	return 0, fmt.Errorf("no route through the blizzard valley")
}

func day24p01(r io.Reader) (string, error) {
	valley, err := parseValley(r)
	if err != nil {
		return "", err
	}
	arrival, err := valley.travel(valley.start, valley.end, 0)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(arrival), nil
}

func day24p02(r io.Reader) (string, error) {
	valley, err := parseValley(r)
	if err != nil {
		return "", err
	}
	minute := 0
	from, to := valley.start, valley.end
	for range 3 {
		minute, err = valley.travel(from, to, minute)
		if err != nil {
			return "", err
		}
		from, to = to, from
	}
	return strconv.Itoa(minute), nil
}
