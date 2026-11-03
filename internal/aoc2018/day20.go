package aoc2018

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func facilityDistances(r io.Reader) (map[grid.Position2D[int]]int, error) {
	input, err := io.ReadAll(io.LimitReader(r, 1<<20+1))
	if err != nil {
		return nil, err
	}
	input = bytes.TrimSpace(input)
	if len(input) > 1<<20 || len(input) < 2 || input[0] != '^' || input[len(input)-1] != '$' {
		return nil, fmt.Errorf("expected a route expression enclosed by ^ and $, at most 1 MiB")
	}

	type branch struct {
		start, ends collections.Set[grid.Position2D[int]]
	}
	var stack []branch
	start := grid.Position2D[int]{}
	positions := collections.NewSet(start)
	doors := map[grid.Position2D[int]]uint8{start: 0}
	directions := [...]grid.Position2D[int]{{Y: -1}, {X: 1}, {Y: 1}, {X: -1}}
	work := 0
	merge := func(destination, source collections.Set[grid.Position2D[int]]) error {
		for pos := range source {
			work++
			if work > 10_000_000 {
				return fmt.Errorf("route expansion exceeded 10000000 operations")
			}
			destination.Add(pos)
		}
		return nil
	}

	for _, token := range input[1 : len(input)-1] {
		switch token {
		case '(':
			stack = append(stack, branch{start: positions, ends: collections.NewSet[grid.Position2D[int]]()})
		case '|', ')':
			if len(stack) == 0 {
				return nil, fmt.Errorf("unmatched route branch %q", token)
			}
			current := &stack[len(stack)-1]
			if err := merge(current.ends, positions); err != nil {
				return nil, err
			}
			if token == '|' {
				positions = current.start
			} else {
				positions = current.ends
				stack = stack[:len(stack)-1]
			}
		default:
			direction := bytes.IndexByte([]byte("NESW"), token)
			if direction < 0 {
				return nil, fmt.Errorf("invalid route token %q", token)
			}
			next := collections.NewSet[grid.Position2D[int]]()
			for pos := range positions {
				work++
				if work > 10_000_000 {
					return nil, fmt.Errorf("route expansion exceeded 10000000 operations")
				}
				target := pos.Add(directions[direction])
				doors[pos] |= 1 << direction
				doors[target] |= 1 << ((direction + 2) % 4)
				next.Add(target)
			}
			if len(doors) > 1_000_000 {
				return nil, fmt.Errorf("facility exceeds 1000000 rooms")
			}
			positions = next
		}
	}
	if len(stack) != 0 {
		return nil, fmt.Errorf("unclosed route branch")
	}

	neighbours := func(pos grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) {
			for direction, offset := range directions {
				if doors[pos]&(1<<direction) != 0 && !yield(pos.Add(offset)) {
					return
				}
			}
		}
	}

	return search.BFSDistances(start, neighbours), nil
}

func day20p01(r io.Reader) (string, error) {
	distances, err := facilityDistances(r)
	if err != nil {
		return "", err
	}

	furthest := 0
	for _, distance := range distances {
		furthest = max(furthest, distance)
	}

	return strconv.Itoa(furthest), nil
}

func day20p02(r io.Reader) (string, error) {
	distances, err := facilityDistances(r)
	if err != nil {
		return "", err
	}

	count := 0
	for _, distance := range distances {
		if distance >= 1000 {
			count++
		}
	}

	return strconv.Itoa(count), nil
}
