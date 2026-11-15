package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

var seaMonster = []string{
	"                  # ",
	"#    ##    ##    ###",
	" #  #  #  #  #  #   ",
}

type imageTile struct {
	id     int
	pixels [][]byte
}

// imageTiles holds square tiles of one size forming a square arrangement, with
// each canonical edge shared by at most two tiles.
type imageTiles struct {
	tiles []imageTile
	side  int
	edges map[string][]int
}

func parseImageTiles(r io.Reader) (imageTiles, error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return imageTiles{}, err
	}
	var tiles []imageTile
	ids := make(map[int]bool)
	for _, lines := range blocks {
		if len(lines) == 0 {
			return imageTiles{}, fmt.Errorf("empty tile")
		}
		idText, ok := strings.CutPrefix(lines[0], "Tile ")
		idText, colon := strings.CutSuffix(idText, ":")
		id, err := strconv.Atoi(idText)
		if !ok || !colon || err != nil || id <= 0 {
			return imageTiles{}, fmt.Errorf("invalid tile header %q", lines[0])
		}
		if ids[id] {
			return imageTiles{}, fmt.Errorf("duplicate tile %d", id)
		}
		ids[id] = true
		size := len(lines) - 1
		if size < 3 || len(tiles) > 0 && size != len(tiles[0].pixels) {
			return imageTiles{}, fmt.Errorf("tile %d has %d rows", id, size)
		}
		pixels := make([][]byte, size)
		for y, line := range lines[1:] {
			if len(line) != size || strings.Trim(line, ".#") != "" {
				return imageTiles{}, fmt.Errorf("tile %d: invalid row %q", id, line)
			}
			pixels[y] = []byte(line)
		}
		tiles = append(tiles, imageTile{id, pixels})
	}
	side := 1
	for side*side < len(tiles) {
		side++
	}
	if side*side != len(tiles) || side < 2 {
		return imageTiles{}, fmt.Errorf("%d tiles do not form a square of at least 2x2", len(tiles))
	}
	edges := make(map[string][]int)
	for i, tile := range tiles {
		for _, edge := range tileEdges(tile.pixels) {
			key := canonicalEdge(edge)
			edges[key] = append(edges[key], i)
			if len(edges[key]) > 2 {
				return imageTiles{}, fmt.Errorf("edge %s is shared by more than two tiles", key)
			}
		}
	}
	return imageTiles{tiles, side, edges}, nil
}

// tileEdges returns the top, right, bottom and left edges of a tile.
func tileEdges(pixels [][]byte) [4]string {
	n := len(pixels)
	left, right := make([]byte, n), make([]byte, n)
	for y, row := range pixels {
		left[y], right[y] = row[0], row[n-1]
	}
	return [4]string{string(pixels[0]), string(right), string(pixels[n-1]), string(left)}
}

func canonicalEdge(edge string) string {
	reversed := []byte(edge)
	slices.Reverse(reversed)
	return min(edge, string(reversed))
}

func (t imageTiles) unmatched(edge string) bool {
	return len(t.edges[canonicalEdge(edge)]) == 1
}

func (t imageTiles) corners() ([]int, error) {
	var corners []int
	for i, tile := range t.tiles {
		unmatched := 0
		for _, edge := range tileEdges(tile.pixels) {
			if t.unmatched(edge) {
				unmatched++
			}
		}
		if unmatched == 2 {
			corners = append(corners, i)
		}
	}
	if len(corners) != 4 {
		return nil, fmt.Errorf("found %d corner tiles, want 4", len(corners))
	}
	return corners, nil
}

func rotateClockwise(pixels [][]byte) [][]byte {
	n := len(pixels)
	rotated := make([][]byte, n)
	for y := range rotated {
		rotated[y] = make([]byte, n)
		for x := range rotated[y] {
			rotated[y][x] = pixels[n-1-x][y]
		}
	}
	return rotated
}

// orientations returns the four rotations of pixels and of its mirror image.
func orientations(pixels [][]byte) [8][][]byte {
	var result [8][][]byte
	mirrored := make([][]byte, len(pixels))
	for y, row := range pixels {
		mirrored[y] = slices.Clone(row)
		slices.Reverse(mirrored[y])
	}
	result[0], result[4] = pixels, mirrored
	for i := 1; i < 4; i++ {
		result[i] = rotateClockwise(result[i-1])
		result[i+4] = rotateClockwise(result[i+3])
	}
	return result
}

