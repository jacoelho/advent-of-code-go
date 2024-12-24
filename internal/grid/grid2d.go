package grid

import (
	"fmt"
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
