package grid

import (
	"fmt"
	"iter"
	"strings"

	"golang.org/x/exp/constraints"
)

type Grid2D[T constraints.Signed, V any] map[Position2D[T]]V

func NewGrid2D[T constraints.Signed, V any](v [][]V) Grid2D[T, V] {
	result := make(Grid2D[T, V])

	for y, row := range v {
		for x, v := range row {
			result[Position2D[T]{X: T(x), Y: T(y)}] = v
		}
	}
	return result
}

func (g *Grid2D[T, V]) Dimensions() (T, T, T, T) {
	var minX, maxX, minY, maxY T

	for pos := range *g {
		minX = pos.X
		maxX = pos.X
		minY = pos.Y
		maxY = pos.Y
		break
	}

	for pos := range *g {
		minX = min(minX, pos.X)
		maxX = max(maxX, pos.X)
		minY = min(minY, pos.Y)
		maxY = max(maxY, pos.Y)
	}

	return minX, maxX, minY, maxY
}

func (g *Grid2D[T, V]) Contains(pos Position2D[T]) bool {
	_, found := (*g)[pos]
	return found
}

func (g *Grid2D[T, V]) ValidNeighbours4(pos Position2D[T]) iter.Seq[Position2D[T]] {
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

func (g *Grid2D[T, V]) ValidNeighbours8(pos Position2D[T]) iter.Seq[Position2D[T]] {
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

func (g *Grid2D[T, V]) TurnRight() Grid2D[T, V] {
	minX, _, minY, maxY := g.Dimensions()
	result := make(Grid2D[T, V])

	for pos, v := range *g {
		// normalize to origin
		x := pos.X - minX
		y := pos.Y - minY
		// apply rotation: (x, y) -> (maxY - minY - y, x)
		newX := (maxY - minY) - y
		newY := x
		result[Position2D[T]{X: newX, Y: newY}] = v
	}

	return result
}

func (g *Grid2D[T, V]) TurnLeft() Grid2D[T, V] {
	minX, maxX, minY, _ := g.Dimensions()
	result := make(Grid2D[T, V])

	for pos, v := range *g {
		// normalize to origin
		x := pos.X - minX
		y := pos.Y - minY
		// apply rotation: (x, y) -> (y, maxX - minX - x)
		newX := y
		newY := (maxX - minX) - x
		result[Position2D[T]{X: newX, Y: newY}] = v
	}

	return result
}

func (g *Grid2D[T, V]) FlipHorizontal() Grid2D[T, V] {
	minX, maxX, minY, _ := g.Dimensions()
	result := make(Grid2D[T, V])

	for pos, v := range *g {
		// normalize to origin
		x := pos.X - minX
		y := pos.Y - minY
		// flip horizontally: (x, y) -> (maxX - minX - x, y)
		newX := (maxX - minX) - x
		newY := y
		result[Position2D[T]{X: newX, Y: newY}] = v
	}

	return result
}

func (g *Grid2D[T, V]) FlipVertical() Grid2D[T, V] {
	minX, _, minY, maxY := g.Dimensions()
	result := make(Grid2D[T, V])

	for pos, v := range *g {
		// normalize to origin
		x := pos.X - minX
		y := pos.Y - minY
		// flip vertically: (x, y) -> (x, maxY - minY - y)
		newX := x
		newY := (maxY - minY) - y
		result[Position2D[T]{X: newX, Y: newY}] = v
	}

	return result
}

func (g *Grid2D[T, V]) GetRow(y T) []V {
	minX, maxX, _, _ := g.Dimensions()
	result := make([]V, 0, maxX-minX+1)

	for x := minX; x <= maxX; x++ {
		if v, exists := (*g)[Position2D[T]{X: x, Y: y}]; exists {
			result = append(result, v)
		}
	}

	return result
}

func (g *Grid2D[T, V]) GetColumn(x T) []V {
	_, _, minY, maxY := g.Dimensions()
	result := make([]V, 0, maxY-minY+1)

	for y := minY; y <= maxY; y++ {
		if v, exists := (*g)[Position2D[T]{X: x, Y: y}]; exists {
			result = append(result, v)
		}
	}

	return result
}

func (g *Grid2D[T, V]) Top() []V {
	_, _, minY, _ := g.Dimensions()
	return g.GetRow(minY)
}

func (g *Grid2D[T, V]) Bottom() []V {
	_, _, _, maxY := g.Dimensions()
	return g.GetRow(maxY)
}

func (g *Grid2D[T, V]) Left() []V {
	minX, _, _, _ := g.Dimensions()
	return g.GetColumn(minX)
}

func (g *Grid2D[T, V]) Right() []V {
	_, maxX, _, _ := g.Dimensions()
	return g.GetColumn(maxX)
}

func (g *Grid2D[T, V]) PrettyPrint(format func(V) string, empty string) {
	minX, maxX, minY, maxY := g.Dimensions()

	sb := new(strings.Builder)
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if v, exists := (*g)[Position2D[T]{X: x, Y: y}]; exists {
				sb.WriteString(format(v))
			} else {
				sb.WriteString(empty)
			}
		}
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
}
