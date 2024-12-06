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

func (g *Grid2D[T, V]) PrettyPrint(format func(V) string, empty string) {
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
