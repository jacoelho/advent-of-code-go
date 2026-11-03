package aoc2023

import (
	"fmt"
	"io"
	"iter"
	"maps"
	"math/big"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func parseGarden(r io.Reader) (grid.Grid2D[int, rune], grid.Position2D[int], error) {
	g := make(grid.Grid2D[int, rune])
	var start grid.Position2D[int]
	var foundStart bool
	width := 0

	lineScanner := scanner.NewScanner(r, func(bytes []byte) (string, error) {
		return string(bytes), nil
	})

	for y, line := range xiter.Enumerate(lineScanner.Values()) {
		if y == 0 {
			width = len(line)
		}
		if width == 0 || len(line) != width {
			return nil, start, fmt.Errorf("garden row %d is empty or has inconsistent width", y+1)
		}
		if len(g)+len(line) > 1_000_000 {
			return nil, start, fmt.Errorf("garden exceeds one million cells")
		}
		for x, ch := range line {
			if ch != '.' && ch != '#' && ch != 'S' {
				return nil, start, fmt.Errorf("invalid garden tile %q at %d,%d", ch, x, y)
			}
			pos := grid.NewPosition2D(x, y)
			g[pos] = ch
			if ch == 'S' {
				if foundStart {
					return nil, start, fmt.Errorf("garden has multiple starting positions")
				}
				start = pos
				foundStart = true
			}
		}
	}

	if err := lineScanner.Err(); err != nil {
		return nil, grid.Position2D[int]{}, err
	}

	if !foundStart {
		return nil, grid.Position2D[int]{}, fmt.Errorf("starting position 'S' not found")
	}

	return g, start, nil
}

func isGardenPlot(g grid.Grid2D[int, rune], pos grid.Position2D[int]) bool {
	ch, exists := g[pos]
	return exists && (ch == '.' || ch == 'S')
}

func gardenNeighbours(g grid.Grid2D[int, rune]) func(grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
	return func(pos grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return xiter.From(grid.Neighbours4(pos)).Filter(func(next grid.Position2D[int]) bool {
			return isGardenPlot(g, next)
		}).All()
	}
}

func countReachableInSteps(distances map[grid.Position2D[int]]int, targetSteps int) int {
	if targetSteps > 0 && len(distances) == 1 {
		return 0
	}
	targetParity := targetSteps % 2
	return xiter.From(maps.Values(distances)).Filter(func(d int) bool {
		return d <= targetSteps && d%2 == targetParity
	}).Count()
}

func day21(r io.Reader, steps int) (string, error) {
	if steps < 0 {
		return "", fmt.Errorf("garden steps must be nonnegative")
	}
	g, start, err := parseGarden(r)
	if err != nil {
		return "", err
	}

	distances := search.BFSDistances(start, gardenNeighbours(g))
	count := countReachableInSteps(distances, steps)

	return strconv.Itoa(count), nil
}

func day21p01(r io.Reader) (string, error) {
	return day21(r, 64)
}

func day21p02(r io.Reader) (string, error) {
	g, start, err := parseGarden(r)
	if err != nil {
		return "", err
	}

	count, err := repeatedGardenPlots(g, start, 26501365)
	if err != nil {
		return "", err
	}
	return count.String(), nil
}

type gardenDistances struct {
	side      int
	distances map[grid.Position2D[int]]int
}

func (d gardenDistances) at(pos grid.Position2D[int]) int {
	x, y := xmath.Modulo(pos.X, d.side), xmath.Modulo(pos.Y, d.side)
	tileX, tileY := (pos.X-x)/d.side, (pos.Y-y)/d.side
	baseX, baseY := max(-1, min(1, tileX)), max(-1, min(1, tileY))
	base := grid.NewPosition2D(baseX*d.side+x, baseY*d.side+y)
	distance, ok := d.distances[base]
	if !ok {
		return -1
	}
	return distance + d.side*(xmath.Abs(tileX-baseX)+xmath.Abs(tileY-baseY))
}

func repeatedGardenPlot(g grid.Grid2D[int, rune], side int, pos grid.Position2D[int]) bool {
	return isGardenPlot(g, grid.NewPosition2D(xmath.Modulo(pos.X, side), xmath.Modulo(pos.Y, side)))
}

func certifyGardenDistances(g grid.Grid2D[int, rune], start grid.Position2D[int], side int) (gardenDistances, error) {
	neighbors := func(pos grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) {
			for _, offset := range grid.OffsetsNeighbours4[int]() {
				next := pos.Add(offset)
				if next.X < -2*side || next.Y < -2*side || next.X >= 3*side || next.Y >= 3*side {
					continue
				}
				if repeatedGardenPlot(g, side, next) && !yield(next) {
					return
				}
			}
		}
	}
	model := gardenDistances{side: side, distances: search.BFSDistances(start, neighbors)}

	// These 25 tile regimes cover every translated edge, including neighbors beyond the patch.
	// Descent to the unique zero proves reachability; the edge bound proves shortest distance.
	for y := -2 * side; y < 3*side; y++ {
		for x := -2 * side; x < 3*side; x++ {
			pos := grid.NewPosition2D(x, y)
			if !repeatedGardenPlot(g, side, pos) {
				continue
			}
			distance := model.at(pos)
			if distance == 0 && pos != start {
				return gardenDistances{}, fmt.Errorf("unsupported garden distances: additional zero at %d,%d", x, y)
			}
			descending := distance == 0
			for _, offset := range grid.OffsetsNeighbours4[int]() {
				next := pos.Add(offset)
				if !repeatedGardenPlot(g, side, next) {
					continue
				}
				nextDistance := model.at(next)
				if (distance < 0) != (nextDistance < 0) || distance >= 0 && xmath.Abs(distance-nextDistance) > 1 {
					return gardenDistances{}, fmt.Errorf("unsupported garden distances: edge at %d,%d", x, y)
				}
				if nextDistance == distance-1 {
					descending = true
				}
			}
			if distance > 0 && !descending {
				return gardenDistances{}, fmt.Errorf("unsupported garden distances: no predecessor at %d,%d", x, y)
			}
		}
	}
	return model, nil
}

