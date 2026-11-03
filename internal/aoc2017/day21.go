package aoc2017

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type fractalPattern struct {
	size int
	bits uint16
}

func parseFractalPattern(text string) (fractalPattern, error) {
	rows := strings.Split(text, "/")
	pattern := fractalPattern{size: len(rows)}
	if pattern.size < 2 || pattern.size > 4 {
		return fractalPattern{}, fmt.Errorf("invalid fractal pattern %q", text)
	}
	for y, row := range rows {
		if len(row) != pattern.size {
			return fractalPattern{}, fmt.Errorf("nonsquare fractal pattern %q", text)
		}
		for x, pixel := range []byte(row) {
			switch pixel {
			case '#':
				pattern.bits |= 1 << (y*pattern.size + x)
			case '.':
			default:
				return fractalPattern{}, fmt.Errorf("invalid fractal pixel %q", pixel)
			}
		}
	}

	return pattern, nil
}

func (p fractalPattern) rotate() fractalPattern {
	result := fractalPattern{size: p.size}
	for y := range p.size {
		for x := range p.size {
			if p.bits&(1<<(y*p.size+x)) != 0 {
				result.bits |= 1 << (x*p.size + p.size - 1 - y)
			}
		}
	}

	return result
}

func (p fractalPattern) reflect() fractalPattern {
	result := fractalPattern{size: p.size}
	for y := range p.size {
		for x := range p.size {
			if p.bits&(1<<(y*p.size+x)) != 0 {
				result.bits |= 1 << (y*p.size + p.size - 1 - x)
			}
		}
	}

	return result
}

func parseFractalRules(r io.Reader) (map[fractalPattern]fractalPattern, error) {
	rules := make(map[fractalPattern]fractalPattern)
	s := bufio.NewScanner(r)
	for s.Scan() {
		left, right, ok := strings.Cut(s.Text(), " => ")
		if !ok {
			return nil, fmt.Errorf("invalid fractal rule %q", s.Text())
		}
		from, err := parseFractalPattern(left)
		if err != nil {
			return nil, err
		}
		to, err := parseFractalPattern(right)
		if err != nil {
			return nil, err
		}
		if from.size > 3 || to.size != from.size+1 {
			return nil, fmt.Errorf("invalid fractal rule sizes: %d => %d", from.size, to.size)
		}
		for range 4 {
			for _, variant := range []fractalPattern{from, from.reflect()} {
				if previous, exists := rules[variant]; exists && previous != to {
					return nil, fmt.Errorf("conflicting fractal rule %q", s.Text())
				}
				rules[variant] = to
			}
			from = from.rotate()
		}
	}

	return rules, s.Err()
}

func enhanceFractal(r io.Reader, iterations int) (int, error) {
	if iterations < 0 || iterations > 18 {
		return 0, fmt.Errorf("fractal iterations must be between 0 and 18")
	}
	rules, err := parseFractalRules(r)
	if err != nil {
		return 0, err
	}

	pixels := []byte{0, 1, 0, 0, 0, 1, 1, 1, 1}
	size := 3
	for range iterations {
		blockSize := 3
		if size%2 == 0 {
			blockSize = 2
		}
		nextSize := size / blockSize * (blockSize + 1)
		next := make([]byte, nextSize*nextSize)
		for by := range size / blockSize {
			for bx := range size / blockSize {
				pattern := fractalPattern{size: blockSize}
				for y := range blockSize {
					for x := range blockSize {
						pattern.bits |= uint16(pixels[(by*blockSize+y)*size+bx*blockSize+x]) << (y*blockSize + x)
					}
				}
				replacement, ok := rules[pattern]
				if !ok {
					return 0, fmt.Errorf("missing %dx%d fractal rule %0*b", blockSize, blockSize, blockSize*blockSize, pattern.bits)
				}
				for y := range replacement.size {
					for x := range replacement.size {
						next[(by*replacement.size+y)*nextSize+bx*replacement.size+x] = byte(replacement.bits >> (y*replacement.size + x) & 1)
					}
				}
			}
		}
		pixels, size = next, nextSize
	}

	count := 0
	for _, pixel := range pixels {
		count += int(pixel)
	}

	return count, nil
}

func day21p01(r io.Reader) (string, error) {
	count, err := enhanceFractal(r, 5)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(count), nil
}

func day21p02(r io.Reader) (string, error) {
	count, err := enhanceFractal(r, 18)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(count), nil
}
