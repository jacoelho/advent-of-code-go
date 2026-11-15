package aoc2020

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type hexTile = grid.Position2D[int]

// hexTilePosition follows a path of e, se, sw, w, nw and ne steps in axial
// coordinates where se is (1, 1) and nw is (-1, -1).
func hexTilePosition(path []byte) (hexTile, error) {
	if len(path) == 0 {
		return hexTile{}, fmt.Errorf("empty tile path")
	}
	var tile hexTile
	for i := 0; i < len(path); i++ {
		switch path[i] {
		case 'e':
			tile.X++
		case 'w':
			tile.X--
		case 'n', 's':
			if i+1 == len(path) || path[i+1] != 'e' && path[i+1] != 'w' {
				return hexTile{}, fmt.Errorf("invalid tile path %q at %d", path, i)
			}
			switch string(path[i : i+2]) {
			case "ne":
				tile.Y--
			case "sw":
				tile.Y++
			case "se":
				tile.X++
				tile.Y++
			case "nw":
				tile.X--
				tile.Y--
			}
			i++
		default:
			return hexTile{}, fmt.Errorf("invalid tile path %q at %d", path, i)
		}
	}
	return tile, nil
}

func initialBlackTiles(r io.Reader) ([]hexTile, error) {
	s := scanner.NewScanner(r, hexTilePosition)
	flipped := make(map[hexTile]bool)
	for tile := range s.Values() {
		flipped[tile] = !flipped[tile]
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(flipped) == 0 {
		return nil, fmt.Errorf("no tile paths")
	}
	var black []hexTile
	for tile, isBlack := range flipped {
		if isBlack {
			black = append(black, tile)
		}
	}
	return black, nil
}

func day24p01(r io.Reader) (string, error) {
	black, err := initialBlackTiles(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(len(black)), nil
}

func day24p02(r io.Reader) (string, error) {
	black, err := initialBlackTiles(r)
	if err != nil {
		return "", err
	}
	const days = 100
	if len(black) == 0 {
		return "0", nil
	}
	minX, maxX, minY, maxY := black[0].X, black[0].X, black[0].Y, black[0].Y
	for _, tile := range black[1:] {
		minX, maxX = min(minX, tile.X), max(maxX, tile.X)
		minY, maxY = min(minY, tile.Y), max(maxY, tile.Y)
	}
	// The black region grows by at most one tile per day in each direction, so
	// a margin of days+1 keeps every live tile and its neighbours inside the grid.
	const margin = days + 1
	width, height := maxX-minX+1+2*margin, maxY-minY+1+2*margin
	floor := make([]bool, width*height)
	for _, tile := range black {
		floor[(tile.Y-minY+margin)*width+tile.X-minX+margin] = true
	}
	neighbours := [6]int{1, -1, -width, width, width + 1, -width - 1}
	next := make([]bool, len(floor))
	count := 0
	for range days {
		count = 0
		for y := 1; y < height-1; y++ {
			for i := y*width + 1; i < (y+1)*width-1; i++ {
				adjacent := 0
				for _, offset := range neighbours {
					if floor[i+offset] {
						adjacent++
					}
				}
				next[i] = adjacent == 2 || floor[i] && adjacent == 1
				if next[i] {
					count++
				}
			}
		}
		floor, next = next, floor
	}
	return strconv.Itoa(count), nil
}
