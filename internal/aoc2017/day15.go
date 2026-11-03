package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseGeneratorSeeds(r io.Reader) ([2]uint64, error) {
	type generatorSeed struct {
		index int
		value uint64
	}

	s := scanner.NewScanner(r, func(line []byte) (generatorSeed, error) {
		name, value, ok := bytes.Cut(line, []byte(" starts with "))
		if !ok {
			return generatorSeed{}, fmt.Errorf("invalid generator description %q", line)
		}

		index := 0
		switch string(name) {
		case "Generator A":
		case "Generator B":
			index = 1
		default:
			return generatorSeed{}, fmt.Errorf("invalid generator name %q", name)
		}

		seed, err := convert.ScanNumber[int64](value)
		if err != nil || seed < 0 || seed >= 2147483647 {
			return generatorSeed{}, fmt.Errorf("invalid generator seed %q", value)
		}

		return generatorSeed{index: index, value: uint64(seed)}, nil
	})

	var seeds [2]uint64
	var seen [2]bool
	for seed := range s.Values() {
		if seen[seed.index] {
			return seeds, fmt.Errorf("duplicate generator %c", 'A'+seed.index)
		}
		seeds[seed.index] = seed.value
		seen[seed.index] = true
	}
	if err := s.Err(); err != nil {
		return seeds, err
	}
	if !seen[0] || !seen[1] {
		return seeds, fmt.Errorf("expected generators A and B")
	}

	return seeds, nil
}

func judgeGenerators(r io.Reader, pairs int, divisors [2]uint64) (string, error) {
	seeds, err := parseGeneratorSeeds(r)
	if err != nil {
		return "", err
	}

	a, b := seeds[0], seeds[1]
	matches := 0
	for range pairs {
		for {
			a = a * 16807 % 2147483647
			if a%divisors[0] == 0 {
				break
			}
		}
		for {
			b = b * 48271 % 2147483647
			if b%divisors[1] == 0 {
				break
			}
		}
		if uint16(a) == uint16(b) {
			matches++
		}
	}

	return strconv.Itoa(matches), nil
}

func day15p01(r io.Reader) (string, error) {
	return judgeGenerators(r, 40_000_000, [2]uint64{1, 1})
}

func day15p02(r io.Reader) (string, error) {
	return judgeGenerators(r, 5_000_000, [2]uint64{4, 8})
}
