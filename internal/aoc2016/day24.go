package aoc2016

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"maps"
	"math"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type cleaningFinish int

const (
	cleaningStop cleaningFinish = iota
	cleaningReturn
)

func cleaningRoute(r io.Reader, finish cleaningFinish) (string, error) {
	open := collections.NewSet[grid.Position2D[int]]()
	points := make(map[byte]grid.Position2D[int])
	s := bufio.NewScanner(r)
	width, height := 0, 0
	for s.Scan() {
		row := s.Text()
		if width == 0 {
			width = len(row)
		}
		if len(row) == 0 || len(row) != width {
			return "", fmt.Errorf("invalid maze row %d", height)
		}
		for x, tile := range []byte(row) {
			position := grid.Position2D[int]{X: x, Y: height}
			switch {
			case tile == '#':
				continue
			case tile == '.':
			case tile >= '0' && tile <= '9':
				if _, exists := points[tile]; exists {
					return "", fmt.Errorf("duplicate location %c", tile)
				}
				points[tile] = position
			default:
				return "", fmt.Errorf("unknown maze tile %q", tile)
			}
			open.Add(position)
		}
		height++
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if _, exists := points['0']; !exists {
		return "", fmt.Errorf("missing maze start")
	}

	labels := slices.Sorted(maps.Keys(points))
	n := len(labels)
	distances := make([][]int, n)
	neighbours := func(position grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) {
			for next := range grid.Neighbours4(position) {
				if open.Contains(next) && !yield(next) {
					return
				}
			}
		}
	}
	for i, label := range labels {
		from := search.BFSDistances(points[label], neighbours)
		distances[i] = make([]int, n)
		for j, target := range labels {
			distance, ok := from[points[target]]
			if !ok {
				return "", fmt.Errorf("location %c is unreachable from %c", target, label)
			}
			distances[i][j] = distance
		}
	}

	const infinity = math.MaxInt / 4
	paths := make([]int, (1<<n)*n)
	for i := range paths {
		paths[i] = infinity
	}
	paths[n] = 0
	for mask := 1; mask < 1<<n; mask++ {
		for last := range n {
			cost := paths[mask*n+last]
			if cost == infinity {
				continue
			}
			for next := range n {
				if mask&(1<<next) != 0 {
					continue
				}
				index := (mask|1<<next)*n + next
				paths[index] = min(paths[index], cost+distances[last][next])
			}
		}
	}

	best := infinity
	for last := range n {
		cost := paths[((1<<n)-1)*n+last]
		if finish == cleaningReturn {
			cost += distances[last][0]
		}
		best = min(best, cost)
	}
	return strconv.Itoa(best), nil
}

func day24p01(r io.Reader) (string, error) {
	return cleaningRoute(r, cleaningStop)
}

func day24p02(r io.Reader) (string, error) {
	return cleaningRoute(r, cleaningReturn)
}
