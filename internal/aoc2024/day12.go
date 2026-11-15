package aoc2024

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func parseGardenPlots(r io.Reader) ([][]byte, error) {
	g, err := grid.ReadBytes(r, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if err != nil {
		return nil, err
	}
	return g.Rows(), nil
}

// day12Directions lists the orthogonal directions in rotation order, so each
// entry and its successor form the two sides of one corner.
var day12Directions = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

// day12FencePrices returns the total price of all regions using perimeter and
// using the number of sides, which equals the number of corners.
func day12FencePrices(plots [][]byte) (int, int) {
	height, width := len(plots), len(plots[0])
	samePlant := func(x, y int, plant byte) bool {
		return x >= 0 && x < width && y >= 0 && y < height && plots[y][x] == plant
	}
	seen := make([]bool, width*height)
	var stack [][2]int
	var perimeterPrice, sidesPrice int
	for startY, row := range plots {
		for startX, plant := range row {
			if seen[startY*width+startX] {
				continue
			}
			seen[startY*width+startX] = true
			stack = append(stack[:0], [2]int{startX, startY})
			var area, perimeter, corners int
			for len(stack) > 0 {
				x, y := stack[len(stack)-1][0], stack[len(stack)-1][1]
				stack = stack[:len(stack)-1]
				area++
				for i, d := range day12Directions {
					nx, ny := x+d[0], y+d[1]
					a := samePlant(nx, ny, plant)
					if !a {
						perimeter++
					} else if !seen[ny*width+nx] {
						seen[ny*width+nx] = true
						stack = append(stack, [2]int{nx, ny})
					}
					e := day12Directions[(i+1)%4]
					b := samePlant(x+e[0], y+e[1], plant)
					diagonal := samePlant(x+d[0]+e[0], y+d[1]+e[1], plant)
					if !a && !b || a && b && !diagonal {
						corners++
					}
				}
			}
			perimeterPrice += area * perimeter
			sidesPrice += area * corners
		}
	}
	return perimeterPrice, sidesPrice
}

func day12p01(r io.Reader) (string, error) {
	plots, err := parseGardenPlots(r)
	if err != nil {
		return "", err
	}
	perimeterPrice, _ := day12FencePrices(plots)
	return strconv.Itoa(perimeterPrice), nil
}

func day12p02(r io.Reader) (string, error) {
	plots, err := parseGardenPlots(r)
	if err != nil {
		return "", err
	}
	_, sidesPrice := day12FencePrices(plots)
	return strconv.Itoa(sidesPrice), nil
}
