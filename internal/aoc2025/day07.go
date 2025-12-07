package aoc2025

import (
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/funcs"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmaps"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

var (
	down  = grid.Position2D[int]{X: 0, Y: 1}
	left  = grid.Position2D[int]{X: -1, Y: 0}
	right = grid.Position2D[int]{X: 1, Y: 0}
)

func parseTachyonManifoldDiagram(r io.Reader) (
	grid.Grid2D[int, rune],
	grid.Position2D[int],
	int,
	error,
) {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		return xslices.Map(func(v byte) rune { return rune(v) }, b), nil
	})

	g := grid.NewGrid2D[int](slices.Collect(s.Values()))
	if err := s.Err(); err != nil {
		return nil, grid.Position2D[int]{}, 0, err
	}
	_, _, _, maxY := g.Dimensions()

	pair, found := xmaps.Find(g, func(_ grid.Position2D[int], v rune) bool {
		return v == 'S'
	})
	if !found {
		panic("start position 'S' not found")
	}

	return g, pair.K, maxY, nil
}

// traceTachyonBeam traces a beam downward until hitting a splitter or exiting.
func traceTachyonBeam(
	g grid.Grid2D[int, rune],
	maxY int,
	pos grid.Position2D[int],
) (grid.Position2D[int], bool) {
	for {
		next := pos.Add(down)
		if next.Y > maxY {
			return grid.Position2D[int]{}, false
		}
		v, exists := g[next]
		if !exists {
			return grid.Position2D[int]{}, false
		}
		switch v {
		case '.', 'S':
			pos = next
		case '^':
			return next, true
		}
	}
}

func day07p01(r io.Reader) (string, error) {
	diagram, start, maxY, err := parseTachyonManifoldDiagram(r)
	if err != nil {
		return "", err
	}

	seen := collections.NewSet[grid.Position2D[int]]()

	var countBeamSplits func(pos grid.Position2D[int]) int
	countBeamSplits = func(pos grid.Position2D[int]) int {
		splitter, found := traceTachyonBeam(diagram, maxY, pos)
		if !found {
			return 0
		}
		if seen.Contains(splitter) {
			return 0
		}
		seen.Add(splitter)

		count := 1
		if leftPos := splitter.Add(left); diagram.Contains(leftPos) {
			count += countBeamSplits(leftPos)
		}
		if rightPos := splitter.Add(right); diagram.Contains(rightPos) {
			count += countBeamSplits(rightPos)
		}
		return count
	}

	return strconv.Itoa(countBeamSplits(start)), nil
}

func day07p02(r io.Reader) (string, error) {
	diagram, start, maxY, err := parseTachyonManifoldDiagram(r)
	if err != nil {
		return "", err
	}

	var countTimelines func(pos grid.Position2D[int]) int
	countTimelines = funcs.Memoize(func(pos grid.Position2D[int]) int {
		splitter, found := traceTachyonBeam(diagram, maxY, pos)
		if !found {
			return 1
		}

		leftTimelines, rightTimelines := 0, 0
		if leftPos := splitter.Add(left); diagram.Contains(leftPos) {
			leftTimelines = countTimelines(leftPos)
		}
		if rightPos := splitter.Add(right); diagram.Contains(rightPos) {
			rightTimelines = countTimelines(rightPos)
		}
		if leftTimelines == 0 && rightTimelines == 0 {
			return 1
		}
		return leftTimelines + rightTimelines
	})

	return strconv.Itoa(countTimelines(start)), nil
}
