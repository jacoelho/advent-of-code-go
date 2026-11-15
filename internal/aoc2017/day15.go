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
		if err != nil || seed < 0 || seed >= generatorModulus {
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

const generatorModulus = 1<<31 - 1

// nextGenerated returns value*factor mod 2^31-1 by folding the high bits, since 2^31 ≡ 1.
func nextGenerated(value, factor uint64) uint64 {
	product := value * factor
	product = product&generatorModulus + product>>31
	if product >= generatorModulus {
		product -= generatorModulus
	}
	return product
}

// judgeGenerators counts low-16-bit matches; each divisor must be a power of two.
func judgeGenerators(r io.Reader, pairs int, divisors [2]uint64) (string, error) {
	seeds, err := parseGeneratorSeeds(r)
	if err != nil {
		return "", err
	}

	a, b := seeds[0], seeds[1]
	maskA, maskB := divisors[0]-1, divisors[1]-1
	matches := 0
	for range pairs {
		a = nextGenerated(a, 16807)
		for a&maskA != 0 {
			a = nextGenerated(a, 16807)
		}
		b = nextGenerated(b, 48271)
		for b&maskB != 0 {
			b = nextGenerated(b, 48271)
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
