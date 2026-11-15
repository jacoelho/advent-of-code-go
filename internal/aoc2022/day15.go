package aoc2022

import (
	"cmp"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type beaconSensor struct {
	position grid.Position2D[int64]
	beacon   grid.Position2D[int64]
	radius   int64
}

type sensorInterval struct{ start, end int64 }

func parseSensors(r io.Reader) ([]beaconSensor, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return nil, err
	}
	pattern := regexp.MustCompile(`^Sensor at x=(-?[0-9]+), y=(-?[0-9]+): closest beacon is at x=(-?[0-9]+), y=(-?[0-9]+)$`)
	var sensors []beaconSensor
	for line := range strings.SplitSeq(input, "\n") {
		fields := pattern.FindStringSubmatch(line)
		if fields == nil {
			return nil, fmt.Errorf("invalid sensor record")
		}
		var values [4]int64
		for i := range values {
			value, err := convert.ScanNumber[int32]([]byte(fields[i+1]))
			if err != nil {
				return nil, fmt.Errorf("sensor coordinate: %w", err)
			}
			values[i] = int64(value)
		}
		position := grid.Position2D[int64]{X: values[0], Y: values[1]}
		beacon := grid.Position2D[int64]{X: values[2], Y: values[3]}
		sensors = append(sensors, beaconSensor{position: position, beacon: beacon, radius: position.Distance(beacon)})
	}
	return sensors, nil
}

func sensorCoverage(sensors []beaconSensor, row int64, intervals []sensorInterval) []sensorInterval {
	intervals = intervals[:0]
	for _, sensor := range sensors {
		reach := sensor.radius - xmath.Abs(sensor.position.Y-row)
		if reach >= 0 {
			intervals = append(intervals, sensorInterval{sensor.position.X - reach, sensor.position.X + reach})
		}
	}
	slices.SortFunc(intervals, func(a, b sensorInterval) int { return cmp.Compare(a.start, b.start) })
	merged := intervals[:0]
	for _, interval := range intervals {
		if len(merged) == 0 || interval.start > merged[len(merged)-1].end+1 {
			merged = append(merged, interval)
		} else {
			merged[len(merged)-1].end = max(merged[len(merged)-1].end, interval.end)
		}
	}
	return merged
}

func excludedBeacons(r io.Reader, row int64) (string, error) {
	sensors, err := parseSensors(r)
	if err != nil {
		return "", err
	}
	intervals := sensorCoverage(sensors, row, make([]sensorInterval, 0, len(sensors)))
	var total int64
	for _, interval := range intervals {
		total += interval.end - interval.start + 1
	}
	beacons := collections.NewSet[int64]()
	for _, sensor := range sensors {
		if sensor.beacon.Y == row {
			beacons.Add(sensor.beacon.X)
		}
	}
	return strconv.FormatInt(total-int64(beacons.Len()), 10), nil
}

// distressBeacon works in rotated coordinates u = x+y, v = x-y, where every
// sensor diamond becomes an axis-aligned square. The square edges split the
// plane into cells that are either wholly covered or wholly uncovered, so only
// uncovered cells need their lattice points (u and v of equal parity) examined.
func distressBeacon(r io.Reader, limit int64) (string, error) {
	if limit < 0 {
		return "", fmt.Errorf("negative search limit %d", limit)
	}
	sensors, err := parseSensors(r)
	if err != nil {
		return "", err
	}
	us := []int64{0, 2*limit + 1}
	vs := []int64{-limit, limit + 1}
	for _, sensor := range sensors {
		u, v := sensor.position.X+sensor.position.Y, sensor.position.X-sensor.position.Y
		us = append(us, min(max(u-sensor.radius, 0), 2*limit+1), min(max(u+sensor.radius+1, 0), 2*limit+1))
		vs = append(vs, min(max(v-sensor.radius, -limit), limit+1), min(max(v+sensor.radius+1, -limit), limit+1))
	}
	slices.Sort(us)
	slices.Sort(vs)
	us, vs = slices.Compact(us), slices.Compact(vs)
	var found []grid.Position2D[int64]
	for i := range len(us) - 1 {
		for j := range len(vs) - 1 {
			u1, u2, v1, v2 := us[i], us[i+1]-1, vs[j], vs[j+1]-1
			covered := slices.ContainsFunc(sensors, func(sensor beaconSensor) bool {
				u, v := sensor.position.X+sensor.position.Y, sensor.position.X-sensor.position.Y
				return xmath.Abs(u1-u) <= sensor.radius && xmath.Abs(v1-v) <= sensor.radius
			})
			if covered {
				continue
			}
			// Inside the search square |v| <= u <= 2*limit-|v|; this loop finds a
			// lattice point within two iterations, so it stops early once two exist.
			nearest := max(0, v1, -v2)
			for u := max(u1, nearest); u <= min(u2, 2*limit-nearest); u++ {
				reach := min(u, 2*limit-u)
				low, high := max(v1, -reach), min(v2, reach)
				low += (low - u) & 1
				for v := low; v <= high; v += 2 {
					found = append(found, grid.Position2D[int64]{X: (u + v) / 2, Y: (u - v) / 2})
					if len(found) > 1 {
						return "", fmt.Errorf("distress beacon is not unique")
					}
				}
			}
		}
	}
	if len(found) == 0 {
		return "", fmt.Errorf("no uncovered beacon position")
	}
	return strconv.FormatInt(found[0].X*4_000_000+found[0].Y, 10), nil
}

func day15p01(r io.Reader) (string, error) {
	return excludedBeacons(r, 2_000_000)
}

func day15p02(r io.Reader) (string, error) {
	return distressBeacon(r, 4_000_000)
}
