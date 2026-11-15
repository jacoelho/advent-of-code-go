package aoc2023

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type race struct {
	time     int64
	distance int64
}

func parseRaces(r io.Reader, combined bool) ([]race, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		return nil, fmt.Errorf("expected time and distance lines")
	}
	var fields [2][]string
	for i, label := range []string{"Time:", "Distance:"} {
		values, ok := strings.CutPrefix(strings.TrimSpace(lines[i]), label)
		if !ok {
			return nil, fmt.Errorf("missing %s label", label)
		}
		fields[i] = strings.Fields(values)
		if len(fields[i]) == 0 {
			return nil, fmt.Errorf("%s has no race values", label)
		}
		for _, value := range fields[i] {
			for _, digit := range value {
				if digit < '0' || digit > '9' {
					return nil, fmt.Errorf("invalid race value %q", value)
				}
			}
		}
	}
	if len(fields[0]) != len(fields[1]) {
		return nil, fmt.Errorf("time and distance counts differ")
	}
	if combined {
		for i := range fields {
			fields[i] = []string{strings.Join(fields[i], "")}
		}
	}
	races := make([]race, len(fields[0]))
	for i := range races {
		time, err := convert.ScanNumber[int64]([]byte(fields[0][i]))
		if err != nil {
			return nil, err
		}
		distance, err := convert.ScanNumber[int64]([]byte(fields[1][i]))
		if err != nil {
			return nil, err
		}
		races[i] = race{time: time, distance: distance}
	}
	return races, nil
}

func waysToWin(r race) (int64, error) {
	if r.time < 2 {
		return 0, nil
	}
	wins := func(hold int64) bool {
		// The denominator is positive on the rising half; division avoids overflow.
		return hold > r.distance/(r.time-hold)
	}
	peak := r.time / 2
	if !wins(peak) {
		return 0, nil
	}
	first, err := search.BinarySearch(int64(1), peak, func(hold int64) (bool, error) { return wins(hold), nil })
	if err != nil {
		return 0, err
	}
	return r.time - 2*first + 1, nil
}

func day06p01(r io.Reader) (string, error) {
	races, err := parseRaces(r, false)
	if err != nil {
		return "", err
	}
	product := int64(1)
	overflow := false
	for _, race := range races {
		ways, err := waysToWin(race)
		if err != nil {
			return "", err
		}
		if ways == 0 {
			product, overflow = 0, false
		} else if !overflow {
			var ok bool
			product, ok = xmath.CheckedMul(product, ways)
			overflow = !ok
		}
	}
	if overflow {
		return "", fmt.Errorf("race product overflows int64")
	}
	return strconv.FormatInt(product, 10), nil
}

func day06p02(r io.Reader) (string, error) {
	races, err := parseRaces(r, true)
	if err != nil {
		return "", err
	}
	ways, err := waysToWin(races[0])
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(ways, 10), nil
}
