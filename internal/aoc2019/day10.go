package aoc2019

import (
	"cmp"
	"io"
	"maps"
	"math"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmaps"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseAsteroidMap(r io.Reader) ([]grid.Position2D[int], error) {
	s := scanner.NewScanner(r, func(line []byte) ([]rune, error) {
		return []rune(string(line)), nil
	})

	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}

	g := grid.NewGrid2D[int](rows)
	asteroidPairs := xmaps.Filter(func(_ grid.Position2D[int], char rune) bool {
		return char == '#'
	}, g)

	return xslices.Map(func(p xmaps.Pair[grid.Position2D[int], rune]) grid.Position2D[int] {
		return p.K
	}, asteroidPairs), nil
}

func normalizeDirection(station, asteroid grid.Position2D[int]) grid.Position2D[int] {
	dx := asteroid.X - station.X
	dy := asteroid.Y - station.Y
	gcd := xmath.GCD(dx, dy)

	return grid.Position2D[int]{
		X: dx / gcd,
		Y: dy / gcd,
	}
}

func countVisibleAsteroids(station grid.Position2D[int], asteroids []grid.Position2D[int]) int {
	directions := collections.NewSet[grid.Position2D[int]]()

	for _, asteroid := range asteroids {
		if asteroid == station {
			continue
		}

		directions.Add(normalizeDirection(station, asteroid))
	}

	return directions.Len()
}

func findBestLocation(asteroids []grid.Position2D[int]) (grid.Position2D[int], int) {
	bestLocation := xslices.MaxBy(func(a, b grid.Position2D[int]) bool {
		return countVisibleAsteroids(a, asteroids) < countVisibleAsteroids(b, asteroids)
	}, asteroids)

	return bestLocation, countVisibleAsteroids(bestLocation, asteroids)
}

func day10p01(r io.Reader) (string, error) {
	asteroids, err := parseAsteroidMap(r)
	if err != nil {
		return "", err
	}

	_, maxVisible := findBestLocation(asteroids)
	return strconv.Itoa(maxVisible), nil
}

func calculateAngle(direction grid.Position2D[int]) float64 {
	angle := math.Atan2(float64(direction.X), float64(-direction.Y))
	if angle < 0 {
		angle += 2 * math.Pi
	}
	return angle
}

func groupAsteroidsByDirection(
	station grid.Position2D[int],
	asteroids []grid.Position2D[int],
) map[grid.Position2D[int]][]grid.Position2D[int] {
	groups := make(map[grid.Position2D[int]][]grid.Position2D[int])

	for _, asteroid := range asteroids {
		if asteroid == station {
			continue
		}

		normalizedDir := normalizeDirection(station, asteroid)
		groups[normalizedDir] = append(groups[normalizedDir], asteroid)
	}

	for dir := range groups {
		slices.SortFunc(groups[dir], func(a, b grid.Position2D[int]) int {
			distA := station.Distance(a)
			distB := station.Distance(b)
			return cmp.Compare(distA, distB)
		})
	}

	return groups
}

func vaporizeAsteroids(
	station grid.Position2D[int],
	asteroids []grid.Position2D[int],
	targetCount int,
) grid.Position2D[int] {
	groups := groupAsteroidsByDirection(station, asteroids)

	directions := slices.Collect(maps.Keys(groups))
	slices.SortFunc(directions, func(a, b grid.Position2D[int]) int {
		return cmp.Compare(calculateAngle(a), calculateAngle(b))
	})

	count := 0
	for {
		for _, dir := range directions {
			if len(groups[dir]) == 0 {
				continue
			}

			vaporized := groups[dir][0]
			groups[dir] = groups[dir][1:]
			count++

			if count == targetCount {
				return vaporized
			}
		}
	}
}

func day10p02(r io.Reader) (string, error) {
	asteroids, err := parseAsteroidMap(r)
	if err != nil {
		return "", err
	}

	station, _ := findBestLocation(asteroids)
	target := vaporizeAsteroids(station, asteroids, 200)

	result := target.X*100 + target.Y
	return strconv.Itoa(result), nil
}
