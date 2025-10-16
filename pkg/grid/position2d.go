package grid

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"golang.org/x/exp/constraints"
)

type Position2D[T constraints.Signed] struct {
	X, Y T
}

func (p Position2D[T]) GetX() T { return p.X }
func (p Position2D[T]) GetY() T { return p.Y }

func NewPosition2D[T constraints.Signed](x, y T) Position2D[T] {
	return Position2D[T]{X: x, Y: y}
}

func (p *Position2D[T]) Add(other Position2D[T]) Position2D[T] {
	return Position2D[T]{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p *Position2D[T]) Sub(other Position2D[T]) Position2D[T] {
	return Position2D[T]{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

func (p *Position2D[T]) Distance(other Position2D[T]) T {
	return xmath.Abs(p.X-other.X) + xmath.Abs(p.Y-other.Y)
}

func (p *Position2D[T]) TurnRight() Position2D[T] {
	return Position2D[T]{
		X: -p.Y,
		Y: p.X,
	}
}

func (p *Position2D[T]) TurnLeft() Position2D[T] {
	return Position2D[T]{
		X: p.Y,
		Y: -p.X,
	}
}

func (p *Position2D[T]) generateNeighbours(offsets []Position2D[T]) iter.Seq[Position2D[T]] {
	return func(yield func(Position2D[T]) bool) {
		for _, offset := range offsets {
			newP := p.Add(offset)
			if newP.X >= 0 && newP.Y >= 0 && !yield(newP) {
				return
			}
		}
	}
}

func OffsetsNeighbours4[T constraints.Signed]() []Position2D[T] {
	return []Position2D[T]{
		{X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1},
	}
}

func Neighbours4[T constraints.Signed](p Position2D[T]) iter.Seq[Position2D[T]] {
	return p.generateNeighbours(OffsetsNeighbours4[T]())
}

func OffsetsNeighbours8[T constraints.Signed]() []Position2D[T] {
	return []Position2D[T]{
		{X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1},
		{X: -1, Y: 0} /*  element */, {X: 1, Y: 0},
		{X: -1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: 1},
	}
}

func Neighbours8[T constraints.Signed](p Position2D[T]) iter.Seq[Position2D[T]] {
	return p.generateNeighbours(OffsetsNeighbours8[T]())
}
