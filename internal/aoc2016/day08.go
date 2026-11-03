package aoc2016

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func swipeScreen(r io.Reader, width, height int) ([][]bool, error) {
	pixels := make([][]bool, height)
	for y := range pixels {
		pixels[y] = make([]bool, width)
	}
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		var a, b int
		switch {
		case strings.HasPrefix(line, "rect "):
			n, err := fmt.Sscanf(line, "rect %dx%d", &a, &b)
			if err != nil || n != 2 || a < 0 || a > width || b < 0 || b > height {
				return nil, fmt.Errorf("invalid screen rectangle %q", line)
			}
			for y := range b {
				for x := range a {
					pixels[y][x] = true
				}
			}
		case strings.HasPrefix(line, "rotate row "):
			n, err := fmt.Sscanf(line, "rotate row y=%d by %d", &a, &b)
			if err != nil || n != 2 || a < 0 || a >= height || b < 0 {
				return nil, fmt.Errorf("invalid row rotation %q", line)
			}
			row := slices.Clone(pixels[a])
			for x, pixel := range row {
				pixels[a][(x+b%width)%width] = pixel
			}
		case strings.HasPrefix(line, "rotate column "):
			n, err := fmt.Sscanf(line, "rotate column x=%d by %d", &a, &b)
			if err != nil || n != 2 || a < 0 || a >= width || b < 0 {
				return nil, fmt.Errorf("invalid column rotation %q", line)
			}
			column := make([]bool, height)
			for y := range height {
				column[y] = pixels[y][a]
			}
			for y, pixel := range column {
				pixels[(y+b%height)%height][a] = pixel
			}
		default:
			return nil, fmt.Errorf("invalid screen instruction %q", line)
		}
	}
	return pixels, s.Err()
}

func screenRows(pixels [][]bool) []string {
	rows := make([]string, len(pixels))
	for y, row := range pixels {
		var text strings.Builder
		for _, lit := range row {
			if lit {
				text.WriteByte('#')
			} else {
				text.WriteByte(' ')
			}
		}
		rows[y] = text.String()
	}
	return rows
}

func day08p01(r io.Reader) (string, error) {
	pixels, err := swipeScreen(r, 50, 6)
	if err != nil {
		return "", err
	}
	total := 0
	for _, row := range pixels {
		for _, lit := range row {
			if lit {
				total++
			}
		}
	}
	return strconv.Itoa(total), nil
}

func day08p02(r io.Reader) (string, error) {
	pixels, err := swipeScreen(r, 50, 6)
	if err != nil {
		return "", err
	}
	return aoc.Decode6RowDisplay(screenRows(pixels), 5)
}
