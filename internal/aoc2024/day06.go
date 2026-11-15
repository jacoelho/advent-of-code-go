package aoc2024

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

// Directions are indexed clockwise from north, so turning right is (d+1)%4.
var (
	guardDX = [4]int{0, 1, 0, -1}
	guardDY = [4]int{-1, 0, 1, 0}
)

type guardMap struct {
	width, height int
	wall          []bool
	start         int
}

// guardObstacle is a cell first entered from state before, encoded as
// 4*cell+direction; the walk up to before never touches cell.
type guardObstacle struct {
	cell, before int
}

func parseGuardMap(r io.Reader) (guardMap, error) {
	g, err := grid.ReadBytes(r, ".#^")
	if err != nil {
		return guardMap{}, err
	}
	start, err := g.Find('^')
	if err != nil {
		return guardMap{}, err
	}
	m := guardMap{width: g.Width(), height: g.Height(), start: g.Index(start), wall: make([]bool, len(g.Cells()))}
	for i, tile := range g.Cells() {
		m.wall[i] = tile == '#'
	}
	return m, nil
}

func (m guardMap) inside(x, y int) bool {
	return x >= 0 && x < m.width && y >= 0 && y < m.height
}

// guardStops[d][cell] is the coordinate along d's axis (y for north and south,
// x for east and west) of the last cell reached walking from cell in direction
// d before a wall or the map edge.
type guardStops [4][]int32

func (m guardMap) stops() guardStops {
	var result guardStops
	n := len(m.wall)
	for dir := range 4 {
		stops := make([]int32, n)
		// North and west neighbours precede a cell in row-major order; south
		// and east neighbours follow it, so those directions fill in reverse.
		reverse := dir == 1 || dir == 2
		for i := range n {
			cell := i
			if reverse {
				cell = n - 1 - i
			}
			x, y := cell%m.width, cell/m.width
			nx, ny := x+guardDX[dir], y+guardDY[dir]
			switch {
			case m.inside(nx, ny) && !m.wall[ny*m.width+nx]:
				stops[cell] = stops[ny*m.width+nx]
			case dir%2 == 0:
				stops[cell] = int32(y)
			default:
				stops[cell] = int32(x)
			}
		}
		result[dir] = stops
	}
	return result
}

// obstacles walks the unobstructed patrol, which may itself loop, and returns
// every cell it enters after the start, in order of first entry.
func (m guardMap) obstacles() []guardObstacle {
	seen := make([]bool, 4*len(m.wall))
	entered := make([]bool, len(m.wall))
	entered[m.start] = true
	var obstacles []guardObstacle
	x, y, dir := m.start%m.width, m.start/m.width, 0
	for {
		state := 4*(y*m.width+x) + dir
		if seen[state] {
			return obstacles
		}
		seen[state] = true
		nx, ny := x+guardDX[dir], y+guardDY[dir]
		if !m.inside(nx, ny) {
			return obstacles
		}
		next := ny*m.width + nx
		if m.wall[next] {
			dir = (dir + 1) % 4
			continue
		}
		if !entered[next] {
			entered[next] = true
			obstacles = append(obstacles, guardObstacle{cell: next, before: state})
		}
		x, y = nx, ny
	}
}

// loopsWith reports whether the patrol loops once obstacle is placed. It
// resumes from obstacle.before and jumps between turns, recording each turn in
// marks under generation.
func (m guardMap) loopsWith(stops guardStops, obstacle guardObstacle, marks []int32, generation int32) bool {
	ox, oy := obstacle.cell%m.width, obstacle.cell/m.width
	cell, dir := obstacle.before/4, obstacle.before%4
	x, y := cell%m.width, cell/m.width
	edges := [4]int{0, m.width - 1, m.height - 1, 0}
	for {
		vertical := dir%2 == 0
		along, obstacleAlong, sameLine := x, ox, oy == y
		if vertical {
			along, obstacleAlong, sameLine = y, oy, ox == x
		}
		sign := guardDX[dir] + guardDY[dir]
		end := int(stops[dir][y*m.width+x])
		if sameLine && (obstacleAlong-along)*sign > 0 && (end-obstacleAlong)*sign >= 0 {
			end = obstacleAlong - sign
		} else if end == edges[dir] {
			return false
		}
		if vertical {
			y = end
		} else {
			x = end
		}
		state := 4*(y*m.width+x) + dir
		if marks[state] == generation {
			return true
		}
		marks[state] = generation
		dir = (dir + 1) % 4
	}
}

func day06p01(r io.Reader) (string, error) {
	m, err := parseGuardMap(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(len(m.obstacles()) + 1), nil
}

func day06p02(r io.Reader) (string, error) {
	m, err := parseGuardMap(r)
	if err != nil {
		return "", err
	}
	stops := m.stops()
	marks := make([]int32, 4*len(m.wall))
	var count int
	for i, obstacle := range m.obstacles() {
		if m.loopsWith(stops, obstacle, marks, int32(i+1)) {
			count++
		}
	}
	return strconv.Itoa(count), nil
}
