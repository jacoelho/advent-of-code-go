package aoc2025

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/funcs"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseRedTilesLayout(r io.Reader) ([]grid.Position2D[int], error) {
	s := scanner.NewScanner(r, func(line []byte) (grid.Position2D[int], error) {
		nums := convert.ExtractDigits[int](string(line))
		if len(nums) != 2 {
			return grid.Position2D[int]{}, fmt.Errorf("expected 2 integers, got %d", len(nums))
		}
		return grid.NewPosition2D(nums[0], nums[1]), nil
	})
	return slices.Collect(s.Values()), s.Err()
}

type rectangle struct {
	bottomLeftCorner, topRightCorner grid.Position2D[int]
}

func newRectangle(p1, p2 grid.Position2D[int]) rectangle {
	return rectangle{
		bottomLeftCorner: grid.Position2D[int]{
			X: min(p1.X, p2.X),
			Y: min(p1.Y, p2.Y),
		},
		topRightCorner: grid.Position2D[int]{
			X: max(p1.X, p2.X),
			Y: max(p1.Y, p2.Y),
		},
	}
}

// Corners returns all four corners of the rectangle in order: bottom-left, top-left, bottom-right, top-right.
func (r rectangle) Corners() []grid.Position2D[int] {
	return []grid.Position2D[int]{
		r.bottomLeftCorner,
		grid.NewPosition2D(r.bottomLeftCorner.X, r.topRightCorner.Y),
		grid.NewPosition2D(r.topRightCorner.X, r.bottomLeftCorner.Y),
		r.topRightCorner,
	}
}

// Edges returns all four edges of the rectangle in order: bottom, right, top, left.
func (r rectangle) Edges() []edge {
	corners := r.Corners()
	return []edge{
		{corners[0], corners[2]},
		{corners[2], corners[3]},
		{corners[3], corners[1]},
		{corners[1], corners[0]},
	}
}

// Area returns the area of the rectangle including boundary points.
func (r rectangle) Area() int {
	dx := r.topRightCorner.X - r.bottomLeftCorner.X
	dy := r.topRightCorner.Y - r.bottomLeftCorner.Y
	return (xmath.Abs(dx) + 1) * (xmath.Abs(dy) + 1)
}

// intersectsPolygon reports whether any edge of the rectangle intersects with any polygon edge.
func (r rectangle) intersectsPolygon(polygonEdges collections.Set[edge]) bool {
	return xslices.Any(func(rectEdge edge) bool {
		for polyEdge := range polygonEdges {
			if rectEdge.intersects(polyEdge) {
				return true
			}
		}
		return false
	}, r.Edges())
}

// crossProduct calculates the cross product of vectors (p2-p1) and (p3-p1).
func crossProduct(p1, p2, p3 grid.Position2D[int]) int {
	return (p2.X-p1.X)*(p3.Y-p1.Y) - (p2.Y-p1.Y)*(p3.X-p1.X)
}

// orientation returns the orientation of three points.
// Returns -1 for clockwise, 0 for collinear, and +1 for counterclockwise.
func orientation(p1, p2, p3 grid.Position2D[int]) int {
	cross := crossProduct(p1, p2, p3)
	switch {
	case cross < 0:
		return -1
	case cross > 0:
		return 1
	default:
		return 0
	}
}

// edge represents a line segment between two points.
type edge struct {
	start, end grid.Position2D[int]
}

// contains reports whether the edge contains the given point.
func (e edge) contains(point grid.Position2D[int]) bool {
	if crossProduct(e.start, e.end, point) != 0 {
		return false
	}

	// is point is within the edge's bounding box.
	return (point.X >= min(e.start.X, e.end.X) && point.X <= max(e.start.X, e.end.X)) &&
		(point.Y >= min(e.start.Y, e.end.Y) && point.Y <= max(e.start.Y, e.end.Y))
}

// intersects reports whether this edge intersects with another edge.
func (e edge) intersects(other edge) bool {
	o1 := orientation(e.start, e.end, other.start)
	o2 := orientation(e.start, e.end, other.end)
	o3 := orientation(other.start, other.end, e.start)
	o4 := orientation(other.start, other.end, e.end)
	return o1*o2 < 0 && o3*o4 < 0
}

// findMaxRectangle finds the maximum area rectangle formed by pairs of tiles
// that satisfy the given validator function.
func findMaxRectangle(tiles []grid.Position2D[int], validator func(rectangle) bool) int {
	var maxArea int
	for pair := range xslices.Pairwise(tiles) {
		rect := newRectangle(pair.V1, pair.V2)
		if validator(rect) {
			maxArea = max(maxArea, rect.Area())
		}
	}
	return maxArea
}

func day09p01(r io.Reader) (string, error) {
	tiles, err := parseRedTilesLayout(r)
	if err != nil {
		return "", err
	}
	maxArea := findMaxRectangle(tiles, func(rectangle) bool { return true })
	return strconv.Itoa(maxArea), nil
}

func polygonEdgesSet(polygon []grid.Position2D[int]) collections.Set[edge] {
	n := len(polygon)
	edges := collections.NewSet[edge]()
	for i := range n {
		edges.Add(edge{
			start: polygon[i],
			end:   polygon[(i+1)%n],
		})
	}
	return edges
}

// isPointInside reports whether point is inside the polygon using the ray casting algorithm.
func isPointInside(polygon []grid.Position2D[int], point grid.Position2D[int]) bool {
	intersectionCount := 0
	numVertices := len(polygon)

	for i := range numVertices {
		vertex1 := polygon[i]
		vertex2 := polygon[(i+1)%numVertices]

		boundaryEdge := edge{start: vertex1, end: vertex2}
		if boundaryEdge.contains(point) {
			return true
		}

		// skip if both endpoints are on same side of ray or segment is horizontal.
		if (vertex1.Y > point.Y) == (vertex2.Y > point.Y) || vertex1.Y == vertex2.Y {
			continue
		}

		// count intersections with horizontal ray going right.
		intersectionX := (vertex2.X-vertex1.X)*(point.Y-vertex1.Y)/(vertex2.Y-vertex1.Y) + vertex1.X
		if point.X < intersectionX {
			intersectionCount++
		}
	}

	return intersectionCount%2 == 1
}

func day09p02(r io.Reader) (string, error) {
	tiles, err := parseRedTilesLayout(r)
	if err != nil {
		return "", err
	}

	isPointInsideMemo := funcs.Memoize(func(point grid.Position2D[int]) bool {
		return isPointInside(tiles, point)
	})

	polygonEdges := polygonEdgesSet(tiles)

	maxArea := findMaxRectangle(tiles, func(rect rectangle) bool {
		cornersInside := xslices.Every(isPointInsideMemo, rect.Corners())
		return cornersInside && !rect.intersectsPolygon(polygonEdges)
	})
	return strconv.Itoa(maxArea), nil
}
