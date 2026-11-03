package aoc2020

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func parseMemoryGame(r io.Reader) ([]int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var starting []int
	for item := range strings.SplitSeq(strings.TrimSpace(string(data)), ",") {
		number, err := strconv.Atoi(item)
		if err != nil || number < 0 {
			return nil, fmt.Errorf("invalid starting number %q", item)
		}
		starting = append(starting, number)
	}
	return starting, nil
}

func memoryGameNumber(starting []int, target int) int {
	if target <= len(starting) {
		return starting[target-1]
	}
	// Generated ages are less than target. Only initial values can fall outside
	// this array, so keep those few values separately instead of sizing by max(seed).
	seen := make([]int32, target)
	large := make(map[int]int32)
	for i, number := range starting[:len(starting)-1] {
		if number < target {
			seen[number] = int32(i + 1)
		} else {
			large[number] = int32(i + 1)
		}
	}
	last := starting[len(starting)-1]
	for turn := len(starting); turn < target; turn++ {
		var previous int32
		if last < target {
			previous, seen[last] = seen[last], int32(turn)
		} else {
			previous, large[last] = large[last], int32(turn)
		}
		last = 0
		if previous != 0 {
			last = turn - int(previous)
		}
	}
	return last
}

func day15p01(r io.Reader) (string, error) {
	starting, err := parseMemoryGame(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(memoryGameNumber(starting, 2020)), nil
}

func day15p02(r io.Reader) (string, error) {
	starting, err := parseMemoryGame(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(memoryGameNumber(starting, 30000000)), nil
}
