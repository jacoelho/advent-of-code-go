package grid

import (
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"golang.org/x/exp/constraints"
)

type Position3D[T constraints.Signed] struct {
	X, Y, Z T
}

func NewPosition3D[T constraints.Signed](x, y, z T) Position3D[T] {
	return Position3D[T]{X: x, Y: y, Z: z}
}

func (p *Position3D[T]) Add(other Position3D[T]) Position3D[T] {
	return Position3D[T]{
		X: p.X + other.X,
		Y: p.Y + other.Y,
		Z: p.Z + other.Z,
	}
}

func (p *Position3D[T]) Sub(other Position3D[T]) Position3D[T] {
	return Position3D[T]{
		X: p.X - other.X,
		Y: p.Y - other.Y,
		Z: p.Z - other.Z,
	}
}

func (p *Position3D[T]) Distance(other Position3D[T]) T {
	return xmath.Abs(p.X-other.X) + xmath.Abs(p.Y-other.Y) + xmath.Abs(p.Z-other.Z)
}

func (p *Position3D[T]) EuclideanDistanceSquared(other Position3D[T]) T {
	dx, dy, dz := p.X-other.X, p.Y-other.Y, p.Z-other.Z
	return dx*dx + dy*dy + dz*dz
}
