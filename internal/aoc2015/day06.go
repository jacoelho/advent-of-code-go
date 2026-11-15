package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type lightMode int

const (
	onOffLights lightMode = iota
	brightnessLights
)

type lightOperation uint8

const (
	turnOnLights lightOperation = iota
	turnOffLights
	toggleLights
)

func parseLightCorner(text string) (grid.Position2D[int], error) {
	p, err := grid.ParsePosition2D[int](text)
	if err != nil {
		return grid.Position2D[int]{}, err
	}
	if !p.Within(1000, 1000) {
		return grid.Position2D[int]{}, fmt.Errorf("light coordinate %q outside the grid", text)
	}
	return p, nil
}

func lightDisplay(r io.Reader, mode lightMode) (string, error) {
	lights := make([]int, 1000*1000)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		var operation lightOperation
		var rest string
		if after, ok := strings.CutPrefix(line, "turn on "); ok {
			operation, rest = turnOnLights, after
		} else if after, ok := strings.CutPrefix(line, "turn off "); ok {
			operation, rest = turnOffLights, after
		} else if after, ok := strings.CutPrefix(line, "toggle "); ok {
			operation, rest = toggleLights, after
		} else {
			return "", fmt.Errorf("invalid lighting instruction %q", line)
		}
		from, to, ok := strings.Cut(rest, " through ")
		if !ok {
			return "", fmt.Errorf("invalid lighting instruction %q", line)
		}
		lo, err := parseLightCorner(from)
		if err != nil {
			return "", fmt.Errorf("lighting instruction %q: %w", line, err)
		}
		hi, err := parseLightCorner(to)
		if err != nil {
			return "", fmt.Errorf("lighting instruction %q: %w", line, err)
		}
		if lo.X > hi.X || lo.Y > hi.Y {
			return "", fmt.Errorf("inverted light rectangle %q", line)
		}
		for y := lo.Y; y <= hi.Y; y++ {
			row := lights[y*1000+lo.X : y*1000+hi.X+1]
			for x := range row {
				light := &row[x]
				if mode == onOffLights {
					switch operation {
					case turnOnLights:
						*light = 1
					case turnOffLights:
						*light = 0
					case toggleLights:
						*light = 1 - *light
					}
				} else {
					switch operation {
					case turnOnLights:
						*light++
					case turnOffLights:
						*light = max(0, *light-1)
					case toggleLights:
						*light += 2
					}
				}
			}
		}
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	total := 0
	for _, light := range lights {
		total += light
	}
	return strconv.Itoa(total), nil
}

func day06p01(r io.Reader) (string, error) {
	return lightDisplay(r, onOffLights)
}

func day06p02(r io.Reader) (string, error) {
	return lightDisplay(r, brightnessLights)
}
