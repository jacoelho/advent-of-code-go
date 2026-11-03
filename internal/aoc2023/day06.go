package aoc2023

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type race struct {
	time     int64
	distance int64
}

func parseRaces(r io.Reader, combined bool) ([]race, error) {
	const inputLimit = 1 << 20
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return nil, err
	}
	if len(data) > inputLimit {
		return nil, fmt.Errorf("race input exceeds 1 MiB")
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
		if len(fields[i]) == 0 || len(fields[i]) > 10_000 {
			return nil, fmt.Errorf("expected 1 to 10000 race values")
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
			if product > math.MaxInt64/ways {
				overflow = true
			} else {
				product *= ways
			}
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
