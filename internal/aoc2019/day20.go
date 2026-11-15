package aoc2019

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const (
	donutStart = 0
	donutEnd   = 1
)

type donutEdge struct {
	to       int
	distance int
}

// donutNode is a labelled passage cell: the start, the end or a portal
// endpoint. Portal endpoints link to their partner on the other rim.
type donutNode struct {
	label   string
	edges   []donutEdge
	partner int
	outer   bool
}

// donutMaze is the graph of labelled passage cells; node donutStart is AA and
// node donutEnd is ZZ.
type donutMaze struct {
	nodes []donutNode
}

type donutGrid struct {
	lines [][]byte
}

func (g donutGrid) at(x, y int) byte {
	if y < 0 || y >= len(g.lines) || x < 0 || x >= len(g.lines[y]) {
		return ' '
	}
	return g.lines[y][x]
}

func isLabel(c byte) bool { return c >= 'A' && c <= 'Z' }

func parseDonutMaze(r io.Reader) (donutMaze, error) {
	s := scanner.NewScanner(r, func(line []byte) ([]byte, error) {
		for x, c := range line {
			if c != ' ' && c != '#' && c != '.' && !isLabel(c) {
				return nil, fmt.Errorf("invalid maze cell %q at column %d", c, x)
			}
		}
		return bytes.Clone(line), nil
	})
	var g donutGrid
	for line := range s.Values() {
		g.lines = append(g.lines, line)
	}
	if err := s.Err(); err != nil {
		return donutMaze{}, err
	}

	minX, minY, maxX, maxY := -1, -1, -1, -1
	for y, line := range g.lines {
		for x, c := range line {
			if c == '#' || c == '.' {
				if minX < 0 || x < minX {
					minX = x
				}
				if minY < 0 {
					minY = y
				}
				maxX, maxY = max(maxX, x), y
			}
		}
	}
	if minX < 0 {
		return donutMaze{}, fmt.Errorf("maze has no walls or passages")
	}

	type portal struct {
		label string
		x, y  int
		outer bool
	}
	var portals []portal
	letters := 0
	for y, line := range g.lines {
		for x, c := range line {
			if !isLabel(c) {
				continue
			}
			letters++
			// A label reads left to right or top to bottom; its passage sits just
			// before the first letter or just after the second.
			for _, d := range [][2]int{{1, 0}, {0, 1}} {
				second := g.at(x+d[0], y+d[1])
				if !isLabel(second) {
					continue
				}
				label := string([]byte{c, second})
				var found []portal
				for _, p := range [][2]int{{x - d[0], y - d[1]}, {x + 2*d[0], y + 2*d[1]}} {
					if g.at(p[0], p[1]) == '.' {
						outer := p[0] == minX || p[0] == maxX || p[1] == minY || p[1] == maxY
						found = append(found, portal{label: label, x: p[0], y: p[1], outer: outer})
					}
				}
				if len(found) != 1 {
					return donutMaze{}, fmt.Errorf("label %s at (%d,%d) needs exactly one adjacent passage", label, x, y)
				}
				portals = append(portals, found[0])
			}
		}
	}
	if letters != 2*len(portals) {
		return donutMaze{}, fmt.Errorf("maze has %d label letters, want two per label", letters)
	}

	maze := donutMaze{nodes: []donutNode{{label: "AA", partner: -1}, {label: "ZZ", partner: -1}}}
	cells := map[[2]int]int{}
	byLabel := map[string]int{}
	for _, p := range portals {
		if _, taken := cells[[2]int{p.x, p.y}]; taken {
			return donutMaze{}, fmt.Errorf("passage (%d,%d) has more than one label", p.x, p.y)
		}
		index := len(maze.nodes)
		switch p.label {
		case "AA":
			index = donutStart
		case "ZZ":
			index = donutEnd
		}
		if index < len(maze.nodes) {
			if _, dup := byLabel[p.label]; dup {
				return donutMaze{}, fmt.Errorf("label %s appears more than once", p.label)
			}
			byLabel[p.label] = index
			cells[[2]int{p.x, p.y}] = index
			continue
		}
		maze.nodes = append(maze.nodes, donutNode{label: p.label, partner: -1, outer: p.outer})
		cells[[2]int{p.x, p.y}] = index
		other, seen := byLabel[p.label]
		if !seen {
			byLabel[p.label] = index
			continue
		}
		if maze.nodes[other].partner >= 0 {
			return donutMaze{}, fmt.Errorf("portal %s appears more than twice", p.label)
		}
		if maze.nodes[other].outer == p.outer {
			return donutMaze{}, fmt.Errorf("portal %s needs one outer and one inner end", p.label)
		}
		maze.nodes[other].partner = index
		maze.nodes[index].partner = other
	}
	for _, label := range []string{"AA", "ZZ"} {
		if _, ok := byLabel[label]; !ok {
			return donutMaze{}, fmt.Errorf("maze has no %s", label)
		}
	}
	for _, node := range maze.nodes[2:] {
		if node.partner < 0 {
			return donutMaze{}, fmt.Errorf("portal %s has only one end", node.label)
		}
	}

	positions := make([][2]int, len(maze.nodes))
	for cell, index := range cells {
		positions[index] = cell
	}
	for source, start := range positions {
		distance := map[[2]int]int{start: 0}
		queue := [][2]int{start}
		for head := 0; head < len(queue); head++ {
			at := queue[head]
			for _, d := range [][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
				next := [2]int{at[0] + d[0], at[1] + d[1]}
				if _, seen := distance[next]; seen || g.at(next[0], next[1]) != '.' {
					continue
				}
				distance[next] = distance[at] + 1
				queue = append(queue, next)
				if target, ok := cells[next]; ok {
					maze.nodes[source].edges = append(maze.nodes[source].edges, donutEdge{to: target, distance: distance[next]})
				}
			}
		}
	}
	return maze, nil
}