func day20p01(r io.Reader) (string, error) {
	tiles, err := parseImageTiles(r)
	if err != nil {
		return "", err
	}
	corners, err := tiles.corners()
	if err != nil {
		return "", err
	}
	product := 1
	for _, corner := range corners {
		product *= tiles.tiles[corner].id
	}
	return strconv.Itoa(product), nil
}

// assemble arranges the tiles row by row from a corner. Each tile's top and
// left edges must match its placed neighbours, or be unmatched on the border.
func (t imageTiles) assemble() ([][][][]byte, error) {
	corners, err := t.corners()
	if err != nil {
		return nil, err
	}
	used := make([]bool, len(t.tiles))
	placed := make([][][][]byte, t.side)
	for y := range placed {
		placed[y] = make([][][]byte, t.side)
		for x := range placed[y] {
			var above, left string
			candidates := corners[:1]
			if y > 0 {
				above = tileEdges(placed[y-1][x])[2]
				candidates = t.edges[canonicalEdge(above)]
			}
			if x > 0 {
				left = tileEdges(placed[y][x-1])[1]
				candidates = t.edges[canonicalEdge(left)]
			}
			fits := func(pixels [][]byte) bool {
				edges := tileEdges(pixels)
				topFits := y > 0 && edges[0] == above || y == 0 && t.unmatched(edges[0])
				leftFits := x > 0 && edges[3] == left || x == 0 && t.unmatched(edges[3])
				return topFits && leftFits
			}
			tile, oriented, ok := t.place(candidates, used, fits)
			if !ok {
				return nil, fmt.Errorf("no tile fits at row %d column %d", y, x)
			}
			used[tile] = true
			placed[y][x] = oriented
		}
	}
	return placed, nil
}

func (t imageTiles) place(candidates []int, used []bool, fits func([][]byte) bool) (int, [][]byte, bool) {
	for _, candidate := range candidates {
		if used[candidate] {
			continue
		}
		for _, oriented := range orientations(t.tiles[candidate].pixels) {
			if fits(oriented) {
				return candidate, oriented, true
			}
		}
	}
	return 0, nil, false
}

func day20p02(r io.Reader) (string, error) {
	tiles, err := parseImageTiles(r)
	if err != nil {
		return "", err
	}
	placed, err := tiles.assemble()
	if err != nil {
		return "", err
	}
	inner := len(tiles.tiles[0].pixels) - 2
	image := make([][]byte, tiles.side*inner)
	for y := range image {
		row := placed[y/inner]
		for _, tile := range row {
			image[y] = append(image[y], tile[y%inner+1][1:inner+1]...)
		}
	}
	for _, oriented := range orientations(image) {
		if roughness, found := waterRoughness(oriented); found {
			return strconv.Itoa(roughness), nil
		}
	}
	return "", fmt.Errorf("no sea monsters found")
}

// waterRoughness counts '#' pixels outside sea monsters, reporting whether any
// monster was found.
func waterRoughness(image [][]byte) (int, bool) {
	monster := make([][]bool, len(image))
	for y := range monster {
		monster[y] = make([]bool, len(image[y]))
	}
	found := false
	for y := 0; y+len(seaMonster) <= len(image); y++ {
		for x := 0; x+len(seaMonster[0]) <= len(image[y]); x++ {
			if !seaMonsterAt(image, x, y) {
				continue
			}
			found = true
			for dy, pattern := range seaMonster {
				for dx := range pattern {
					if pattern[dx] == '#' {
						monster[y+dy][x+dx] = true
					}
				}
			}
		}
	}
	roughness := 0
	for y, row := range image {
		for x, pixel := range row {
			if pixel == '#' && !monster[y][x] {
				roughness++
			}
		}
	}
	return roughness, found
}

func seaMonsterAt(image [][]byte, x, y int) bool {
	for dy, pattern := range seaMonster {
		for dx := range pattern {
			if pattern[dx] == '#' && image[y+dy][x+dx] != '#' {
				return false
			}
		}
	}
	return true
}
