package aoc2025

import (
	"cmp"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseRedTilesLayout(r io.Reader) ([]grid.Position2D[int64], error) {
	s := scanner.NewScanner(r, func(line []byte) (grid.Position2D[int64], error) {
		xText, yText, ok := strings.Cut(string(line), ",")
		if !ok {
			return grid.Position2D[int64]{}, fmt.Errorf("red tile needs two comma-separated coordinates")
		}
		x, err := convert.ScanNumber[int64]([]byte(xText))
		if err != nil {
			return grid.Position2D[int64]{}, fmt.Errorf("red tile x coordinate: %w", err)
		}
		y, err := convert.ScanNumber[int64]([]byte(yText))
		if err != nil {
			return grid.Position2D[int64]{}, fmt.Errorf("red tile y coordinate: %w", err)
		}
		return grid.NewPosition2D(x, y), nil
	})
	var tiles []grid.Position2D[int64]
	for tile := range s.Values() {
		if len(tiles) == 1024 {
			return nil, fmt.Errorf("layout exceeds 1024 red tiles")
		}
		tiles = append(tiles, tile)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(tiles) < 2 {
		return nil, fmt.Errorf("layout needs at least two red tiles")
	}
	seen := make(map[grid.Position2D[int64]]bool)
	for _, tile := range tiles {
		if seen[tile] {
			return nil, fmt.Errorf("duplicate red tile")
		}
		seen[tile] = true
	}
	return tiles, nil
}

type tileRectangle struct {
	min, max grid.Position2D[int64]
}

func newTileRectangle(a, b grid.Position2D[int64]) tileRectangle {
	return tileRectangle{
		min: grid.NewPosition2D(min(a.X, b.X), min(a.Y, b.Y)),
		max: grid.NewPosition2D(max(a.X, b.X), max(a.Y, b.Y)),
	}
}

func (r tileRectangle) area() (int64, error) {
	// Unsigned differences retain the full distance across signed coordinate extremes.
	dx := uint64(r.max.X) - uint64(r.min.X)
	dy := uint64(r.max.Y) - uint64(r.min.Y)
	if dx >= math.MaxInt64 || dy >= math.MaxInt64 {
		return 0, fmt.Errorf("rectangle area exceeds int64")
	}
	width, height := int64(dx)+1, int64(dy)+1
	if width > math.MaxInt64/height {
		return 0, fmt.Errorf("rectangle area exceeds int64")
	}
	return width * height, nil
}

func tilePolygonEdges(tiles []grid.Position2D[int64]) ([]tileRectangle, error) {
	if len(tiles) < 4 {
		return nil, fmt.Errorf("tile polygon needs at least four vertices")
	}
	edges := make([]tileRectangle, len(tiles))
	for i, a := range tiles {
		b := tiles[(i+1)%len(tiles)]
		if a == b || a.X != b.X && a.Y != b.Y {
			return nil, fmt.Errorf("tile polygon edges must be nonzero and axis-aligned")
		}
		edges[i] = newTileRectangle(a, b)
	}
	for i, a := range edges {
		for j := i + 1; j < len(edges); j++ {
			b := edges[j]
			left, right := max(a.min.X, b.min.X), min(a.max.X, b.max.X)
			bottom, top := max(a.min.Y, b.min.Y), min(a.max.Y, b.max.Y)
			if left > right || bottom > top {
				continue
			}
			adjacent := j == i+1 || i == 0 && j == len(edges)-1
			if !adjacent || left != right || bottom != top {
				return nil, fmt.Errorf("tile polygon must be simple")
			}
		}
	}
	return edges, nil
}

func tileAxis(values []int64) []int64 {
	slices.Sort(values)
	values = slices.Compact(values)
	axis := make([]int64, 0, 2*len(values)-1)
	for i, value := range values {
		axis = append(axis, value)
		if i+1 < len(values) && uint64(values[i+1])-uint64(value) > 1 {
			axis = append(axis, value+1)
		}
	}
	return axis
}

type tileSpan struct {
	low, high int64
}

func coloredTileRow(edges []tileRectangle, y int64) ([]tileSpan, error) {
	var crossings []int64
	var spans []tileSpan
	for _, edge := range edges {
		if edge.min.X == edge.max.X {
			if edge.min.Y <= y && y < edge.max.Y {
				crossings = append(crossings, edge.min.X)
			}
			if edge.min.Y <= y && y <= edge.max.Y {
				spans = append(spans, tileSpan{low: edge.min.X, high: edge.min.X})
			}
		} else if edge.min.Y == y {
			spans = append(spans, tileSpan{low: edge.min.X, high: edge.max.X})
		}
	}
	slices.Sort(crossings)
	if len(crossings)%2 != 0 {
		return nil, fmt.Errorf("tile polygon has an unmatched scanline crossing")
	}
	for i := 0; i < len(crossings); i += 2 {
		spans = append(spans, tileSpan{low: crossings[i], high: crossings[i+1]})
	}
	slices.SortFunc(spans, func(a, b tileSpan) int { return cmp.Compare(a.low, b.low) })
	merged := spans[:0]
	for _, span := range spans {
		if len(merged) > 0 {
			last := &merged[len(merged)-1]
			if span.low <= last.high || last.high < math.MaxInt64 && span.low == last.high+1 {
				last.high = max(last.high, span.high)
				continue
			}
		}
		merged = append(merged, span)
	}
	return merged, nil
}

type redTileFloor struct {
	xs, ys  []int64
	stride  int
	outside []uint32
}

func newRedTileFloor(tiles []grid.Position2D[int64]) (*redTileFloor, error) {
	edges, err := tilePolygonEdges(tiles)
	if err != nil {
		return nil, err
	}
	xs, ys := make([]int64, len(tiles)), make([]int64, len(tiles))
	for i, tile := range tiles {
		xs[i], ys[i] = tile.X, tile.Y
	}
	xs, ys = tileAxis(xs), tileAxis(ys)
	if len(xs) > 4_000_000/len(ys) {
		return nil, fmt.Errorf("tile floor exceeds four million compressed cells")
	}
	floor := &redTileFloor{
		xs:      xs,
		ys:      ys,
		stride:  len(xs) + 1,
		outside: make([]uint32, (len(xs)+1)*(len(ys)+1)),
	}
	// Color is constant within every vertex singleton or intervening integer band.
	for row, y := range ys {
		spans, err := coloredTileRow(edges, y)
		if err != nil {
			return nil, err
		}
		span := 0
		var outside uint32
		for column, x := range xs {
			for span < len(spans) && spans[span].high < x {
				span++
			}
			if span == len(spans) || x < spans[span].low {
				outside++
			}
			floor.outside[(row+1)*floor.stride+column+1] = floor.outside[row*floor.stride+column+1] + outside
		}
	}
	return floor, nil
}

func tileAxisIndex(axis []int64, value int64) int {
	index, found := slices.BinarySearch(axis, value)
	if !found {
		index--
	}
	return index
}

func (f *redTileFloor) contains(r tileRectangle) bool {
	if r.min.X < f.xs[0] || r.max.X > f.xs[len(f.xs)-1] || r.min.Y < f.ys[0] || r.max.Y > f.ys[len(f.ys)-1] {
		return false
	}
	left, right := tileAxisIndex(f.xs, r.min.X), tileAxisIndex(f.xs, r.max.X)+1
	bottom, top := tileAxisIndex(f.ys, r.min.Y), tileAxisIndex(f.ys, r.max.Y)+1
	return f.outside[top*f.stride+right]+f.outside[bottom*f.stride+left] ==
		f.outside[top*f.stride+left]+f.outside[bottom*f.stride+right]
}

func largestTileRectangle(tiles []grid.Position2D[int64], allowed func(tileRectangle) bool) (int64, error) {
	var largest int64
	for pair := range xslices.Pairwise(tiles) {
		rectangle := newTileRectangle(pair.V1, pair.V2)
		if !allowed(rectangle) {
			continue
		}
		area, err := rectangle.area()
		if err != nil {
			return 0, err
		}
		largest = max(largest, area)
	}
	return largest, nil
}

func day09p01(r io.Reader) (string, error) {
	tiles, err := parseRedTilesLayout(r)
	if err != nil {
		return "", err
	}
	area, err := largestTileRectangle(tiles, func(tileRectangle) bool { return true })
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(area, 10), nil
}

func day09p02(r io.Reader) (string, error) {
	tiles, err := parseRedTilesLayout(r)
	if err != nil {
		return "", err
	}
	floor, err := newRedTileFloor(tiles)
	if err != nil {
		return "", err
	}
	area, err := largestTileRectangle(tiles, floor.contains)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(area, 10), nil
}
