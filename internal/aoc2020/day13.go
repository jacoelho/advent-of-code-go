package aoc2020

import (
	"fmt"
	"io"
	"math/big"
	"strings"
)

type busDeparture struct {
	offset   int64
	interval *big.Int
}

func parseBusSchedule(r io.Reader) (*big.Int, []busDeparture, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Fields(string(data))
	if len(lines) != 2 {
		return nil, nil, fmt.Errorf("expected departure time and bus schedule")
	}
	earliest, ok := new(big.Int).SetString(lines[0], 10)
	if !ok || earliest.Sign() < 0 {
		return nil, nil, fmt.Errorf("invalid departure time %q", lines[0])
	}
	var buses []busDeparture
	var offset int64
	for item := range strings.SplitSeq(lines[1], ",") {
		if item != "x" {
			interval, ok := new(big.Int).SetString(item, 10)
			if !ok || interval.Sign() <= 0 {
				return nil, nil, fmt.Errorf("invalid bus interval %q", item)
			}
			buses = append(buses, busDeparture{offset, interval})
		}
		offset++
	}
	if len(buses) == 0 {
		return nil, nil, fmt.Errorf("schedule contains no buses")
	}
	return earliest, buses, nil
}

func day13p01(r io.Reader) (string, error) {
	earliest, buses, err := parseBusSchedule(r)
	if err != nil {
		return "", err
	}
	var bestWait, bestBus *big.Int
	for _, bus := range buses {
		wait := new(big.Int).Neg(earliest)
		wait.Mod(wait, bus.interval)
		if bestWait == nil || wait.Cmp(bestWait) < 0 {
			bestWait, bestBus = wait, bus.interval
		}
	}
	return new(big.Int).Mul(bestWait, bestBus).String(), nil
}

func day13p02(r io.Reader) (string, error) {
	_, buses, err := parseBusSchedule(r)
	if err != nil {
		return "", err
	}
	timestamp, period := new(big.Int), big.NewInt(1)
	for _, bus := range buses {
		divisor := new(big.Int).GCD(nil, nil, period, bus.interval)
		difference := new(big.Int).Neg(big.NewInt(bus.offset))
		difference.Sub(difference, timestamp)
		if new(big.Int).Mod(difference, divisor).Sign() != 0 {
			return "", fmt.Errorf("incompatible bus interval %s at offset %d", bus.interval, bus.offset)
		}
		modulus := new(big.Int).Quo(bus.interval, divisor)
		if modulus.Cmp(big.NewInt(1)) > 0 {
			step := new(big.Int).Quo(period, divisor)
			inverse := new(big.Int).ModInverse(step, modulus)
			multiple := new(big.Int).Quo(difference, divisor)
			multiple.Mul(multiple, inverse).Mod(multiple, modulus)
			timestamp.Add(timestamp, multiple.Mul(multiple, period))
		}
		period.Mul(period, modulus)
		timestamp.Mod(timestamp, period)
	}
	return timestamp.String(), nil
}
