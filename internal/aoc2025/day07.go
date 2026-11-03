package aoc2025

import (
	"fmt"
	"io"
	"iter"
	"math/big"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/funcs"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func parseTachyonManifoldDiagram(r io.Reader) (
	grid.Grid2D[int, rune],
	grid.Position2D[int],
	error,
) {
	s := scanner.NewScanner(r, func(b []byte) (string, error) {
		return string(b), nil
	})
	g := make(grid.Grid2D[int, rune])
	var start grid.Position2D[int]
	found := false
	width, y := 0, 0
	for row := range s.Values() {
		if y == 0 {
			width = len(row)
		}
		if len(row) == 0 || len(row) != width {
			return nil, start, fmt.Errorf("manifold rows must be nonempty and have equal width")
		}
		if len(g)+len(row) > 100_000 || y == 4096 {
			return nil, start, fmt.Errorf("manifold exceeds 100000 cells or 4096 rows")
		}
		for x, tile := range row {
			if tile != '.' && tile != '^' && tile != 'S' {
				return nil, start, fmt.Errorf("invalid manifold tile %q", tile)
			}
			pos := grid.NewPosition2D(x, y)
			g[pos] = tile
			if tile == 'S' {
				if found {
					return nil, start, fmt.Errorf("manifold has multiple starting positions")
				}
				start, found = pos, true
			}
		}
		y++
	}
	if err := s.Err(); err != nil {
		return nil, grid.Position2D[int]{}, err
	}
	if !found {
		return nil, grid.Position2D[int]{}, fmt.Errorf("start position 'S' not found")
	}

	return g, start, nil
}

func traceTachyonBeam(
	g grid.Grid2D[int, rune],
	pos grid.Position2D[int],
) (grid.Position2D[int], bool) {
	for {
		next := grid.NewPosition2D(pos.X, pos.Y+1)
		v, exists := g[next]
		if !exists {
			return grid.Position2D[int]{}, false
		}
		switch v {
		case '.':
			pos = next
		case '^':
			return next, true
		}
	}
}

func day07p01(r io.Reader) (string, error) {
	diagram, start, err := parseTachyonManifoldDiagram(r)
	if err != nil {
		return "", err
	}

	first, found := traceTachyonBeam(diagram, start)
	if !found {
		return "0", nil
	}
	children := func(splitter grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) {
			for _, dx := range [...]int{-1, 1} {
				position := grid.NewPosition2D(splitter.X+dx, splitter.Y)
				if !diagram.Contains(position) {
					continue
				}
				if next, found := traceTachyonBeam(diagram, position); found && !yield(next) {
					return
				}
			}
		}
	}
	count := 0
	for range search.DFS(first, children) {
		count++
	}
	return strconv.Itoa(count), nil
}

func day07p02(r io.Reader) (string, error) {
	diagram, start, err := parseTachyonManifoldDiagram(r)
	if err != nil {
		return "", err
	}

	var countTimelines func(pos grid.Position2D[int]) *big.Int
	countTimelines = funcs.Memoize(func(pos grid.Position2D[int]) *big.Int {
		splitter, found := traceTachyonBeam(diagram, pos)
		if !found {
			return big.NewInt(1)
		}
		leftTimelines := countTimelines(grid.NewPosition2D(splitter.X-1, splitter.Y))
		rightTimelines := countTimelines(grid.NewPosition2D(splitter.X+1, splitter.Y))
		return new(big.Int).Add(leftTimelines, rightTimelines)
	})

	return countTimelines(start).String(), nil
}
