package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type lightMode int

const (
	onOffLights lightMode = iota
	brightnessLights
)

func lightDisplay(r io.Reader, mode lightMode) (string, error) {
	grid := make([]int, 1000*1000)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		operation, coordinates := "", ""
		for _, candidate := range []string{"turn on", "turn off", "toggle"} {
			if rest, ok := strings.CutPrefix(line, candidate+" "); ok {
				operation, coordinates = candidate, rest
				break
			}
		}
		var x1, y1, x2, y2 int
		n, err := fmt.Sscanf(coordinates, "%d,%d through %d,%d", &x1, &y1, &x2, &y2)
		if operation == "" || err != nil || n != 4 || x1 < 0 || y1 < 0 || x2 >= 1000 || y2 >= 1000 || x1 > x2 || y1 > y2 {
			return "", fmt.Errorf("invalid lighting instruction %q", line)
		}
		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				light := &grid[y*1000+x]
				if mode == onOffLights {
					switch operation {
					case "turn on":
						*light = 1
					case "turn off":
						*light = 0
					case "toggle":
						*light = 1 - *light
					}
				} else {
					switch operation {
					case "turn on":
						*light++
					case "turn off":
						*light = max(0, *light-1)
					case "toggle":
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
	for _, light := range grid {
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
