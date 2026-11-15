package aoc2023

import (
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

var gardenOffsets = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

type garden struct {
	rows           [][]byte
	startX, startY int
}

func parseGarden(r io.Reader) (garden, error) {
	g, err := grid.ReadBytes(r, ".#S")
	if err != nil {
		return garden{}, err
	}
	start, err := g.Find('S')
	if err != nil {
		return garden{}, err
	}
	return garden{rows: g.Rows(), startX: start.X, startY: start.Y}, nil
}

func (g garden) width() int  { return len(g.rows[0]) }
func (g garden) height() int { return len(g.rows) }

// plot reports whether x,y is a garden plot in the garden repeated infinitely.
func (g garden) plot(x, y int) bool {
	return g.rows[xmath.Modulo(y, g.height())][xmath.Modulo(x, g.width())] != '#'
}

func (g garden) hasExit() bool {
	for _, offset := range gardenOffsets {
		if g.plot(g.startX+offset[0], g.startY+offset[1]) {
			return true
		}
	}
	return false
}

// distances returns the step distances from the start to every cell of the
// width×height window whose corner is x0,y0, walking only inside the window.
// Cells are row-major and unreachable cells hold -1.
func (g garden) distances(x0, y0, width, height int) []int {
	distance := make([]int, width*height)
	for i := range distance {
		distance[i] = -1
	}
	startIndex := (g.startY-y0)*width + g.startX - x0
	distance[startIndex] = 0
	queue := []int{startIndex}
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		x, y := current%width, current/width
		for _, offset := range gardenOffsets {
			nx, ny := x+offset[0], y+offset[1]
			if nx < 0 || ny < 0 || nx >= width || ny >= height || !g.plot(nx+x0, ny+y0) {
				continue
			}
			if next := ny*width + nx; distance[next] < 0 {
				distance[next] = distance[current] + 1
				queue = append(queue, next)
			}
		}
	}
	return distance
}

func day21(r io.Reader, steps int) (string, error) {
	if steps < 0 {
		return "", fmt.Errorf("garden steps must be nonnegative")
	}
	g, err := parseGarden(r)
	if err != nil {
		return "", err
	}
	count, reachable := 0, 0
	for _, d := range g.distances(0, 0, g.width(), g.height()) {
		if d >= 0 {
			reachable++
		}
		if d >= 0 && d <= steps && d%2 == steps%2 {
			count++
		}
	}
	if steps > 0 && reachable == 1 {
		return "0", nil
	}
	return strconv.Itoa(count), nil
}

func day21p01(r io.Reader) (string, error) {
	return day21(r, 64)
}

func day21p02(r io.Reader) (string, error) {
	g, err := parseGarden(r)
	if err != nil {
		return "", err
	}
	count, err := repeatedGardenPlots(g, 26501365)
	if err != nil {
		return "", err
	}
	return count.String(), nil
}

// gardenDistances holds exact distances in the 5×5 tile patch around the
// start tile and extends them to farther tiles by whole-tile translations.
type gardenDistances struct {
	side  int
	patch []int
}

func (d gardenDistances) at(x, y int) int {
	modX, modY := xmath.Modulo(x, d.side), xmath.Modulo(y, d.side)
	tileX, tileY := (x-modX)/d.side, (y-modY)/d.side
	baseX, baseY := max(-1, min(1, tileX)), max(-1, min(1, tileY))
	size := 5 * d.side
	distance := d.patch[((baseY+2)*d.side+modY)*size+(baseX+2)*d.side+modX]
	if distance < 0 {
		return -1
	}
	return distance + d.side*(xmath.Abs(tileX-baseX)+xmath.Abs(tileY-baseY))
}

func certifyGardenDistances(g garden, side int) (gardenDistances, error) {
	model := gardenDistances{side: side, patch: g.distances(-2*side, -2*side, 5*side, 5*side)}

	// These 25 tile regimes cover every translated edge, including neighbors beyond the patch.
	// Descent to the unique zero proves reachability; the edge bound proves shortest distance.
	for y := -2 * side; y < 3*side; y++ {
		for x := -2 * side; x < 3*side; x++ {
			if !g.plot(x, y) {
				continue
			}
			distance := model.at(x, y)
			if distance == 0 && (x != g.startX || y != g.startY) {
				return gardenDistances{}, fmt.Errorf("unsupported garden distances: additional zero at %d,%d", x, y)
			}
			descending := distance == 0
			for _, offset := range gardenOffsets {
				nx, ny := x+offset[0], y+offset[1]
				if !g.plot(nx, ny) {
					continue
				}
				nextDistance := model.at(nx, ny)
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

func repeatedGardenPlots(g garden, steps int) (*big.Int, error) {
	if steps < 0 {
		return nil, fmt.Errorf("garden steps must be nonnegative")
	}
	side := g.width()
	if side != g.height() || side%2 == 0 {
		return nil, fmt.Errorf("unsupported garden shape: need an odd square")
	}
	if !g.hasExit() {
		if steps == 0 {
			return big.NewInt(1), nil
		}
		return new(big.Int), nil
	}
	model, err := certifyGardenDistances(g, side)
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
					if !g.plot(x, y) {
						continue
					}
					distance := model.at(tileX*side+x, tileY*side+y)
					total.Add(total, gardenCopies(distance, dimensions, side, steps))
				}
			}
		}
	}
	return total, nil
}
