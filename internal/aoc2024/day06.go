package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
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

func guardPosition(g grid.Grid2D[int, rune]) (grid.Position2D[int], error) {
	var position grid.Position2D[int]
	found := false
	for p, tile := range g {
		if tile != '^' {
			continue
		}
		if found {
			return position, fmt.Errorf("guard map has multiple guards")
		}
		position, found = p, true
	}
	if !found {
		return position, fmt.Errorf("guard not found")
	}
	return position, nil
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
	m, err := parseGuardMap(r)
	if err != nil {
		return "", err
	}
	startPosition, err := guardPosition(m)
	if err != nil {
		return "", err
	}

	positions, _ := followGuard(m, startPosition, collections.NewSet[grid.Position2D[int]]())
	count := xiter.From(xiter.Unique(xiter.From(positions.Iter()).Map(func(in [2]grid.Position2D[int]) grid.Position2D[int] {
		return in[0]
	}).All())).Count()

	return strconv.Itoa(count), nil
}

func day06p02(r io.Reader) (string, error) {
	m, err := parseGuardMap(r)
	if err != nil {
		return "", err
	}
	startPosition, err := guardPosition(m)
	if err != nil {
		return "", err
	}

	positions, _ := followGuard(m, startPosition, collections.NewSet[grid.Position2D[int]]())
	candidates := xiter.Unique(xiter.From(positions.Iter()).Map(func(in [2]grid.Position2D[int]) grid.Position2D[int] {
		return in[0]
	}).All())

	var count int
	for candidate := range candidates {
		if candidate == startPosition {
			continue
		}
		if _, cycles := followGuard(m, startPosition, collections.NewSet(candidate)); cycles {
			count++
		}
	}

	return strconv.Itoa(count), nil
}
