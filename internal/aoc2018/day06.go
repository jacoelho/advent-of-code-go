package aoc2018

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"slices"
	"sort"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func coordinateSites(r io.Reader) ([]grid.Position2D[int], error) {
	s := bufio.NewScanner(r)
	var sites []grid.Position2D[int]
	for s.Scan() {
		xText, yText, ok := bytes.Cut(s.Bytes(), []byte(","))
		if !ok {
			return nil, fmt.Errorf("invalid coordinate %q", s.Text())
		}
		x, err := convert.ScanNumber[int](bytes.TrimSpace(xText))
		if err != nil {
			return nil, err
		}
		y, err := convert.ScanNumber[int](bytes.TrimSpace(yText))
		if err != nil {
			return nil, err
		}
		sites = append(sites, grid.Position2D[int]{X: x, Y: y})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(sites) == 0 {
		return nil, fmt.Errorf("no coordinate sites")
	}
	return sites, nil
}

func siteBounds(sites []grid.Position2D[int]) (grid.Position2D[int], grid.Position2D[int]) {
	lower, upper := sites[0], sites[0]
	for _, site := range sites[1:] {
		lower.X, lower.Y = min(lower.X, site.X), min(lower.Y, site.Y)
		upper.X, upper.Y = max(upper.X, site.X), max(upper.Y, site.Y)
	}
	return lower, upper
}

func day06p01(r io.Reader) (string, error) {
	sites, err := coordinateSites(r)
	if err != nil {
		return "", err
	}

	lower, upper := siteBounds(sites)
	areas := make([]int, len(sites))
	infinite := make([]bool, len(sites))
	for y := lower.Y; y <= upper.Y; y++ {
		for x := lower.X; x <= upper.X; x++ {
			point := grid.Position2D[int]{X: x, Y: y}
			owner, distance := -1, math.MaxInt
			for i, site := range sites {
				d := point.Distance(site)
				if d < distance {
					owner, distance = i, d
				} else if d == distance {
					owner = -1
				}
			}
			if owner == -1 {
				continue
			}
			areas[owner]++
			if x == lower.X || x == upper.X || y == lower.Y || y == upper.Y {
				infinite[owner] = true
			}
		}
	}
	largest := 0
	for site, area := range areas {
		if !infinite[site] {
			largest = max(largest, area)
		}
	}
	return strconv.Itoa(largest), nil
}

func safeRegionSize(r io.Reader, limit int) (int, error) {
	sites, err := coordinateSites(r)
	if err != nil {
		return 0, err
	}
	if limit <= 0 {
		return 0, nil
	}

	lower, upper := siteBounds(sites)
	// Outside this margin, the distance along one axis alone reaches the limit.
	margin := (limit - 1) / len(sites)
	yDistances := make([]int, 0, upper.Y-lower.Y+2*margin+1)
	for y := lower.Y - margin; y <= upper.Y+margin; y++ {
		distance := 0
		for _, site := range sites {
			distance += xmath.Abs(y - site.Y)
		}
		yDistances = append(yDistances, distance)
	}
	slices.Sort(yDistances)

	area := 0
	for x := lower.X - margin; x <= upper.X+margin; x++ {
		distance := 0
		for _, site := range sites {
			distance += xmath.Abs(x - site.X)
		}
		area += sort.SearchInts(yDistances, limit-distance)
	}
	return area, nil
}

func day06p02(r io.Reader) (string, error) {
	area, err := safeRegionSize(r, 10000)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(area), nil
}
