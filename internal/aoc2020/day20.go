package aoc2020

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

var seaMonsterPattern = []grid.Position2D[int]{
	{X: 18, Y: 0},
	{X: 0, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 11, Y: 1}, {X: 12, Y: 1}, {X: 17, Y: 1}, {X: 18, Y: 1}, {X: 19, Y: 1},
	{X: 1, Y: 2}, {X: 4, Y: 2}, {X: 7, Y: 2}, {X: 10, Y: 2}, {X: 13, Y: 2}, {X: 16, Y: 2},
}

type imageTile struct {
	id   int
	grid grid.Grid2D[int, rune]
}

func parseImageTile(s string) (imageTile, error) {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) < 2 {
		return imageTile{}, fmt.Errorf("invalid tile format")
	}

	digits := convert.ExtractDigits[int](lines[0])
	if len(digits) == 0 {
		return imageTile{}, fmt.Errorf("failed to parse tile ID")
	}

	rows := make([][]rune, len(lines)-1)
	for i, line := range lines[1:] {
		rows[i] = []rune(line)
	}

	return imageTile{
		id:   digits[0],
		grid: grid.NewGrid2D[int](rows),
	}, nil
}

func parseImageTiles(r io.Reader) ([]imageTile, error) {
	s := bufio.NewScanner(r)
	s.Split(scanner.SplitBySeparator([]byte{'\n', '\n'}))

	var tiles []imageTile
	for s.Scan() {
		tile, err := parseImageTile(s.Text())
		if err != nil {
			return nil, err
		}
		tiles = append(tiles, tile)
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return tiles, nil
}

func reverseBorder(border string) string {
	runes := []rune(border)
	slices.Reverse(runes)
	return string(runes)
}

func getBorders(tile imageTile) []string {
	return []string{
		string(tile.grid.Top()),
		string(tile.grid.Bottom()),
		string(tile.grid.Left()),
		string(tile.grid.Right()),
	}
}

func bordersMatch(border1, border2 string) bool {
	return border1 == border2 || border1 == reverseBorder(border2)
}

func countMatchingBorders(tile imageTile, allTiles []imageTile) int {
	tileBorders := getBorders(tile)

	matchCount := 0
	for _, border := range tileBorders {
		if xiter.FromSlice(allTiles).Any(func(otherTile imageTile) bool {
			if tile.id == otherTile.id {
				return false
			}
			otherBorders := getBorders(otherTile)
			return xiter.FromSlice(otherBorders).Any(func(otherBorder string) bool {
				return bordersMatch(border, otherBorder)
			})
		}) {
			matchCount++
		}
	}

	return matchCount
}

func findCorners(tiles []imageTile) xiter.Seq[imageTile] {
	return xiter.FromSlice(tiles).Filter(func(tile imageTile) bool {
		return countMatchingBorders(tile, tiles) == 2
	})
}

func day20p01(r io.Reader) (string, error) {
	tiles, err := parseImageTiles(r)
	if err != nil {
		return "", err
	}

	result := findCorners(tiles).
		Map(func(tile imageTile) int { return tile.id }).
		Reduce(1, func(product, id int) int { return product * id })

	return strconv.Itoa(result), nil
}

// Generate all 8 orientations of a tile grid
func allOrientations(g grid.Grid2D[int, rune]) iter.Seq[grid.Grid2D[int, rune]] {
	return func(yield func(grid.Grid2D[int, rune]) bool) {
		current := g
		for range 4 {
			if !yield(current) {
				return
			}
			current = current.TurnRight()
		}

		current = g.FlipHorizontal()
		for range 4 {
			if !yield(current) {
				return
			}
			current = current.TurnRight()
		}
	}
}

func findTileWithBorder(border string, tiles []imageTile, excludeID int) (int, bool) {
	for _, tile := range tiles {
		if tile.id == excludeID {
			continue
		}
		borders := getBorders(tile)
		for _, b := range borders {
			if bordersMatch(border, b) {
				return tile.id, true
			}
		}
	}
	return 0, false
}

func findMatchingTileOrientation(tiles []imageTile, used map[int]bool, leftBorder, topBorder string) (imageTile, bool) {
	for _, tile := range tiles {
		if used[tile.id] {
			continue
		}

		for oriented := range allOrientations(tile.grid) {
			matchesLeft := leftBorder == "" || string(oriented.Left()) == leftBorder
			matchesTop := topBorder == "" || string(oriented.Top()) == topBorder

			if matchesLeft && matchesTop {
				return imageTile{id: tile.id, grid: oriented}, true
			}
		}
	}
	return imageTile{}, false
}

func assembleTiles(tiles []imageTile) map[grid.Position2D[int]]imageTile {
	cornerTile, ok := findCorners(tiles).First()
	if !ok {
		return nil
	}

	// Orient the first corner so matching borders face right and down.

	var startOriented grid.Grid2D[int, rune]

	for oriented := range allOrientations(cornerTile.grid) {
		rightBorder := string(oriented.Right())
		bottomBorder := string(oriented.Bottom())

		_, hasRight := findTileWithBorder(rightBorder, tiles, cornerTile.id)
		_, hasBottom := findTileWithBorder(bottomBorder, tiles, cornerTile.id)

		if hasRight && hasBottom {
			startOriented = oriented
			break
		}
	}

	placed := make(map[grid.Position2D[int]]imageTile)
	placed[grid.Position2D[int]{X: 0, Y: 0}] = imageTile{id: cornerTile.id, grid: startOriented}
	used := make(map[int]bool)
	used[cornerTile.id] = true

	gridSize := 1
	for gridSize*gridSize < len(tiles) {
		gridSize++
	}

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if x == 0 && y == 0 {
				continue
			}

			var leftBorder, topBorder string
			if x > 0 {
				if leftTile, exists := placed[grid.Position2D[int]{X: x - 1, Y: y}]; exists {
					leftBorder = string(leftTile.grid.Right())
				}
			}
			if y > 0 {
				if topTile, exists := placed[grid.Position2D[int]{X: x, Y: y - 1}]; exists {
					topBorder = string(topTile.grid.Bottom())
				}
			}

			if leftBorder == "" && topBorder == "" {
				continue
			}

			if matchedTile, found := findMatchingTileOrientation(tiles, used, leftBorder, topBorder); found {
				placed[grid.Position2D[int]{X: x, Y: y}] = matchedTile
				used[matchedTile.id] = true
			}
		}
	}

	return placed
}

