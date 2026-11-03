package aoc2015

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func deliverPresents(r io.Reader, couriers int) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	positions := make([]grid.Position2D[int], couriers)
	visited := collections.NewSet(grid.Position2D[int]{})
	for i, direction := range bytes.TrimSpace(input) {
		position := &positions[i%couriers]
		switch direction {
		case '^':
			position.Y++
		case 'v':
			position.Y--
		case '>':
			position.X++
		case '<':
			position.X--
		default:
			return "", fmt.Errorf("invalid delivery direction %q", direction)
		}
		visited.Add(*position)
	}
	return strconv.Itoa(visited.Len()), nil
}

func day03p01(r io.Reader) (string, error) {
	return deliverPresents(r, 1)
}

func day03p02(r io.Reader) (string, error) {
	return deliverPresents(r, 2)
}
