package aoc2020

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type busDeparture struct {
	offset, interval int64
}

func parseBusSchedule(r io.Reader) (int64, []busDeparture, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, nil, err
	}
	lines := strings.Fields(string(data))
	if len(lines) != 2 {
		return 0, nil, fmt.Errorf("expected departure time and bus schedule")
	}
	earliest, err := strconv.ParseInt(lines[0], 10, 64)
	if err != nil || earliest < 0 {
		return 0, nil, fmt.Errorf("invalid departure time %q", lines[0])
	}
	var buses []busDeparture
	var offset int64
	for item := range strings.SplitSeq(lines[1], ",") {
		if item != "x" {
			interval, err := strconv.ParseInt(item, 10, 64)
			if err != nil || interval <= 0 {
				return 0, nil, fmt.Errorf("invalid bus interval %q", item)
			}
			buses = append(buses, busDeparture{offset, interval})
		}
		offset++
	}
	if len(buses) == 0 {
		return 0, nil, fmt.Errorf("schedule contains no buses")
	}
	return earliest, buses, nil
}

func day13p01(r io.Reader) (string, error) {
	earliest, buses, err := parseBusSchedule(r)
	if err != nil {
		return "", err
	}
	bestWait, bestBus := int64(math.MaxInt64), int64(0)
	for _, bus := range buses {
		if wait := xmath.Modulo(-earliest, bus.interval); wait < bestWait {
			bestWait, bestBus = wait, bus.interval
		}
	}
	product, ok := xmath.CheckedMul(bestWait, bestBus)
	if !ok {
		return "", fmt.Errorf("wait %d times bus %d overflows int64", bestWait, bestBus)
	}
	return strconv.FormatInt(product, 10), nil
}

func day13p02(r io.Reader) (string, error) {
	_, buses, err := parseBusSchedule(r)
	if err != nil {
		return "", err
	}
	timestamp, period := int64(0), int64(1)
	for _, bus := range buses {
		var err error
		timestamp, period, err = xmath.CombineCongruences(timestamp, period, -bus.offset, bus.interval)
		if err != nil {
			return "", fmt.Errorf("bus interval %d at offset %d: %w", bus.interval, bus.offset, err)
		}
	}
	return strconv.FormatInt(timestamp, 10), nil
}
