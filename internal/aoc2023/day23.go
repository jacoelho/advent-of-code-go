package aoc2023

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type trailSegment struct {
	toIdx    int
	distance int
}

type trailNetwork struct {
	segments [][]trailSegment
	startIdx int
	endIdx   int
}

type trailConditions uint8

const (
	slipperyTrail trailConditions = iota
	dryTrail
)

// trailSteps lists the moves in the order of the slope characters "^>v<".
var trailSteps = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type trailMap struct {
	rows       [][]byte
	start, end [2]int
}

func parseTrail(r io.Reader) (trailMap, error) {
	g, err := grid.ReadBytes(r, ".#^>v<")
	if err != nil {
		return trailMap{}, err
	}
	rows := g.Rows()
	startX, endX := bytes.IndexByte(rows[0], '.'), bytes.IndexByte(rows[len(rows)-1], '.')
	if startX < 0 || endX < 0 {
		return trailMap{}, fmt.Errorf("trail start or end not found")
	}
	return trailMap{rows: rows, start: [2]int{startX, 0}, end: [2]int{endX, len(rows) - 1}}, nil
}

func (t trailMap) open(x, y int) bool {
	return y >= 0 && y < len(t.rows) && x >= 0 && x < len(t.rows[y]) && t.rows[y][x] != '#'
}

// neighbors appends the cells reachable from x,y in one step.
func (t trailMap) neighbors(dst [][2]int, x, y int, conditions trailConditions) [][2]int {
	slope := bytes.IndexByte([]byte("^>v<"), t.rows[y][x])
	for direction, step := range trailSteps {
		if conditions == slipperyTrail && slope >= 0 && direction != slope {
			continue
		}
		if nx, ny := x+step[0], y+step[1]; t.open(nx, ny) {
			dst = append(dst, [2]int{nx, ny})
		}
	}
	return dst
}

func day23p01(r io.Reader) (string, error) {
	return longestHike(r, slipperyTrail)
}

func day23p02(r io.Reader) (string, error) {
	return longestHike(r, dryTrail)
}

func longestHike(r io.Reader, conditions trailConditions) (string, error) {
	trail, err := parseTrail(r)
	if err != nil {
		return "", err
	}
	length, err := findLongestHike(buildTrailNetwork(trail, conditions))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(length), nil
}

// buildTrailNetwork contracts every corridor between junctions into one
// segment. Junctions are the start, the end and every cell without exactly two
// open neighbours.
func buildTrailNetwork(trail trailMap, conditions trailConditions) trailNetwork {
	width := len(trail.rows[0])
	junction := make([]int, width*len(trail.rows))
	var junctions [][2]int
	var buffer [][2]int
	for y, row := range trail.rows {
		for x, tile := range row {
			junction[y*width+x] = -1
			cell := [2]int{x, y}
			if tile == '#' {
				continue
			}
			if cell == trail.start || cell == trail.end || len(trail.neighbors(buffer[:0], x, y, dryTrail)) != 2 {
				junction[y*width+x] = len(junctions)
				junctions = append(junctions, cell)
			}
		}
	}

	network := trailNetwork{
		segments: make([][]trailSegment, len(junctions)),
		startIdx: junction[trail.start[1]*width+trail.start[0]],
		endIdx:   junction[trail.end[1]*width+trail.end[0]],
	}
	for i, from := range junctions {
		for _, first := range trail.neighbors(nil, from[0], from[1], conditions) {
			previous, current := from, first
			for distance := 1; ; distance++ {
				if to := junction[current[1]*width+current[0]]; to >= 0 {
					network.segments[i] = append(network.segments[i], trailSegment{toIdx: to, distance: distance})
					break
				}
				// Corridor cells have two open neighbours; a slope may leave none forward.
				buffer = trail.neighbors(buffer[:0], current[0], current[1], conditions)
				next := -1
				for k, candidate := range buffer {
					if candidate != previous {
						next = k
						break
					}
				}
				if next < 0 {
					break
				}
				previous, current = current, buffer[next]
			}
		}
	}
	return network
}

func findLongestHike(network trailNetwork) (int, error) {
	// When a single junction leads to the end, a hike that reaches it and
	// turns elsewhere can never return, so that junction only exits to the end.
	lastJunction := -1
	for from, segments := range network.segments {
		for _, segment := range segments {
			if segment.toIdx == network.endIdx && from != network.endIdx {
				if lastJunction >= 0 && lastJunction != from {
					lastJunction = -2
				} else if lastJunction == -1 {
					lastJunction = from
				}
			}
		}
	}

	visited := make([]bool, len(network.segments))
	longest := -1
	var explore func(int, int)
	explore = func(current, distance int) {
		if current == network.endIdx {
			longest = max(longest, distance)
			return
		}
		visited[current] = true
		for _, segment := range network.segments[current] {
			if visited[segment.toIdx] || current == lastJunction && segment.toIdx != network.endIdx {
				continue
			}
			explore(segment.toIdx, distance+segment.distance)
		}
		visited[current] = false
	}
	explore(network.startIdx, 0)
	if longest < 0 {
		return 0, fmt.Errorf("no hiking route reaches the end")
	}
	return longest, nil
}
