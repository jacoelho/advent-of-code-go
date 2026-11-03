package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type burrow [27]byte

func roomCell(room, level int) int { return 11 + 4*room + level }

func parseBurrow(r io.Reader) (burrow, int, error) {
	var state burrow
	var lines []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		lines = append(lines, s.Text())
		if len(lines) > 7 {
			return state, 0, fmt.Errorf("burrow supports room depths two and four")
		}
	}
	if err := s.Err(); err != nil {
		return state, 0, err
	}
	depth := len(lines) - 3
	if depth != 2 && depth != 4 {
		return state, 0, fmt.Errorf("burrow supports room depths two and four")
	}
	if lines[0] != "#############" || lines[len(lines)-1] != "  #########" {
		return state, 0, fmt.Errorf("invalid burrow boundary")
	}
	counts := [5]int{}
	readCell := func(c byte, index int) error {
		if c == '.' {
			return nil
		}
		if c < 'A' || c > 'D' {
			return fmt.Errorf("invalid amphipod %q", c)
		}
		state[index] = c - 'A' + 1
		counts[state[index]]++
		return nil
	}
	hall := lines[1]
	if len(hall) != 13 || hall[0] != '#' || hall[12] != '#' {
		return state, 0, fmt.Errorf("invalid burrow hallway")
	}
	for x := range 11 {
		if err := readCell(hall[x+1], x); err != nil {
			return state, 0, err
		}
	}
	for level := range depth {
		line := lines[level+2]
		template := "  #.#.#.#.#"
		if level == 0 {
			template = "###.#.#.#.###"
		}
		if len(line) != len(template) {
			return state, 0, fmt.Errorf("invalid room row %d", level+1)
		}
		for i := range line {
			if template[i] != '.' && line[i] != template[i] {
				return state, 0, fmt.Errorf("invalid room wall at row %d", level+1)
			}
		}
		for room := range 4 {
			if err := readCell(line[3+2*room], roomCell(room, level)); err != nil {
				return state, 0, err
			}
		}
	}
	for room := range 4 {
		if state[2+2*room] != 0 {
			return state, 0, fmt.Errorf("amphipod cannot stop at a doorway")
		}
		occupied := false
		for level := range depth {
			if state[roomCell(room, level)] != 0 {
				occupied = true
			} else if occupied {
				return state, 0, fmt.Errorf("room has an empty cell below an amphipod")
			}
		}
		if counts[room+1] != depth {
			return state, 0, fmt.Errorf("expected %d amphipods of type %c", depth, 'A'+room)
		}
	}
	return state, depth, nil
}

func (state burrow) roomDestination(room, depth int) int {
	destination := -1
	for level := range depth {
		value := state[roomCell(room, level)]
		if value == 0 {
			destination = level
		} else if int(value) != room+1 {
			return -1
		}
	}
	return destination
}

func (state burrow) hallwayClear(from, to int) bool {
	step := 1
	if from > to {
		step = -1
	}
	for x := from + step; x != to+step; x += step {
		if state[x] != 0 {
			return false
		}
	}
	return true
}

func (state burrow) settled(room, level, depth int) bool {
	for i := level; i < depth; i++ {
		if int(state[roomCell(room, i)]) != room+1 {
			return false
		}
	}
	return true
}

func amphipodEnergy(kind byte) int {
	return [...]int{0, 1, 10, 100, 1000}[kind]
}

func (state burrow) remainingEnergy(depth int) int {
	total := 0
	for x := range 11 {
		if kind := state[x]; kind != 0 {
			total += (xmath.Abs(x-2*int(kind)) + 1) * amphipodEnergy(kind)
		}
	}
	for room := range 4 {
		for level := range depth {
			kind := state[roomCell(room, level)]
			if kind == 0 || state.settled(room, level, depth) {
				continue
			}
			// Ignore blockers and deeper destination cells to keep this a lower bound.
			total += (level + 2 + xmath.Abs(2*(room+1)-2*int(kind))) * amphipodEnergy(kind)
		}
	}
	return total
}

func (state burrow) moves(depth int, visit func(burrow, int)) {
	move := func(from, to, steps int) {
		next := state
		next[to], next[from] = state[from], 0
		visit(next, steps*amphipodEnergy(state[from]))
	}

	for x := range 11 {
		kind := state[x]
		if kind == 0 {
			continue
		}
		room := int(kind) - 1
		level := state.roomDestination(room, depth)
		if level >= 0 && state.hallwayClear(x, 2+2*room) {
			move(x, roomCell(room, level), xmath.Abs(x-(2+2*room))+level+1)
		}
	}

	for room := range 4 {
		level := 0
		for level < depth && state[roomCell(room, level)] == 0 {
			level++
		}
		if level == depth || state.settled(room, level, depth) {
			continue
		}
		from := roomCell(room, level)
		entrance := 2 + 2*room
		destination := int(state[from]) - 1
		if destination != room {
			landing := state.roomDestination(destination, depth)
			if landing >= 0 && state.hallwayClear(entrance, 2+2*destination) {
				move(from, roomCell(destination, landing), level+landing+2+2*xmath.Abs(room-destination))
			}
		}
		for _, x := range [...]int{0, 1, 3, 5, 7, 9, 10} {
			if state.hallwayClear(entrance, x) {
				move(from, x, level+1+xmath.Abs(entrance-x))
			}
		}
	}
}

func organizeBurrow(start burrow, depth int) (int, error) {
	type entry struct {
		state       burrow
		cost, score int
	}
	queue := collections.NewHeap(func(a, b entry) bool { return a.score < b.score })
	queue.Push(entry{state: start, score: start.remainingEnergy(depth)})
	best := map[burrow]int{start: 0}
	const limit = 2_000_000
	expanded := 0
	for current := range queue.PopSeq() {
		if best[current.state] != current.cost {
			continue
		}
		if current.state.remainingEnergy(depth) == 0 {
			return current.cost, nil
		}
		expanded++
		if expanded > limit {
			return 0, fmt.Errorf("burrow search exceeds %d expanded states", limit)
		}
		current.state.moves(depth, func(next burrow, energy int) {
			cost := current.cost + energy
			if previous, exists := best[next]; exists && previous <= cost {
				return
			}
			best[next] = cost
			queue.Push(entry{state: next, cost: cost, score: cost + next.remainingEnergy(depth)})
		})
		if len(best) > limit || queue.Len() > limit {
			return 0, fmt.Errorf("burrow search exceeds %d stored states", limit)
		}
	}
	return 0, fmt.Errorf("burrow cannot be organized")
}

func day23p01(r io.Reader) (string, error) {
	state, depth, err := parseBurrow(r)
	if err != nil {
		return "", err
	}
	energy, err := organizeBurrow(state, depth)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(energy), nil
}

func day23p02(r io.Reader) (string, error) {
	state, depth, err := parseBurrow(r)
	if err != nil {
		return "", err
	}
	if depth != 2 {
		return "", fmt.Errorf("unfolding requires a burrow with two room rows")
	}
	for x := range 11 {
		if state[x] != 0 {
			return "", fmt.Errorf("unfolding requires all amphipods in rooms")
		}
	}
	for room, inserted := range [4][2]byte{{4, 4}, {3, 2}, {2, 1}, {1, 3}} {
		state[roomCell(room, 3)] = state[roomCell(room, 1)]
		state[roomCell(room, 1)] = inserted[0]
		state[roomCell(room, 2)] = inserted[1]
	}
	energy, err := organizeBurrow(state, 4)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(energy), nil
}