// shortestPath returns the fewest steps from AA to ZZ. In a recursive maze an
// inner portal leads one level deeper and an outer portal one level up; outer
// portals are walls at the outermost level, and ZZ counts only there.
func (m donutMaze) shortestPath(recursive bool) (int, error) {
	type state struct{ node, level int }
	type visit struct {
		state
		distance int
	}
	dist := map[state]int{{donutStart, 0}: 0}
	queue := collections.NewHeap(func(a, b visit) bool { return a.distance < b.distance })
	queue.Push(visit{state: state{donutStart, 0}})
	relax := func(next state, cost int) {
		if known, ok := dist[next]; !ok || cost < known {
			dist[next] = cost
			queue.Push(visit{state: next, distance: cost})
		}
	}
	for current := range queue.PopSeq() {
		if current.distance != dist[current.state] {
			continue
		}
		if current.state == (state{donutEnd, 0}) {
			return current.distance, nil
		}
		node := m.nodes[current.node]
		for _, edge := range node.edges {
			relax(state{edge.to, current.level}, current.distance+edge.distance)
		}
		if node.partner < 0 {
			continue
		}
		level := current.level
		if recursive {
			if node.outer {
				level--
			} else {
				level++
			}
		}
		if level >= 0 {
			relax(state{node.partner, level}, current.distance+1)
		}
	}
	return 0, fmt.Errorf("ZZ is unreachable from AA")
}

func day20p01(r io.Reader) (string, error) {
	maze, err := parseDonutMaze(r)
	if err != nil {
		return "", err
	}
	steps, err := maze.shortestPath(false)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(steps), nil
}

func day20p02(r io.Reader) (string, error) {
	maze, err := parseDonutMaze(r)
	if err != nil {
		return "", err
	}
	// The recursive search never ends when ZZ is unreachable even through
	// portals at a single level, so rule that out first.
	if _, err := maze.shortestPath(false); err != nil {
		return "", err
	}
	steps, err := maze.shortestPath(true)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(steps), nil
}
