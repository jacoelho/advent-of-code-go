package aoc2022

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

// blizzardValley numbers interior cells row-major from zero, followed by the
// entrance and the exit, which blizzards never reach.
type blizzardValley struct {
	cells          []byte
	width, height  int
	entrance, exit int
	neighbours     [][]int
}

func parseValley(r io.Reader) (blizzardValley, error) {
	var valley blizzardValley
	input, err := puzzleInput(r)
	if err != nil {
		return valley, err
	}
	rows := strings.Split(input, "\n")
	if len(rows) < 3 {
		return valley, fmt.Errorf("valley requires top and bottom walls with an interior")
	}
	for _, row := range rows {
		if len(row) < 3 || len(row) != len(rows[0]) {
			return valley, fmt.Errorf("valley must be a rectangle at least three cells wide")
		}
	}
	valley.width, valley.height = len(rows[0])-2, len(rows)-2
	entrance, exit := -1, -1
	for y, row := range rows {
		for x := range len(row) {
			c := row[x]
			if x == 0 || x == len(row)-1 || y == 0 || y == len(rows)-1 {
				if c == '#' {
					continue
				}
				if c != '.' || x == 0 || x == len(row)-1 || y != 0 && y != len(rows)-1 {
					return valley, fmt.Errorf("invalid valley boundary at %d,%d", x, y)
				}
				if y == 0 {
					if entrance >= 0 {
						return valley, fmt.Errorf("valley has more than one entrance")
					}
					entrance = x - 1
				} else {
					if exit >= 0 {
						return valley, fmt.Errorf("valley has more than one exit")
					}
					exit = x - 1
				}
				continue
			}
			if c != '.' && c != '>' && c != '<' && c != '^' && c != 'v' {
				return valley, fmt.Errorf("invalid valley interior cell %q", c)
			}
			valley.cells = append(valley.cells, c)
		}
	}
	if entrance < 0 || exit < 0 {
		return valley, fmt.Errorf("valley requires one entrance and one exit")
	}
	w, h, n := valley.width, valley.height, len(valley.cells)
	valley.entrance, valley.exit = n, n+1
	valley.neighbours = make([][]int, n+2)
	for i := range n {
		x, y := i%w, i/w
		valley.neighbours[i] = append(valley.neighbours[i], i)
		if x > 0 {
			valley.neighbours[i] = append(valley.neighbours[i], i-1)
		}
		if x < w-1 {
			valley.neighbours[i] = append(valley.neighbours[i], i+1)
		}
		if y > 0 {
			valley.neighbours[i] = append(valley.neighbours[i], i-w)
		}
		if y < h-1 {
			valley.neighbours[i] = append(valley.neighbours[i], i+w)
		}
	}
	link := func(gate, cell int) {
		valley.neighbours[gate] = []int{gate, cell}
		valley.neighbours[cell] = append(valley.neighbours[cell], gate)
	}
	link(valley.entrance, entrance)
	link(valley.exit, (h-1)*w+exit)
	return valley, nil
}

func (valley blizzardValley) blocked(cell, minute int) bool {
	if cell >= len(valley.cells) {
		return false
	}
	w, h := valley.width, valley.height
	x, y := cell%w, cell/w
	// Inverse motion checks the original map without storing every hazard phase.
	return valley.cells[y*w+xmath.Modulo(x-minute, w)] == '>' ||
		valley.cells[y*w+xmath.Modulo(x+minute, w)] == '<' ||
		valley.cells[xmath.Modulo(y-minute, h)*w+x] == 'v' ||
		valley.cells[xmath.Modulo(y+minute, h)*w+x] == '^'
}

// travel advances the set of reachable cells one minute at a time. Waiting at
// the departure gate is always safe, so the set one blizzard period later is a
// superset of the current one; an unchanged set after a full period means the
// destination is unreachable.
func (valley blizzardValley) travel(from, to, minute int) (int, error) {
	period := xmath.LCM(valley.width, valley.height)
	reachable := make([]bool, len(valley.neighbours))
	next := make([]bool, len(reachable))
	snapshot := make([]bool, len(reachable))
	reachable[from] = true
	for elapsed := 0; !reachable[to]; elapsed++ {
		if elapsed%period == 0 {
			if elapsed > 0 && slices.Equal(reachable, snapshot) {
				return 0, fmt.Errorf("no route through the blizzard valley")
			}
			copy(snapshot, reachable)
		}
		minute++
		clear(next)
		for cell, ok := range reachable {
			if !ok {
				continue
			}
			for _, neighbour := range valley.neighbours[cell] {
				if !next[neighbour] && !valley.blocked(neighbour, minute) {
					next[neighbour] = true
				}
			}
		}
		reachable, next = next, reachable
	}
	return minute, nil
}

func day24p01(r io.Reader) (string, error) {
	valley, err := parseValley(r)
	if err != nil {
		return "", err
	}
	arrival, err := valley.travel(valley.entrance, valley.exit, 0)
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
	from, to := valley.entrance, valley.exit
	for range 3 {
		minute, err = valley.travel(from, to, minute)
		if err != nil {
			return "", err
		}
		from, to = to, from
	}
	return strconv.Itoa(minute), nil
}
