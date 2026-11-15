package aoc2024

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type wordSearch [][]byte

func parseWordSearch(r io.Reader) (wordSearch, error) {
	g, err := grid.ReadBytes(r, "XMAS")
	if err != nil {
		return nil, err
	}
	return g.Rows(), nil
}

func (w wordSearch) at(x, y int) byte {
	if y < 0 || y >= len(w) || x < 0 || x >= len(w[y]) {
		return 0
	}
	return w[y][x]
}

func day04p01(r io.Reader) (string, error) {
	w, err := parseWordSearch(r)
	if err != nil {
		return "", err
	}

	var count int
	for y, row := range w {
		for x, letter := range row {
			if letter != 'X' {
				continue
			}
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if w.at(x+dx, y+dy) == 'M' && w.at(x+2*dx, y+2*dy) == 'A' && w.at(x+3*dx, y+3*dy) == 'S' {
						count++
					}
				}
			}
		}
	}

	return strconv.Itoa(count), nil
}

func day04p02(r io.Reader) (string, error) {
	w, err := parseWordSearch(r)
	if err != nil {
		return "", err
	}

	var count int
	for y, row := range w {
		for x, letter := range row {
			if letter != 'A' {
				continue
			}
			diagonal := [2]byte{w.at(x-1, y-1), w.at(x+1, y+1)}
			antiDiagonal := [2]byte{w.at(x-1, y+1), w.at(x+1, y-1)}
			if wordSearchMS(diagonal) && wordSearchMS(antiDiagonal) {
				count++
			}
		}
	}

	return strconv.Itoa(count), nil
}

func wordSearchMS(ends [2]byte) bool {
	return ends == [2]byte{'M', 'S'} || ends == [2]byte{'S', 'M'}
}
