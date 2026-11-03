package aoc2016

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
)

type capsuleDisc struct{ number, positions, start int64 }

func capsuleRelease(r io.Reader, extraPositions int64) (string, error) {
	pattern := regexp.MustCompile(`^Disc #(\d+) has (\d+) positions; at time=0, it is at position (\d+)\.$`)
	var discs []capsuleDisc
	seen := make(map[int64]bool)
	var last int64
	s := bufio.NewScanner(r)
	for s.Scan() {
		match := pattern.FindStringSubmatch(s.Text())
		if match == nil {
			return "", fmt.Errorf("invalid disc description %q", s.Text())
		}
		var values [3]int64
		for i := range values {
			value, err := strconv.ParseInt(match[i+1], 10, 32)
			if err != nil {
				return "", fmt.Errorf("invalid disc number: %w", err)
			}
			values[i] = value
		}
		disc := capsuleDisc{values[0], values[1], values[2]}
		if disc.number <= 0 || disc.positions <= 0 || disc.start >= disc.positions || seen[disc.number] {
			return "", fmt.Errorf("invalid or duplicate disc %d", disc.number)
		}
		seen[disc.number] = true
		last = max(last, disc.number)
		discs = append(discs, disc)
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(discs) == 0 {
		return "", fmt.Errorf("empty disc list")
	}
	if extraPositions > 0 {
		discs = append(discs, capsuleDisc{last + 1, extraPositions, 0})
	}
	time, period := new(big.Int), big.NewInt(1)
	for _, disc := range discs {
		size := big.NewInt(disc.positions)
		divisor := new(big.Int).GCD(nil, nil, period, size)
		difference := new(big.Int).Sub(big.NewInt(-disc.number-disc.start), time)
		if new(big.Int).Mod(difference, divisor).Sign() != 0 {
			return "", fmt.Errorf("disc %d has incompatible alignment", disc.number)
		}
		modulus := new(big.Int).Quo(size, divisor)
		if modulus.Cmp(big.NewInt(1)) > 0 {
			step := new(big.Int).Quo(period, divisor)
			inverse := new(big.Int).ModInverse(step, modulus)
			multiple := new(big.Int).Quo(difference, divisor)
			multiple.Mul(multiple, inverse).Mod(multiple, modulus)
			time.Add(time, multiple.Mul(multiple, period))
		}
		period.Mul(period, modulus)
		time.Mod(time, period)
	}
	return time.String(), nil
}

func day15p01(r io.Reader) (string, error) {
	return capsuleRelease(r, 0)
}

func day15p02(r io.Reader) (string, error) {
	return capsuleRelease(r, 11)
}
