package aoc2016

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
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
	time, period := 0, 1
	for _, disc := range discs {
		var err error
		time, period, err = xmath.CombineCongruences(time, period, -int(disc.number+disc.start), int(disc.positions))
		if err != nil {
			return "", fmt.Errorf("disc %d alignment: %w", disc.number, err)
		}
	}
	return strconv.Itoa(time), nil
}

func day15p01(r io.Reader) (string, error) {
	return capsuleRelease(r, 0)
}

func day15p02(r io.Reader) (string, error) {
	return capsuleRelease(r, 11)
}
