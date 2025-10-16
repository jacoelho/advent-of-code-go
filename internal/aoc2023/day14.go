package aoc2023

import (
	"io"
	"iter"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type direction int

const (
	tiltNorth direction = iota
	tiltWest
	tiltSouth
	tiltEast
)

type axis int

const (
	vertical axis = iota
	horizontal
)

type directionType int

const (
	forward  directionType = iota // north/west (min to max)
	backward                      // south/east (max to min)
)

type rollDirection struct {
	axis          axis
	directionType directionType
}

func parsePlatform(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(line []byte) ([]rune, error) {
		return []rune(string(line)), nil
	})

	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}

	return grid.NewGrid2D[int](rows), nil
}

func tilt(platform grid.Grid2D[int, rune], dir direction) grid.Grid2D[int, rune] {
	result := maps.Clone(platform)
	minX, maxX, minY, maxY := platform.Dimensions()

	params := map[direction]rollDirection{
		tiltNorth: {vertical, forward},
		tiltSouth: {vertical, backward},
		tiltWest:  {horizontal, forward},
		tiltEast:  {horizontal, backward},
	}[dir]

	if params.axis == vertical {
		// North/South: iterate over columns
		for x := minX; x <= maxX; x++ {
			var nextY int
			if params.directionType == backward {
				nextY = maxY
			} else {
				nextY = minY
			}

			var yRange func() iter.Seq[int]
			if params.directionType == backward {
				yRange = func() iter.Seq[int] {
					return func(yield func(int) bool) {
						for y := maxY; y >= minY; y-- {
							if !yield(y) {
								return
							}
						}
					}
				}
			} else {
				yRange = func() iter.Seq[int] {
					return func(yield func(int) bool) {
						for y := minY; y <= maxY; y++ {
							if !yield(y) {
								return
							}
						}
					}
				}
			}

			for y := range yRange() {
				pos := grid.NewPosition2D(x, y)
				val := platform[pos]
				switch val {
				case '#':
					if params.directionType == backward {
						nextY = y - 1
					} else {
						nextY = y + 1
					}
				case 'O':
					if y != nextY {
						result[pos] = '.'
						result[grid.NewPosition2D(x, nextY)] = 'O'
					}
					if params.directionType == backward {
						nextY--
					} else {
						nextY++
					}
				}
			}
		}
	} else {
		for y := minY; y <= maxY; y++ {
			var nextX int
			if params.directionType == backward {
				nextX = maxX
			} else {
				nextX = minX
			}

			var xRange func() iter.Seq[int]
			if params.directionType == backward {
				xRange = func() iter.Seq[int] {
					return func(yield func(int) bool) {
						for x := maxX; x >= minX; x-- {
							if !yield(x) {
								return
							}
						}
					}
				}
			} else {
				xRange = func() iter.Seq[int] {
					return func(yield func(int) bool) {
						for x := minX; x <= maxX; x++ {
							if !yield(x) {
								return
							}
						}
					}
				}
			}

			for x := range xRange() {
				pos := grid.NewPosition2D(x, y)
				val := platform[pos]
				switch val {
				case '#':
					if params.directionType == backward {
						nextX = x - 1
					} else {
						nextX = x + 1
					}
				case 'O':
					if x != nextX {
						result[pos] = '.'
						result[grid.NewPosition2D(nextX, y)] = 'O'
					}
					if params.directionType == backward {
						nextX--
					} else {
						nextX++
					}
				}
			}
		}
	}

	return result
}

func spinCycle(platform grid.Grid2D[int, rune]) grid.Grid2D[int, rune] {
	platform = tilt(platform, tiltNorth)
	platform = tilt(platform, tiltWest)
	platform = tilt(platform, tiltSouth)
	platform = tilt(platform, tiltEast)
	return platform
}

func calculateLoad(g grid.Grid2D[int, rune]) int {
	_, _, _, maxY := g.Dimensions()
	load := 0
	for pos, val := range g {
		if val == 'O' {
			load += maxY - pos.Y + 1
		}
	}
	return load
}

func platformState(platform grid.Grid2D[int, rune]) string {
	minX, maxX, minY, maxY := platform.Dimensions()
	var sb strings.Builder
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			sb.WriteRune(platform[grid.NewPosition2D(x, y)])
		}
		if y < maxY {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func day14p01(r io.Reader) (string, error) {
	platform, err := parsePlatform(r)
	if err != nil {
		return "", err
	}
	platform = tilt(platform, tiltNorth)
	load := calculateLoad(platform)
	return strconv.Itoa(load), nil
}

func day14p02(r io.Reader) (string, error) {
	platform, err := parsePlatform(r)
	if err != nil {
		return "", err
	}

	seen := make(map[string]int)
	cycles := 1000000000

	for i := range cycles {
		platform = spinCycle(platform)
		key := platformState(platform)

		if prevI, found := seen[key]; found {
			cycleLength := i - prevI
			remaining := cycles - i - 1
			finalPos := remaining % cycleLength

			for j := 0; j < finalPos; j++ {
				platform = spinCycle(platform)
			}
			break
		}
		seen[key] = i
	}

	load := calculateLoad(platform)
	return strconv.Itoa(load), nil
}
