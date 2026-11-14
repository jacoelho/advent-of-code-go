// Package grid provides integer positions in two and three dimensions and a
// sparse two-dimensional grid keyed by position.
package grid

import (
	"fmt"
	"iter"
	"strings"

	"golang.org/x/exp/constraints"
)

// Grid2D is a sparse two-dimensional grid mapping positions to values. Cells
// absent from the map are empty.
type Grid2D[T constraints.Signed, V any] map[Position2D[T]]V

// NewGrid2D returns a grid with v[y][x] stored at position (x, y).
func NewGrid2D[T constraints.Signed, V any](v [][]V) Grid2D[T, V] {
	result := make(Grid2D[T, V])

	for y, row := range v {
		for x, v := range row {
			result[Position2D[T]{X: T(x), Y: T(y)}] = v
		}
	}
	return result
}

// Dimensions returns the bounding box of the occupied positions as minX,
// maxX, minY, maxY. An empty grid returns all zeros.
func (g Grid2D[T, V]) Dimensions() (T, T, T, T) {
	var minX, maxX, minY, maxY T
	first := true
	for pos := range g {
		if first {
			minX, maxX, minY, maxY = pos.X, pos.X, pos.Y, pos.Y
			first = false
			continue
		}
		minX, maxX = min(minX, pos.X), max(maxX, pos.X)
		minY, maxY = min(minY, pos.Y), max(maxY, pos.Y)
	}
	return minX, maxX, minY, maxY
}

// Contains reports whether pos holds a value.
func (g Grid2D[T, V]) Contains(pos Position2D[T]) bool {
	_, found := g[pos]
	return found
}

// ValidNeighbours4 returns an iterator over the orthogonal neighbours of pos
// that hold a value.
func (g Grid2D[T, V]) ValidNeighbours4(pos Position2D[T]) iter.Seq[Position2D[T]] {
	return func(yield func(Position2D[T]) bool) {
		for neighbor := range Neighbours4(pos) {
			if g.Contains(neighbor) {
				if !yield(neighbor) {
					return
				}
			}
		}
	}
}

// ValidNeighbours8 returns an iterator over the orthogonal and diagonal
// neighbours of pos that hold a value.
func (g Grid2D[T, V]) ValidNeighbours8(pos Position2D[T]) iter.Seq[Position2D[T]] {
	return func(yield func(Position2D[T]) bool) {
		for neighbor := range Neighbours8(pos) {
			if g.Contains(neighbor) {
				if !yield(neighbor) {
					return
				}
			}
		}
	}
}

// TurnRight returns a new grid rotated 90 degrees clockwise, with Y increasing
// downward. The result's bounding box starts at the origin.
func (g Grid2D[T, V]) TurnRight() Grid2D[T, V] {
	return g.remap(func(x, y, maxX, maxY T) Position2D[T] {
		return Position2D[T]{X: maxY - y, Y: x}
	})
}

// TurnLeft returns a new grid rotated 90 degrees counterclockwise, with Y
// increasing downward. The result's bounding box starts at the origin.
func (g Grid2D[T, V]) TurnLeft() Grid2D[T, V] {
	return g.remap(func(x, y, maxX, maxY T) Position2D[T] {
		return Position2D[T]{X: y, Y: maxX - x}
	})
}

// FlipHorizontal returns a new grid mirrored left to right. The result's
// bounding box starts at the origin.
func (g Grid2D[T, V]) FlipHorizontal() Grid2D[T, V] {
	return g.remap(func(x, y, maxX, maxY T) Position2D[T] {
		return Position2D[T]{X: maxX - x, Y: y}
	})
}

// FlipVertical returns a new grid mirrored top to bottom. The result's
// bounding box starts at the origin.
func (g Grid2D[T, V]) FlipVertical() Grid2D[T, V] {
	return g.remap(func(x, y, maxX, maxY T) Position2D[T] {
		return Position2D[T]{X: x, Y: maxY - y}
	})
}

// remap moves every cell to place(x, y, maxX, maxY), with all coordinates
// relative to the bounding box's top-left corner.
func (g Grid2D[T, V]) remap(place func(x, y, maxX, maxY T) Position2D[T]) Grid2D[T, V] {
	minX, maxX, minY, maxY := g.Dimensions()
	result := make(Grid2D[T, V], len(g))
	for pos, v := range g {
		result[place(pos.X-minX, pos.Y-minY, maxX-minX, maxY-minY)] = v
	}
	return result
}

// GetRow returns the values in row y from left to right, omitting empty cells.
func (g Grid2D[T, V]) GetRow(y T) []V {
	minX, maxX, _, _ := g.Dimensions()
	return g.row(y, minX, maxX)
}

// GetColumn returns the values in column x from top to bottom, omitting empty
// cells.
func (g Grid2D[T, V]) GetColumn(x T) []V {
	_, _, minY, maxY := g.Dimensions()
	return g.column(x, minY, maxY)
}

func (g Grid2D[T, V]) row(y, minX, maxX T) []V {
	result := make([]V, 0, maxX-minX+1)
	for x := minX; x <= maxX; x++ {
		if v, exists := g[Position2D[T]{X: x, Y: y}]; exists {
			result = append(result, v)
		}
	}
	return result
}

func (g Grid2D[T, V]) column(x, minY, maxY T) []V {
	result := make([]V, 0, maxY-minY+1)
	for y := minY; y <= maxY; y++ {
		if v, exists := g[Position2D[T]{X: x, Y: y}]; exists {
			result = append(result, v)
		}
	}
	return result
}

// Top returns the values in the topmost occupied row, as GetRow does.
func (g Grid2D[T, V]) Top() []V {
	minX, maxX, minY, _ := g.Dimensions()
	return g.row(minY, minX, maxX)
}

// Bottom returns the values in the bottommost occupied row, as GetRow does.
func (g Grid2D[T, V]) Bottom() []V {
	minX, maxX, _, maxY := g.Dimensions()
	return g.row(maxY, minX, maxX)
}

// Left returns the values in the leftmost occupied column, as GetColumn does.
func (g Grid2D[T, V]) Left() []V {
	minX, _, minY, maxY := g.Dimensions()
	return g.column(minX, minY, maxY)
}

// Right returns the values in the rightmost occupied column, as GetColumn
// does.
func (g Grid2D[T, V]) Right() []V {
	_, maxX, minY, maxY := g.Dimensions()
	return g.column(maxX, minY, maxY)
}

// PrettyPrint writes the grid's bounding box to standard output, one row per
// line, rendering each value with format and each empty cell as empty.
func (g Grid2D[T, V]) PrettyPrint(format func(V) string, empty string) {
	minX, maxX, minY, maxY := g.Dimensions()

	sb := new(strings.Builder)
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if v, exists := g[Position2D[T]{X: x, Y: y}]; exists {
				sb.WriteString(format(v))
			} else {
				sb.WriteString(empty)
			}
		}
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
}
