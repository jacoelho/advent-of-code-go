package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func convertDirection(r rune) (grid.Position2D[int], error) {
	switch r {
	case '<':
		return grid.Position2D[int]{X: -1, Y: 0}, nil
	case '>':
		return grid.Position2D[int]{X: 1, Y: 0}, nil
	case 'v':
		return grid.Position2D[int]{X: 0, Y: 1}, nil
	case '^':
		return grid.Position2D[int]{X: 0, Y: -1}, nil
	default:
		return grid.Position2D[int]{}, fmt.Errorf("invalid warehouse movement %q", r)
	}
}

func warehouseIdentityTile(tile byte) []byte {
	return []byte{tile}
}

func warehouseScaleUpTile(tile byte) []byte {
	switch tile {
	case 'O':
		return []byte{'[', ']'}
	case '@':
		return []byte{'@', '.'}
	}
	return []byte{tile, tile}
}

type day15Warehouse struct {
	tiles [][]byte
	robot grid.Position2D[int]
}

func parseWarehouse(r io.Reader, transform func(byte) []byte) (day15Warehouse, []grid.Position2D[int], error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return day15Warehouse{}, nil, err
	}
	if len(blocks) < 2 {
		return day15Warehouse{}, nil, fmt.Errorf("warehouse needs a map and movement section")
	}

	rows := blocks[0]
	if len(rows) < 3 || len(rows[0]) < 3 {
		return day15Warehouse{}, nil, fmt.Errorf("warehouse needs an enclosed map")
	}
	width := len(rows[0])
	var w day15Warehouse
	robots := 0
	w.tiles = make([][]byte, len(rows))
	for y, row := range rows {
		if len(row) != width {
			return day15Warehouse{}, nil, fmt.Errorf("warehouse row %d has width %d, want %d", y+1, len(row), width)
		}
		for x := range len(row) {
			tile := row[x]
			if tile != '#' && tile != '.' && tile != 'O' && tile != '@' {
				return day15Warehouse{}, nil, fmt.Errorf("warehouse row %d has invalid tile %q", y+1, tile)
			}
			if (x == 0 || x == width-1 || y == 0 || y == len(rows)-1) && tile != '#' {
				return day15Warehouse{}, nil, fmt.Errorf("warehouse boundary must be walls")
			}
			if tile == '@' {
				robots++
				w.robot = grid.Position2D[int]{X: len(w.tiles[y]), Y: y}
			}
			w.tiles[y] = append(w.tiles[y], transform(tile)...)
		}
	}
	if robots != 1 {
		return day15Warehouse{}, nil, fmt.Errorf("warehouse has %d robots, want 1", robots)
	}

	var movements []grid.Position2D[int]
	for _, line := range slices.Concat(blocks[1:]...) {
		for _, ch := range line {
			direction, err := convertDirection(ch)
			if err != nil {
				return day15Warehouse{}, nil, err
			}
			movements = append(movements, direction)
		}
	}
	return w, movements, nil
}

func (w *day15Warehouse) tile(p grid.Position2D[int]) byte {
	return w.tiles[p.Y][p.X]
}

// push moves the robot one step in direction, shifting every box in the way,
// unless any of them is blocked by a wall. The walled boundary keeps every
// visited cell inside the map.
func (w *day15Warehouse) push(direction grid.Position2D[int]) {
	right := grid.Position2D[int]{X: 1}
	cells := []grid.Position2D[int]{w.robot}
	enqueue := func(p grid.Position2D[int]) {
		if !slices.Contains(cells, p) {
			cells = append(cells, p)
		}
	}
	for i := 0; i < len(cells); i++ {
		next := cells[i].Add(direction)
		switch w.tile(next) {
		case '#':
			return
		case 'O':
			enqueue(next)
		case '[':
			enqueue(next)
			enqueue(next.Add(right))
		case ']':
			enqueue(next)
			enqueue(next.Sub(right))
		}
	}
	// Cells are discovered in order of distance along direction, so moving
	// them in reverse empties each destination before it is filled.
	for _, cell := range slices.Backward(cells) {
		to := cell.Add(direction)
		w.tiles[to.Y][to.X] = w.tile(cell)
		w.tiles[cell.Y][cell.X] = '.'
	}
	w.robot = w.robot.Add(direction)
}

func day15GPSSum(r io.Reader, transform func(byte) []byte) (string, error) {
	warehouse, movements, err := parseWarehouse(r, transform)
	if err != nil {
		return "", err
	}
	for _, direction := range movements {
		warehouse.push(direction)
	}
	var total int
	for y, row := range warehouse.tiles {
		for x, tile := range row {
			if tile == 'O' || tile == '[' {
				total += 100*y + x
			}
		}
	}
	return strconv.Itoa(total), nil
}

func day15p01(r io.Reader) (string, error) {
	return day15GPSSum(r, warehouseIdentityTile)
}

func day15p02(r io.Reader) (string, error) {
	return day15GPSSum(r, warehouseScaleUpTile)
}
