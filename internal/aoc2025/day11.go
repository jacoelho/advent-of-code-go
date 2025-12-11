package aoc2025

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/funcs"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseDevicesList(r io.Reader) (map[string][]string, error) {
	s := scanner.NewScanner(r, func(line []byte) ([]string, error) {
		fields := strings.FieldsFunc(string(line), func(r rune) bool {
			return r == ':' || r == ' '
		})
		return fields, nil
	})

	entries := slices.Collect(s.Values())
	result := make(map[string][]string, len(entries))
	for _, entry := range entries {
		result[entry[0]] = entry[1:]
	}

	return result, s.Err()
}

func countPathsToTarget(graph map[string][]string, target string, device string) int {
	var countPaths func(string) int
	countPaths = funcs.Memoize(func(d string) int {
		if d == target {
			return 1
		}

		neighbors, exists := graph[d]
		if !exists {
			return 0
		}

		count := 0
		for _, neighbor := range neighbors {
			count += countPaths(neighbor)
		}
		return count
	})
	return countPaths(device)
}

func day11p01(r io.Reader) (string, error) {
	graph, err := parseDevicesList(r)
	if err != nil {
		return "", err
	}

	count := countPathsToTarget(graph, "out", "you")
	return strconv.Itoa(count), nil
}

type deviceState struct {
	device string
	state  int
}

// allVisitedMask returns a bitmask with all bits set for the given count of required nodes.
func allVisitedMask(requiredCount int) int {
	return (1 << requiredCount) - 1
}

// setVisitedBit sets the bit at the given index in the state.
func setVisitedBit(state int, bitIndex int) int {
	return state | (1 << bitIndex)
}

func countPathsWithRequiredNodes(
	graph map[string][]string,
	required []string,
	target string,
	device string,
) int {
	requiredMap := make(map[string]int, len(required))
	for i, node := range required {
		requiredMap[node] = i
	}

	allVisitedMaskValue := allVisitedMask(len(required))

	var countPaths func(deviceState) int
	countPaths = funcs.Memoize(func(ds deviceState) int {
		if ds.device == target {
			if ds.state == allVisitedMaskValue {
				return 1
			}
			return 0
		}

		neighbors, exists := graph[ds.device]
		if !exists {
			return 0
		}

		count := 0
		for _, neighbor := range neighbors {
			newState := ds.state
			if bitIndex, ok := requiredMap[neighbor]; ok {
				newState = setVisitedBit(newState, bitIndex)
			}
			count += countPaths(deviceState{device: neighbor, state: newState})
		}
		return count
	})
	return countPaths(deviceState{device: device, state: 0})
}

func day11p02(r io.Reader) (string, error) {
	graph, err := parseDevicesList(r)
	if err != nil {
		return "", err
	}

	required := []string{"dac", "fft"}

	count := countPathsWithRequiredNodes(graph, required, "out", "svr")
	return strconv.Itoa(count), nil
}
