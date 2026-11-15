package grid

import (
	"errors"
	"fmt"
	"iter"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"golang.org/x/exp/constraints"
)

// Position2D is a point or offset on a two-dimensional integer grid, with Y
// increasing downward.
type Position2D[T constraints.Signed] struct {
	X, Y T
}

// GetX returns p.X.
func (p Position2D[T]) GetX() T { return p.X }

// GetY returns p.Y.
func (p Position2D[T]) GetY() T { return p.Y }

// NewPosition2D returns the position (x, y).
func NewPosition2D[T constraints.Signed](x, y T) Position2D[T] {
	return Position2D[T]{X: x, Y: y}
}

// Add returns the component-wise sum of p and other.
func (p Position2D[T]) Add(other Position2D[T]) Position2D[T] {
	return Position2D[T]{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

// Sub returns the component-wise difference p - other.
func (p Position2D[T]) Sub(other Position2D[T]) Position2D[T] {
	return Position2D[T]{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

// Distance returns the Manhattan distance between p and other.
func (p Position2D[T]) Distance(other Position2D[T]) T {
	return xmath.Abs(p.X-other.X) + xmath.Abs(p.Y-other.Y)
}

// TurnRight returns p rotated 90 degrees clockwise about the origin, with Y
// increasing downward.
func (p Position2D[T]) TurnRight() Position2D[T] {
	return Position2D[T]{
		X: -p.Y,
		Y: p.X,
	}
}

// TurnLeft returns p rotated 90 degrees counterclockwise about the origin,
// with Y increasing downward.
func (p Position2D[T]) TurnLeft() Position2D[T] {
	return Position2D[T]{
		X: p.Y,
		Y: -p.X,
	}
}

// Within reports whether p lies in the rectangle [0, width) × [0, height).
func (p Position2D[T]) Within(width, height T) bool {
	return 0 <= p.X && p.X < width && 0 <= p.Y && p.Y < height
}

func (p Position2D[T]) generateNeighbours(offsets []Position2D[T]) iter.Seq[Position2D[T]] {
	return func(yield func(Position2D[T]) bool) {
		for _, offset := range offsets {
			if !yield(p.Add(offset)) {
				return
			}
		}
	}
}

// OffsetsNeighbours4 returns the four orthogonal unit offsets. Each call
// returns a fresh slice.
func OffsetsNeighbours4[T constraints.Signed]() []Position2D[T] {
	return []Position2D[T]{
		{X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1},
	}
}

// Neighbours4 returns an iterator over the orthogonal neighbours of p.
func Neighbours4[T constraints.Signed](p Position2D[T]) iter.Seq[Position2D[T]] {
	return p.generateNeighbours(OffsetsNeighbours4[T]())
}

// OffsetsNeighbours8 returns the eight orthogonal and diagonal unit offsets in
// row-major order. Each call returns a fresh slice.
func OffsetsNeighbours8[T constraints.Signed]() []Position2D[T] {
	return []Position2D[T]{
		{X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1},
		{X: -1, Y: 0}, {X: 1, Y: 0},
		{X: -1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: 1},
	}
}

// Neighbours8 returns an iterator over the orthogonal and diagonal neighbours
// of p.
func Neighbours8[T constraints.Signed](p Position2D[T]) iter.Seq[Position2D[T]] {
	return p.generateNeighbours(OffsetsNeighbours8[T]())
}

// ParsePosition2D parses "x,y", two base-10 integers that fit in T separated
// by a single comma, with no surrounding space.
func ParsePosition2D[T constraints.Signed](s string) (Position2D[T], error) {
	xs, ys, ok := strings.Cut(s, ",")
	if !ok {
		return Position2D[T]{}, fmt.Errorf("position %q is not x,y", s)
	}
	x, errX := parseCoordinate[T](xs)
	y, errY := parseCoordinate[T](ys)
	if err := errors.Join(errX, errY); err != nil {
		return Position2D[T]{}, fmt.Errorf("position %q: %w", s, err)
	}
	return Position2D[T]{X: x, Y: y}, nil
}

func parseCoordinate[T constraints.Signed](s string) (T, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err == nil && int64(T(v)) != v {
		err = fmt.Errorf("coordinate %s out of range", s)
	}
	return T(v), err
}