func gardenCopies(distance, dimensions, side, steps int) *big.Int {
	if distance < 0 || distance > steps {
		return new(big.Int)
	}
	remaining := steps - distance
	parity := remaining % 2
	if dimensions == 0 {
		if parity == 0 {
			return big.NewInt(1)
		}
		return new(big.Int)
	}
	maximum := remaining / side
	if maximum < parity {
		return new(big.Int)
	}
	count := (maximum-parity)/2 + 1
	result := big.NewInt(int64(count))
	if dimensions == 2 {
		// For each admissible sum k, a quadrant contains k+1 pairs of tile offsets.
		result.Mul(result, big.NewInt(int64(count+parity)))
	}
	return result
}

func repeatedGardenPlots(g grid.Grid2D[int, rune], start grid.Position2D[int], steps int) (*big.Int, error) {
	if steps < 0 {
		return nil, fmt.Errorf("garden steps must be nonnegative")
	}
	_, maxX, _, maxY := g.Dimensions()
	side := maxX + 1
	if side != maxY+1 || side%2 == 0 || side > 201 {
		return nil, fmt.Errorf("unsupported garden shape: need an odd square at most 201 cells wide")
	}

	hasExit := false
	for _, offset := range grid.OffsetsNeighbours4[int]() {
		if repeatedGardenPlot(g, side, start.Add(offset)) {
			hasExit = true
			break
		}
	}
	if !hasExit {
		if steps == 0 {
			return big.NewInt(1), nil
		}
		return new(big.Int), nil
	}
	model, err := certifyGardenDistances(g, start, side)
	if err != nil {
		return nil, err
	}

	total := new(big.Int)
	for tileY := -1; tileY <= 1; tileY++ {
		for tileX := -1; tileX <= 1; tileX++ {
			dimensions := 0
			if tileX != 0 {
				dimensions++
			}
			if tileY != 0 {
				dimensions++
			}
			for y := range side {
				for x := range side {
					if !isGardenPlot(g, grid.NewPosition2D(x, y)) {
						continue
					}
					distance := model.at(grid.NewPosition2D(tileX*side+x, tileY*side+y))
					total.Add(total, gardenCopies(distance, dimensions, side, steps))
				}
			}
		}
	}
	return total, nil
}
