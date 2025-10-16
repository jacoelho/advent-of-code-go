package aoc2023

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type Hailstone struct {
	position grid.Position3D[int]
	velocity grid.Position3D[int]
}

func (h Hailstone) zAt(time float64, rockVZ int) float64 {
	return float64(h.position.Z) + time*float64(h.velocity.Z+rockVZ)
}

func parsePosition3D(s string) (grid.Position3D[int], error) {
	digits := convert.ExtractDigits[int](s)
	if len(digits) != 3 {
		return grid.Position3D[int]{}, fmt.Errorf("invalid coordinates")
	}
	return grid.NewPosition3D(digits[0], digits[1], digits[2]), nil
}

func parseHailstone(line []byte) (Hailstone, error) {
	parts := strings.Split(string(line), "@")
	if len(parts) != 2 {
		return Hailstone{}, fmt.Errorf("invalid hailstone format")
	}

	position, err := parsePosition3D(parts[0])
	if err != nil {
		return Hailstone{}, err
	}

	velocity, err := parsePosition3D(parts[1])
	if err != nil {
		return Hailstone{}, err
	}

	return Hailstone{position: position, velocity: velocity}, nil
}

type intersection struct {
	x, y   float64
	t1, t2 float64
}

func calculateIntersection2D(h1, h2 Hailstone) *intersection {
	if h1.velocity.X == 0 || h2.velocity.X == 0 {
		return nil
	}

	slope1 := float64(h1.velocity.Y) / float64(h1.velocity.X)
	slope2 := float64(h2.velocity.Y) / float64(h2.velocity.X)

	if slope1 == slope2 {
		return nil
	}

	c1 := float64(h1.position.Y) - slope1*float64(h1.position.X)
	c2 := float64(h2.position.Y) - slope2*float64(h2.position.X)

	x := (c2 - c1) / (slope1 - slope2)
	y := slope1*x + c1

	t1 := (x - float64(h1.position.X)) / float64(h1.velocity.X)
	t2 := (x - float64(h2.position.X)) / float64(h2.velocity.X)

	if t1 < 0 || t2 < 0 {
		return nil
	}

	return &intersection{x: x, y: y, t1: t1, t2: t2}
}

func day24p01(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseHailstone)
	hailstones := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return "", err
	}

	count := xiter.CountBy(func(pair xslices.Pair[Hailstone, Hailstone]) bool {
		return intersects2D(pair.V1, pair.V2, 200000000000000, 400000000000000)
	}, xslices.Pairwise(hailstones))

	return strconv.Itoa(count), nil
}

func intersects2D(h1, h2 Hailstone, minBound, maxBound float64) bool {
	inter := calculateIntersection2D(h1, h2)
	if inter == nil {
		return false
	}

	return inter.x >= minBound && inter.x <= maxBound &&
		inter.y >= minBound && inter.y <= maxBound
}

func allIntersectionsMatch(intersections []*intersection) bool {
	if len(intersections) == 0 {
		return false
	}

	return xslices.Every(func(inter *intersection) bool {
		return inter.x == intersections[0].x && inter.y == intersections[0].y
	}, intersections[1:])
}

func findRockVelocityZ(hailstones []Hailstone, intersections []*intersection) (int, bool) {
	const searchRange = 1000
	for rockVZ := -searchRange; rockVZ <= searchRange; rockVZ++ {
		z1 := hailstones[1].zAt(intersections[0].t1, rockVZ)
		z2 := hailstones[2].zAt(intersections[1].t1, rockVZ)
		z3 := hailstones[3].zAt(intersections[2].t1, rockVZ)

		if z1 == z2 && z2 == z3 {
			return rockVZ, true
		}
	}
	return 0, false
}

func day24p02(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseHailstone)
	hailstones := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(hailstones) < 4 {
		return "", fmt.Errorf("need at least 4 hailstones")
	}

	const searchRange = 1000
	for rockVX := -searchRange; rockVX <= searchRange; rockVX++ {
		for rockVY := -searchRange; rockVY <= searchRange; rockVY++ {
			adjustedStone := Hailstone{
				position: hailstones[0].position,
				velocity: hailstones[0].velocity.Add(grid.NewPosition3D(rockVX, rockVY, 0)),
			}
			var intercepts []*intersection

			for i := 1; i < 4; i++ {
				h := Hailstone{
					position: hailstones[i].position,
					velocity: hailstones[i].velocity.Add(grid.NewPosition3D(rockVX, rockVY, 0)),
				}
				inter := calculateIntersection2D(h, adjustedStone)
				if inter == nil {
					intercepts = nil
					break
				}
				intercepts = append(intercepts, inter)
			}

			if len(intercepts) != 3 {
				continue
			}

			if !allIntersectionsMatch(intercepts) {
				continue
			}

			if rockVZ, found := findRockVelocityZ(hailstones, intercepts); found {
				z1 := hailstones[1].zAt(intercepts[0].t1, rockVZ)
				result := int(intercepts[0].x + intercepts[0].y + z1)
				return strconv.Itoa(result), nil
			}
		}
	}

	return "", fmt.Errorf("no solution found")
}
