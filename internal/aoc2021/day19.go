package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func parseScanners(r io.Reader) ([][]grid.Position3D[int64], error) {
	s := bufio.NewScanner(r)
	var scanners [][]grid.Position3D[int64]
	ids := make(map[int]bool)
	var seen map[grid.Position3D[int64]]bool
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		if id, ok := strings.CutPrefix(line, "--- scanner "); ok {
			id, ok = strings.CutSuffix(id, " ---")
			if !ok {
				return nil, fmt.Errorf("invalid scanner header %q", line)
			}
			n, err := convert.ScanNumber[int]([]byte(id))
			if err != nil || n < 0 || ids[n] {
				return nil, fmt.Errorf("invalid or duplicate scanner ID %q", id)
			}
			if len(scanners) == 128 || len(scanners) > 0 && len(scanners[len(scanners)-1]) == 0 {
				return nil, fmt.Errorf("empty scanner or more than 128 scanners")
			}
			ids[n] = true
			scanners = append(scanners, nil)
			seen = make(map[grid.Position3D[int64]]bool)
			continue
		}
		if len(scanners) == 0 {
			return nil, fmt.Errorf("beacon before first scanner header")
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid beacon %q", line)
		}
		var coordinates [3]int64
		for i, part := range parts {
			value, err := convert.ScanNumber[int32]([]byte(part))
			if err != nil {
				return nil, fmt.Errorf("beacon coordinate: %w", err)
			}
			coordinates[i] = int64(value)
		}
		point := grid.NewPosition3D(coordinates[0], coordinates[1], coordinates[2])
		if seen[point] {
			return nil, fmt.Errorf("duplicate beacon %q", line)
		}
		seen[point] = true
		last := len(scanners) - 1
		if len(scanners[last]) == 512 {
			return nil, fmt.Errorf("scanner exceeds 512 beacons")
		}
		scanners[last] = append(scanners[last], point)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(scanners) == 0 || len(scanners[len(scanners)-1]) == 0 {
		return nil, fmt.Errorf("missing scanner beacons")
	}
	return scanners, nil
}

func scannerRotations() [24][3]int {
	var rotations [24][3]int
	count := 0
	for _, axes := range [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}} {
		parity := 1
		for i := range 3 {
			for j := i + 1; j < 3; j++ {
				if axes[i] > axes[j] {
					parity = -parity
				}
			}
		}
		for _, x := range []int{-1, 1} {
			for _, y := range []int{-1, 1} {
				// Positive determinant excludes reflected coordinate systems.
				z := parity * x * y
				rotations[count] = [3]int{x * (axes[0] + 1), y * (axes[1] + 1), z * (axes[2] + 1)}
				count++
			}
		}
	}
	return rotations
}

func rotateBeacon(point grid.Position3D[int64], rotation [3]int) grid.Position3D[int64] {
	coordinates := [3]int64{point.X, point.Y, point.Z}
	var result [3]int64
	for i, axis := range rotation {
		if axis < 0 {
			result[i] = -coordinates[-axis-1]
		} else {
			result[i] = coordinates[axis-1]
		}
	}
	return grid.NewPosition3D(result[0], result[1], result[2])
}

func locateScanner(known, scanner []grid.Position3D[int64], work *int) (grid.Position3D[int64], []grid.Position3D[int64], bool, error) {
	if len(known) < 12 || len(scanner) < 12 {
		return grid.Position3D[int64]{}, nil, false, nil
	}
	if len(known)*len(scanner) > 2_000_000 {
		return grid.Position3D[int64]{}, nil, false, fmt.Errorf("scanner alignment exceeds 2000000 translation candidates")
	}
	rotated := make([]grid.Position3D[int64], len(scanner))
	for _, rotation := range scannerRotations() {
		for i, point := range scanner {
			rotated[i] = rotateBeacon(point, rotation)
		}
		translations := make(map[grid.Position3D[int64]]int)
		for _, fixed := range known {
			for _, point := range rotated {
				*work++
				if *work > 100_000_000 {
					return grid.Position3D[int64]{}, nil, false, fmt.Errorf("scanner alignment exceeds 100000000 comparisons")
				}
				offset := fixed.Sub(point)
				translations[offset]++
				if translations[offset] == 12 {
					for i, beacon := range rotated {
						rotated[i] = beacon.Add(offset)
					}
					return offset, rotated, true, nil
				}
			}
		}
	}
	return grid.Position3D[int64]{}, nil, false, nil
}

func alignScanners(r io.Reader) (int, int64, error) {
	scanners, err := parseScanners(r)
	if err != nil {
		return 0, 0, err
	}
	known := scanners[0]
	beacons := make(map[grid.Position3D[int64]]bool, len(known))
	for _, point := range known {
		beacons[point] = true
	}
	positions := []grid.Position3D[int64]{{}}
	aligned := make([]bool, len(scanners))
	aligned[0] = true
	work := 0
	for len(positions) < len(scanners) {
		progress := false
		for i, scanner := range scanners {
			if aligned[i] {
				continue
			}
			position, moved, found, err := locateScanner(known, scanner, &work)
			if err != nil {
				return 0, 0, err
			}
			if !found {
				continue
			}
			for _, beacon := range moved {
				if !beacons[beacon] {
					beacons[beacon] = true
					known = append(known, beacon)
				}
			}
			positions = append(positions, position)
			aligned[i], progress = true, true
		}
		if !progress {
			return 0, 0, fmt.Errorf("cannot align %d scanners with 12 distinct overlapping beacons", len(scanners)-len(positions))
		}
	}
	var distance int64
	for i, left := range positions {
		for _, right := range positions[i+1:] {
			distance = max(distance, left.Distance(right))
		}
	}
	return len(known), distance, nil
}

func day19p01(r io.Reader) (string, error) {
	count, _, err := alignScanners(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(count), nil
}

func day19p02(r io.Reader) (string, error) {
	_, distance, err := alignScanners(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(distance, 10), nil
}
