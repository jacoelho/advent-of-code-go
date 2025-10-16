package aoc2023

import (
	"fmt"
	"io"
	"iter"
	"maps"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parseGarden(r io.Reader) (grid.Grid2D[int, rune], grid.Position2D[int], error) {
	g := make(grid.Grid2D[int, rune])
	var start grid.Position2D[int]
	var foundStart bool

	lineScanner := scanner.NewScanner(r, func(bytes []byte) (string, error) {
		return string(bytes), nil
	})

	for y, line := range xiter.Enumerate(lineScanner.Values()) {
		for x, ch := range line {
			pos := grid.NewPosition2D(x, y)
			g[pos] = ch
			if ch == 'S' {
				start = pos
				foundStart = true
			}
		}
	}

	if err := lineScanner.Err(); err != nil {
		return nil, grid.Position2D[int]{}, err
	}

	if !foundStart {
		return nil, grid.Position2D[int]{}, fmt.Errorf("starting position 'S' not found")
	}

	return g, start, nil
}

func isGardenPlot(g grid.Grid2D[int, rune], pos grid.Position2D[int]) bool {
	ch, exists := g[pos]
	return exists && (ch == '.' || ch == 'S')
}

func gardenNeighbours(g grid.Grid2D[int, rune]) func(grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
	return func(pos grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return xiter.Filter(func(next grid.Position2D[int]) bool {
			return isGardenPlot(g, next)
		}, grid.Neighbours4(pos))
	}
}

func countReachableInSteps(distances map[grid.Position2D[int]]int, targetSteps int) int {
	targetParity := targetSteps % 2
	return xiter.CountBy(func(d int) bool {
		return d <= targetSteps && d%2 == targetParity
	}, maps.Values(distances))
}

func day21(r io.Reader, steps int) (string, error) {
	g, start, err := parseGarden(r)
	if err != nil {
		return "", err
	}

	distances := search.BFSDistances(start, gardenNeighbours(g))
	count := countReachableInSteps(distances, steps)

	return strconv.Itoa(count), nil
}

func day21p01(r io.Reader) (string, error) {
	return day21(r, 64)
}

func day21p02(r io.Reader) (string, error) {
	g, start, err := parseGarden(r)
	if err != nil {
		return "", err
	}

	_, maxX, _, _ := g.Dimensions()
	width := maxX + 1
	gridSize := width // 131

	// count all positions by running BFS to completion
	steps := search.BFSDistances(start, gardenNeighbours(g))

	// count blocks (all reachable positions)
	oddBlock := 0
	evenBlock := 0
	oddCorners := 0
	evenCorners := 0

	for _, dist := range steps {
		if dist%2 == 1 {
			oddBlock++
			if dist > 65 {
				oddCorners++
			}
		} else {
			evenBlock++
			if dist > 65 {
				evenCorners++
			}
		}
	}

	stepCount := 26501365
	n := (stepCount - (gridSize / 2)) / gridSize

	even := n * n
	odd := (n + 1) * (n + 1)

	result := (odd * oddBlock) + (even * evenBlock) - ((n + 1) * oddCorners) + (n * evenCorners)

	return strconv.Itoa(result), nil
}
