package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func enhanceImage(r io.Reader, steps int) (string, error) {
	s := bufio.NewScanner(r)
	var algorithm [512]bool
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("missing image enhancement algorithm")
	}
	if len(s.Bytes()) != len(algorithm) {
		return "", fmt.Errorf("image enhancement algorithm must have 512 pixels")
	}
	for i, pixel := range s.Bytes() {
		if pixel != '.' && pixel != '#' {
			return "", fmt.Errorf("invalid algorithm pixel %q", pixel)
		}
		algorithm[i] = pixel == '#'
	}
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("missing image separator")
	}
	if s.Text() != "" {
		return "", fmt.Errorf("expected blank image separator")
	}

	exceptions := make(map[grid.Position2D[int]]bool)
	width, height := 0, 0
	for s.Scan() {
		line := s.Bytes()
		if height == 0 {
			width = len(line)
		}
		if width == 0 || len(line) != width || (height+1)*width > 1_000_000 {
			return "", fmt.Errorf("image must be rectangular and at most 1000000 pixels")
		}
		for x, pixel := range line {
			if pixel != '.' && pixel != '#' {
				return "", fmt.Errorf("invalid image pixel %q", pixel)
			}
			if pixel == '#' {
				exceptions[grid.Position2D[int]{X: x, Y: height}] = true
			}
		}
		height++
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if height == 0 {
		return "", fmt.Errorf("missing image")
	}

	background, work := false, 0
	for range steps {
		nextBackground := algorithm[0]
		if background {
			nextBackground = algorithm[511]
		}
		if len(exceptions) == 0 {
			background = nextBackground
			continue
		}
		xmin, xmax, ymin, ymax := 0, 0, 0, 0
		first := true
		for point := range exceptions {
			if first {
				xmin, xmax, ymin, ymax = point.X, point.X, point.Y, point.Y
				first = false
			}
			xmin, xmax = min(xmin, point.X), max(xmax, point.X)
			ymin, ymax = min(ymin, point.Y), max(ymax, point.Y)
		}
		area := (xmax - xmin + 3) * (ymax - ymin + 3)
		work += area
		if area > 4_000_000 || work > 100_000_000 {
			return "", fmt.Errorf("image enhancement exceeds area or work limit")
		}
		next := make(map[grid.Position2D[int]]bool)
		for y := ymin - 1; y <= ymax+1; y++ {
			for x := xmin - 1; x <= xmax+1; x++ {
				index := 0
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						index <<= 1
						if exceptions[grid.Position2D[int]{X: x + dx, Y: y + dy}] != background {
							index++
						}
					}
				}
				if algorithm[index] != nextBackground {
					next[grid.Position2D[int]{X: x, Y: y}] = true
				}
			}
		}
		exceptions, background = next, nextBackground
	}
	if background {
		return "", fmt.Errorf("enhanced image contains infinitely many lit pixels")
	}
	return strconv.Itoa(len(exceptions)), nil
}

func day20p01(r io.Reader) (string, error) {
	return enhanceImage(r, 2)
}

func day20p02(r io.Reader) (string, error) {
	return enhanceImage(r, 50)
}