// remove borders from a tile (remove first and last row/column)
func removeBorder(g grid.Grid2D[int, rune]) grid.Grid2D[int, rune] {
	minX, maxX, minY, maxY := g.Dimensions()
	result := make(grid.Grid2D[int, rune])

	for y := minY + 1; y < maxY; y++ {
		for x := minX + 1; x < maxX; x++ {
			if val, exists := g[grid.Position2D[int]{X: x, Y: y}]; exists {
				// normalize to start at 0,0
				newPos := grid.Position2D[int]{X: x - minX - 1, Y: y - minY - 1}
				result[newPos] = val
			}
		}
	}

	return result
}

func mergeToSingleGrid(assembled map[grid.Position2D[int]]imageTile) grid.Grid2D[int, rune] {
	if len(assembled) == 0 {
		return grid.Grid2D[int, rune]{}
	}

	result := make(grid.Grid2D[int, rune])

	for tilePos, tile := range assembled {
		borderless := removeBorder(tile.grid)
		minX, maxX, minY, maxY := borderless.Dimensions()
		tileWidth := maxX - minX + 1

		// copy borderless tile to final position
		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				val := borderless[grid.Position2D[int]{X: x, Y: y}]
				finalX := tilePos.X*tileWidth + (x - minX)
				finalY := tilePos.Y*tileWidth + (y - minY)
				result[grid.Position2D[int]{X: finalX, Y: finalY}] = val
			}
		}
	}

	return result
}

func findSeaMonsters(g grid.Grid2D[int, rune]) []grid.Position2D[int] {
	minX, maxX, minY, maxY := g.Dimensions()
	var monsters []grid.Position2D[int]

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			found := true
			for _, offset := range seaMonsterPattern {
				pos := grid.Position2D[int]{X: x + offset.X, Y: y + offset.Y}
				val, exists := g[pos]
				if !exists || val != '#' {
					found = false
					break
				}
			}
			if found {
				monsters = append(monsters, grid.Position2D[int]{X: x, Y: y})
			}
		}
	}

	return monsters
}

// count water roughness (# not part of sea monsters)
func countRoughness(g grid.Grid2D[int, rune], monsters []grid.Position2D[int]) int {
	allHash := collections.NewSet[grid.Position2D[int]]()
	for pos, val := range g {
		if val == '#' {
			allHash.Add(pos)
		}
	}

	monsterPositions := collections.NewSet[grid.Position2D[int]]()
	for _, monsterStart := range monsters {
		for _, offset := range seaMonsterPattern {
			pos := grid.Position2D[int]{X: monsterStart.X + offset.X, Y: monsterStart.Y + offset.Y}
			monsterPositions.Add(pos)
		}
	}

	return allHash.Difference(monsterPositions).Len()
}

func day20p02(r io.Reader) (string, error) {
	tiles, err := parseImageTiles(r)
	if err != nil {
		return "", err
	}

	assembled := assembleTiles(tiles)
	if len(assembled) != len(tiles) {
		return "", fmt.Errorf("failed to assemble all tiles: placed %d of %d", len(assembled), len(tiles))
	}

	fullImage := mergeToSingleGrid(assembled)

	for oriented := range allOrientations(fullImage) {
		monsters := findSeaMonsters(oriented)
		if len(monsters) > 0 {
			roughness := countRoughness(oriented, monsters)
			return strconv.Itoa(roughness), nil
		}
	}

	return "", fmt.Errorf("no sea monsters found")
}
