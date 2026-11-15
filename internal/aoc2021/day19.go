package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

const sharedBeacons = 12

func parseScanners(r io.Reader) ([][]grid.Position3D[int64], error) {
	s := bufio.NewScanner(r)
	var scanners [][]grid.Position3D[int64]
	ids := make(map[int]bool)
	var seen map[grid.Position3D[int64]]bool
	separated := false
	for s.Scan() {
		line := s.Text()
		if line == "" {
			separated = true
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
			if len(scanners) > 0 && len(scanners[len(scanners)-1]) == 0 {
				return nil, fmt.Errorf("scanner %d has no beacons", len(scanners)-1)
			}
			ids[n] = true
			scanners = append(scanners, nil)
			seen = make(map[grid.Position3D[int64]]bool)
			separated = false
			continue
		}
		if len(scanners) == 0 || separated {
			return nil, fmt.Errorf("beacon %q outside a scanner section", line)
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

// beaconSpan is the sorted absolute per-axis delta between two beacons of one
// scanner, which no rotation or translation of that scanner changes.
type beaconSpan [3]int64

// scannerFingerprint returns the sorted multiset of spans between all beacon
// pairs of a scanner.
func scannerFingerprint(beacons []grid.Position3D[int64]) []beaconSpan {
	spans := make([]beaconSpan, 0, len(beacons)*(len(beacons)-1)/2)
	for i, a := range beacons {
		for _, b := range beacons[i+1:] {
			span := beaconSpan{xmath.Abs(a.X - b.X), xmath.Abs(a.Y - b.Y), xmath.Abs(a.Z - b.Z)}
			slices.Sort(span[:])
			spans = append(spans, span)
		}
	}
	slices.SortFunc(spans, func(a, b beaconSpan) int { return slices.Compare(a[:], b[:]) })
	return spans
}

// sharedSpans counts the multiset intersection of two fingerprints. Scanners
// sharing n beacons share at least n(n-1)/2 spans, so a lower count rules out
// an overlap without trying any rotation.
func sharedSpans(a, b []beaconSpan) int {
	shared := 0
	for len(a) > 0 && len(b) > 0 {
		switch c := slices.Compare(a[0][:], b[0][:]); {
		case c < 0:
			a = a[1:]
		case c > 0:
			b = b[1:]
		default:
			shared++
			a, b = a[1:], b[1:]
		}
	}
	return shared
}

// locateScanner finds the rotation and translation that place at least
// sharedBeacons of scanner onto known, returning the scanner position and its
// beacons in the frame of known.
func locateScanner(known, scanner []grid.Position3D[int64]) (grid.Position3D[int64], []grid.Position3D[int64], bool) {
	rotated := make([]grid.Position3D[int64], len(scanner))
	translations := make(map[grid.Position3D[int64]]int)
	for _, rotation := range scannerRotations() {
		for i, point := range scanner {
			rotated[i] = rotateBeacon(point, rotation)
		}
		clear(translations)
		for _, fixed := range known {
			for _, point := range rotated {
				offset := fixed.Sub(point)
				translations[offset]++
				if translations[offset] == sharedBeacons {
					for i, beacon := range rotated {
						rotated[i] = beacon.Add(offset)
					}
					return offset, rotated, true
				}
			}
		}
	}
	return grid.Position3D[int64]{}, nil, false
}

func alignScanners(r io.Reader) (int, int64, error) {
	scanners, err := parseScanners(r)
	if err != nil {
		return 0, 0, err
	}
	fingerprints := make([][]beaconSpan, len(scanners))
	for i, beacons := range scanners {
		fingerprints[i] = scannerFingerprint(beacons)
	}
	placed := make([][]grid.Position3D[int64], len(scanners))
	positions := make([]grid.Position3D[int64], len(scanners))
	placed[0] = scanners[0]
	queue := []int{0}
	for head := 0; head < len(queue); head++ {
		anchor := queue[head]
		for i, scanner := range scanners {
			if placed[i] != nil || sharedSpans(fingerprints[anchor], fingerprints[i]) < sharedBeacons*(sharedBeacons-1)/2 {
				continue
			}
			position, moved, found := locateScanner(placed[anchor], scanner)
			if !found {
				continue
			}
			placed[i], positions[i] = moved, position
			queue = append(queue, i)
		}
	}
	if len(queue) != len(scanners) {
		return 0, 0, fmt.Errorf("cannot align %d scanners: each needs %d beacons in common with an aligned scanner", len(scanners)-len(queue), sharedBeacons)
	}
	beacons := make(map[grid.Position3D[int64]]struct{})
	for _, moved := range placed {
		for _, beacon := range moved {
			beacons[beacon] = struct{}{}
		}
	}
	var distance int64
	for i, left := range positions {
		for _, right := range positions[i+1:] {
			distance = max(distance, left.Distance(right))
		}
	}
	return len(beacons), distance, nil
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
