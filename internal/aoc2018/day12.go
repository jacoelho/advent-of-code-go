package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func plantPots(r io.Reader, generations int64) (string, error) {
	s := bufio.NewScanner(r)
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("missing initial plant state")
	}
	state, ok := strings.CutPrefix(s.Text(), "initial state: ")
	if !ok || len(state) == 0 || strings.Trim(state, ".#") != "" {
		return "", fmt.Errorf("invalid initial plant state")
	}

	var rules, seenRules [32]bool
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		if len(line) != 10 || line[5:9] != " => " || strings.Trim(line[:5], ".#") != "" || line[9] != '.' && line[9] != '#' {
			return "", fmt.Errorf("invalid plant rule %q", line)
		}
		pattern := 0
		for _, cell := range line[:5] {
			pattern <<= 1
			if cell == '#' {
				pattern++
			}
		}
		if seenRules[pattern] {
			return "", fmt.Errorf("duplicate plant rule %q", line[:5])
		}
		seenRules[pattern], rules[pattern] = true, line[9] == '#'
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if rules[0] {
		return "", fmt.Errorf("plant rule creates infinitely many plants")
	}

	type occurrence struct {
		generation, origin int64
	}
	seen := make(map[string]occurrence)
	origin, generation := int64(0), int64(0)
	updates, storedBytes := 0, 0
	for generation < generations {
		first := strings.IndexByte(state, '#')
		if first < 0 {
			return "0", nil
		}
		origin += int64(first)
		state = state[first : strings.LastIndexByte(state, '#')+1]

		if previous, exists := seen[state]; exists {
			period := generation - previous.generation
			cycles := (generations - generation) / period
			if cycles > 0 {
				origin += cycles * (origin - previous.origin)
				generation += cycles * period
				continue
			}
		} else {
			storedBytes += len(state)
			if storedBytes > 64_000_000 {
				return "", fmt.Errorf("plant cycle history exceeds 64000000 bytes")
			}
			seen[state] = occurrence{generation, origin}
		}
		if updates == 100_000 || len(state) > 1_000_000 {
			return "", fmt.Errorf("plant simulation exceeded its generation or width limit")
		}
		updates++

		padded := "...." + state + "...."
		next := make([]byte, len(state)+4)
		for i := range next {
			pattern := 0
			for _, cell := range padded[i : i+5] {
				pattern <<= 1
				if cell == '#' {
					pattern++
				}
			}
			next[i] = '.'
			if rules[pattern] {
				next[i] = '#'
			}
		}
		state, origin = string(next), origin-2
		generation++
	}

	var sum int64
	for i, cell := range state {
		if cell == '#' {
			sum += origin + int64(i)
		}
	}

	return strconv.FormatInt(sum, 10), nil
}

func day12p01(r io.Reader) (string, error) {
	return plantPots(r, 20)
}

func day12p02(r io.Reader) (string, error) {
	return plantPots(r, 50_000_000_000)
}
