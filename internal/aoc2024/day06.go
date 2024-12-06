package aoc2024

import (
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

func parseGuardMap(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		res := make([]rune, len(b))
		for i, v := range b {
			res[i] = rune(v)
		}
		return res, nil
	})
	return grid.NewGrid2D[int, rune](slices.Collect(s.Values())), s.Err()
}

func guardPosition(g grid.Grid2D[int, rune]) grid.Position2D[int] {
	for p, v := range g {
		if v != '#' && v != '.' {
			return p
		}
	}
	panic("guard not found")
}

func followGuard(
	g grid.Grid2D[int, rune],
	start grid.Position2D[int],
) (iter.Seq[grid.Position2D[int]], bool) {
	visited := collections.NewSet[[2]grid.Position2D[int]]()
	location := start
	direction := grid.Position2D[int]{X: 0, Y: -1}

	for g[location] != 0 && !visited.Contains([2]grid.Position2D[int]{location, direction}) {
		visited.Add([2]grid.Position2D[int]{location, direction})
		next := location.Add(direction)

		if g[next] == '#' {
			direction = direction.TurnRight()
		} else {
			location = next
		}
	}

	positionsVisited := xiter.Unique(xiter.Map(func(in [2]grid.Position2D[int]) grid.Position2D[int] {
		return in[0]
	}, visited.Iter()))

	return positionsVisited, g[location] != 0
}

func day06p01(r io.Reader) (string, error) {
	m := aoc.Must(parseGuardMap(r))
	startPosition := guardPosition(m)

	positions, _ := followGuard(m, startPosition)
	count := xiter.Length(positions)

	return strconv.Itoa(count), nil
}

func day06p02(r io.Reader) (string, error) {
	m := aoc.Must(parseGuardMap(r))
	startPosition := guardPosition(m)

	candidates, _ := followGuard(m, startPosition)

	total := 0
	for candidate := range candidates {
		if candidate == startPosition {
			continue
		}
		m[candidate] = '#'

		_, cycle := followGuard(m, startPosition)
		if cycle {
			total++
		}
		m[candidate] = '.'
	}

	return strconv.Itoa(total), nil
}
