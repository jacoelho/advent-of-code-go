package aoc2022

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func lavaCubes(r io.Reader) (map[grid.Position3D[int64]]bool, error) {
	s := scanner.NewScanner(r, func(line []byte) (grid.Position3D[int64], error) {
		parts := bytes.Split(line, []byte(","))
		if len(parts) != 3 {
			return grid.Position3D[int64]{}, fmt.Errorf("invalid lava cube %q", line)
		}
		var coordinates [3]int64
		for i, part := range parts {
			value, err := convert.ScanNumber[int32](part)
			if err != nil {
				return grid.Position3D[int64]{}, fmt.Errorf("lava coordinate: %w", err)
			}
			coordinates[i] = int64(value)
		}
		return grid.NewPosition3D(coordinates[0], coordinates[1], coordinates[2]), nil
	})
	cubes := make(map[grid.Position3D[int64]]bool)
	count := 0
	for cube := range s.Values() {
		count++
		if count > 100_000 {
			return nil, fmt.Errorf("lava input exceeds 100000 cubes")
		}
		cubes[cube] = true
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(cubes) == 0 {
		return nil, fmt.Errorf("empty lava droplet")
	}
	return cubes, nil
}

func lavaNeighbours(point grid.Position3D[int64]) [6]grid.Position3D[int64] {
	return [6]grid.Position3D[int64]{
		{X: point.X - 1, Y: point.Y, Z: point.Z},
		{X: point.X + 1, Y: point.Y, Z: point.Z},
		{X: point.X, Y: point.Y - 1, Z: point.Z},
		{X: point.X, Y: point.Y + 1, Z: point.Z},
		{X: point.X, Y: point.Y, Z: point.Z - 1},
		{X: point.X, Y: point.Y, Z: point.Z + 1},
	}
}

func day18p01(r io.Reader) (string, error) {
	cubes, err := lavaCubes(r)
	if err != nil {
		return "", err
	}
	area := 0
	for cube := range cubes {
		for _, next := range lavaNeighbours(cube) {
			if !cubes[next] {
				area++
			}
		}
	}
	return strconv.Itoa(area), nil
}

func day18p02(r io.Reader) (string, error) {
	cubes, err := lavaCubes(r)
	if err != nil {
		return "", err
	}
	var low, high grid.Position3D[int64]
	first := true
	for cube := range cubes {
		if first {
			low, high, first = cube, cube, false
		}
		low.X, high.X = min(low.X, cube.X), max(high.X, cube.X)
		low.Y, high.Y = min(low.Y, cube.Y), max(high.Y, cube.Y)
		low.Z, high.Z = min(low.Z, cube.Z), max(high.Z, cube.Z)
	}
	low = low.Sub(grid.NewPosition3D[int64](1, 1, 1))
	high = high.Add(grid.NewPosition3D[int64](1, 1, 1))
	width, height, depth := high.X-low.X+1, high.Y-low.Y+1, high.Z-low.Z+1
	volume := int64(1)
	for _, size := range []int64{width, height, depth} {
		if size > 4_000_000/volume {
			return "", fmt.Errorf("lava exterior exceeds 4000000 cells")
		}
		volume *= size
	}
	index := func(point grid.Position3D[int64]) int64 {
		return ((point.Z-low.Z)*height+point.Y-low.Y)*width + point.X - low.X
	}
	visited := make([]bool, volume)
	visited[0] = true
	queue := collections.NewDeque[grid.Position3D[int64]](1024)
	queue.PushBack(low)
	area := 0
	for queue.Size() > 0 {
		point, _ := queue.PopFront()
		for _, next := range lavaNeighbours(point) {
			if next.X < low.X || next.X > high.X || next.Y < low.Y || next.Y > high.Y || next.Z < low.Z || next.Z > high.Z {
				continue
			}
			if cubes[next] {
				area++
				continue
			}
			i := index(next)
			if !visited[i] {
				visited[i] = true
				queue.PushBack(next)
			}
		}
	}
	return strconv.Itoa(area), nil
}
