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
			value, err := convert.ScanNumber[int64]([]byte(fields[i+1]))
			if err != nil {
				return nil, err
			}
			if value < -1_000_000_000 || value > 1_000_000_000 {
				return nil, fmt.Errorf("sensor coordinate exceeds one billion")
			}
			values[i] = value
		}
		position := grid.Position2D[int64]{X: values[0], Y: values[1]}
		beacon := grid.Position2D[int64]{X: values[2], Y: values[3]}
		sensors = append(sensors, beaconSensor{position: position, beacon: beacon, radius: position.Distance(beacon)})
		if len(sensors) > 128 {
			return nil, fmt.Errorf("too many sensors")
		}
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
	if row < -1_000_000_000 || row > 1_000_000_000 {
		return "", fmt.Errorf("scan row exceeds coordinate bound")
	}
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

func distressBeacon(r io.Reader, limit int64) (string, error) {
	if limit < 0 || limit > 4_000_000 {
		return "", fmt.Errorf("search limit must be between zero and four million")
	}
	sensors, err := parseSensors(r)
	if err != nil {
		return "", err
	}
	if int64(len(sensors))*(limit+1) > 160_000_000 {
		return "", fmt.Errorf("sensor search exceeds work limit")
	}
	intervals := make([]sensorInterval, 0, len(sensors))
	var answer grid.Position2D[int64]
	found := false
	for y := int64(0); y <= limit; y++ {
		intervals = sensorCoverage(sensors, y, intervals)
		next := int64(0)
		gap := func(end int64) error {
			if next > end {
				return nil
			}
			if next != end || found {
				return fmt.Errorf("distress beacon is not unique")
			}
			answer = grid.Position2D[int64]{X: next, Y: y}
			found = true
			return nil
		}
		for _, interval := range intervals {
			if interval.end < 0 {
				continue
			}
			if interval.start > limit {
				break
			}
			if err := gap(min(limit, interval.start-1)); err != nil {
				return "", err
			}
			next = max(next, interval.end+1)
			if next > limit {
				break
			}
		}
		if err := gap(limit); err != nil {
			return "", err
		}
	}
	if !found {
		return "", fmt.Errorf("no uncovered beacon position")
	}
	return strconv.FormatInt(answer.X*4_000_000+answer.Y, 10), nil
}

func day15p01(r io.Reader) (string, error) {
	return excludedBeacons(r, 2_000_000)
}

func day15p02(r io.Reader) (string, error) {
	return distressBeacon(r, 4_000_000)
}
