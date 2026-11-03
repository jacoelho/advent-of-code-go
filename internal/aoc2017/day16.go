package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func dancePower(permutation []int, rounds uint64) []int {
	result := make([]int, len(permutation))
	for start := range permutation {
		length := uint64(1)
		for next := permutation[start]; next != start; next = permutation[next] {
			length++
		}
		result[start] = start
		for range rounds % length {
			result[start] = permutation[result[start]]
		}
	}

	return result
}

func dancePrograms(r io.Reader, size int, rounds uint64) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	input = bytes.TrimSpace(input)
	if len(input) == 0 || size < 1 || size > 26 {
		return "", fmt.Errorf("invalid dance input or program count")
	}

	positions, labels := make([]int, size), make([]int, size)
	for i := range size {
		positions[i], labels[i] = i, i
	}

	for move := range strings.SplitSeq(string(input), ",") {
		if len(move) < 2 {
			return "", fmt.Errorf("invalid dance move %q", move)
		}
		switch move[0] {
		case 's':
			spin, err := strconv.Atoi(move[1:])
			if err != nil || spin < 0 {
				return "", fmt.Errorf("invalid spin %q", move)
			}
			spin %= size
			rotated := make([]int, size)
			for i := range size {
				rotated[(i+spin)%size] = positions[i]
			}
			positions = rotated
		case 'x':
			left, right, ok := strings.Cut(move[1:], "/")
			a, aErr := strconv.Atoi(left)
			b, bErr := strconv.Atoi(right)
			if !ok || aErr != nil || bErr != nil || a < 0 || a >= size || b < 0 || b >= size {
				return "", fmt.Errorf("invalid exchange %q", move)
			}
			positions[a], positions[b] = positions[b], positions[a]
		case 'p':
			if len(move) != 4 || move[2] != '/' || move[1] < 'a' || int(move[1]-'a') >= size || move[3] < 'a' || int(move[3]-'a') >= size {
				return "", fmt.Errorf("invalid partner move %q", move)
			}
			a, b := int(move[1]-'a'), int(move[3]-'a')
			for i, label := range labels {
				switch label {
				case a:
					labels[i] = b
				case b:
					labels[i] = a
				}
			}
		default:
			return "", fmt.Errorf("invalid dance move %q", move)
		}
	}

	// Position moves commute with label renaming, so each permutation can be powered separately.
	positions, labels = dancePower(positions, rounds), dancePower(labels, rounds)
	result := make([]byte, size)
	for i, position := range positions {
		result[i] = byte(labels[position]) + 'a'
	}

	return string(result), nil
}

func day16p01(r io.Reader) (string, error) {
	return dancePrograms(r, 16, 1)
}

func day16p02(r io.Reader) (string, error) {
	return dancePrograms(r, 16, 1_000_000_000)
}
