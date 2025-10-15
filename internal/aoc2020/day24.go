package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type direction uint8

const (
	east direction = iota
	southeast
	southwest
	west
	northwest
	northeast
)

var directionMap = map[string]direction{
	"e":  east,
	"w":  west,
	"se": southeast,
	"sw": southwest,
	"ne": northeast,
	"nw": northwest,
}

var directionsOffsets = map[direction]grid.Position2D[int]{
	east:      {X: 1, Y: 0},
	west:      {X: -1, Y: 0},
	northeast: {X: 0, Y: -1},
	southwest: {X: 0, Y: 1},
	southeast: {X: 1, Y: 1},
	northwest: {X: -1, Y: -1},
}

func parseDirections(s string) ([]direction, error) {
	var dirs []direction

	for i := 0; i < len(s); {
		if i+1 < len(s) {
			if dir, ok := directionMap[s[i:i+2]]; ok {
				dirs = append(dirs, dir)
				i += 2
				continue
			}
		}

		if dir, ok := directionMap[s[i:i+1]]; ok {
			dirs = append(dirs, dir)
			i++
		} else {
			return nil, fmt.Errorf("invalid char '%c' at %d", s[i], i)
		}
	}
	return dirs, nil
}

func parseTilePaths(r io.Reader) ([][]direction, error) {
	s := scanner.NewScanner(r, func(b []byte) ([]direction, error) {
		return parseDirections(string(b))
	})
	return slices.Collect(s.Values()), s.Err()
}

func findTile(directions []direction) grid.Position2D[int] {
	return xslices.Reduce(func(pos grid.Position2D[int], dir direction) grid.Position2D[int] {
		return pos.Add(directionsOffsets[dir])
	}, grid.Position2D[int]{X: 0, Y: 0}, directions)
}

func initialBlackTiles(r io.Reader) (map[grid.Position2D[int]]bool, error) {
	tiles, err := parseTilePaths(r)
	if err != nil {
		return nil, err
	}

	flipCount := make(map[grid.Position2D[int]]int)

	for _, path := range tiles {
		pos := findTile(path)
		flipCount[pos]++
	}

	// black tiles are flipped an odd number of times
	blackTiles := make(map[grid.Position2D[int]]bool)
	for pos, count := range flipCount {
		if count%2 == 1 {
			blackTiles[pos] = true
		}
	}

	return blackTiles, nil
}

func neighbors(pos grid.Position2D[int]) []grid.Position2D[int] {
	result := make([]grid.Position2D[int], 0, 6)
	for _, offset := range directionsOffsets {
		result = append(result, pos.Add(offset))
	}
	return result
}

func blackNeighborCount(pos grid.Position2D[int], blackTiles map[grid.Position2D[int]]bool) int {
	count := 0
	for _, neighbor := range neighbors(pos) {
		if blackTiles[neighbor] {
			count++
		}
	}
	return count
}

func nextDay(blackTiles map[grid.Position2D[int]]bool) map[grid.Position2D[int]]bool {
	tilesToCheck := make(map[grid.Position2D[int]]bool)
	for tile := range blackTiles {
		tilesToCheck[tile] = true
		for _, neighbor := range neighbors(tile) {
			tilesToCheck[neighbor] = true
		}
	}

	nextBlack := make(map[grid.Position2D[int]]bool)
	for tile := range tilesToCheck {
		count := blackNeighborCount(tile, blackTiles)
		isBlack := blackTiles[tile]

		if isBlack {
			// black stays black if it has 1 or 2 black neighbors
			if count == 1 || count == 2 {
				nextBlack[tile] = true
			}
		} else {
			// white flips to black if it has exactly 2 black neighbors
			if count == 2 {
				nextBlack[tile] = true
			}
		}
	}
	return nextBlack
}

func day24p01(r io.Reader) (string, error) {
	blackTiles, err := initialBlackTiles(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(len(blackTiles)), nil
}

func day24p02(r io.Reader) (string, error) {
	blackTiles, err := initialBlackTiles(r)
	if err != nil {
		return "", err
	}

	for range 100 {
		blackTiles = nextDay(blackTiles)
	}

	return strconv.Itoa(len(blackTiles)), nil
}
