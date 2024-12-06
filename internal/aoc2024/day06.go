package aoc2024

import (
	"io"
	"slices"
	"strconv"
	"sync/atomic"

	"github.com/jacoelho/advent-of-code-go/internal/conc"

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
	position grid.Position2D[int],
	extraObstacles collections.Set[grid.Position2D[int]],
) (collections.Set[[2]grid.Position2D[int]], bool) {
	seen := collections.NewSet[[2]grid.Position2D[int]]()
	direction := grid.Position2D[int]{X: 0, Y: -1}

	keyFn := func(position, direction grid.Position2D[int]) [2]grid.Position2D[int] {
		return [2]grid.Position2D[int]{position, direction}
	}

	for g[position] != 0 && !seen.Contains(keyFn(position, direction)) {
		seen.Add(keyFn(position, direction))
		next := position.Add(direction)

		if g[next] == '#' || extraObstacles.Contains(next) {
			direction = direction.TurnRight()
		} else {
			position = next
		}
	}

	return seen, g[position] != 0
}

func day06p01(r io.Reader) (string, error) {
	m := aoc.Must(parseGuardMap(r))
	startPosition := guardPosition(m)

	positions, _ := followGuard(m, startPosition, collections.NewSet[grid.Position2D[int]]())
	count := xiter.Length(xiter.Unique(xiter.Map(func(in [2]grid.Position2D[int]) grid.Position2D[int] {
		return in[0]
	}, positions.Iter())))

	return strconv.Itoa(count), nil
}

func day06p02(r io.Reader) (string, error) {
	m := aoc.Must(parseGuardMap(r))
	startPosition := guardPosition(m)

	positions, _ := followGuard(m, startPosition, collections.NewSet[grid.Position2D[int]]())
	candidates := xiter.Unique(xiter.Map(func(in [2]grid.Position2D[int]) grid.Position2D[int] {
		return in[0]
	}, positions.Iter()))

	var wg conc.WaitGroup
	var count atomic.Int64
	for candidate := range candidates {
		if candidate == startPosition {
			continue
		}

		wg.Go(func() {
			if _, cycles := followGuard(m, startPosition, collections.NewSet(candidate)); cycles {
				count.Add(1)
			}
		})

	}
	wg.Wait()

	return strconv.FormatInt(count.Load(), 10), nil
}
