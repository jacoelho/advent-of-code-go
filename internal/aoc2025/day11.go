package aoc2025

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/funcs"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type deviceEntry struct {
	device  string
	outputs []string
}

func parseDevicesList(r io.Reader) (map[string]collections.Set[string], error) {
	s := scanner.NewScanner(r, func(line []byte) (deviceEntry, error) {
		fields := strings.FieldsFunc(string(line), func(r rune) bool {
			return r == ':' || r == ' '
		})
		return deviceEntry{
			device:  fields[0],
			outputs: fields[1:],
		}, nil
	})

	entries := slices.Collect(s.Values())
	result := make(map[string]collections.Set[string], len(entries))
	for _, entry := range entries {
		result[entry.device] = collections.NewSet(entry.outputs...)
	}

	return result, s.Err()
}

func countPathsToTarget(graph map[string]collections.Set[string], target string, device string) int {
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
		for neighbor := range neighbors {
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
	graph map[string]collections.Set[string],
	required collections.Set[string],
	target string,
	device string,
) int {
	requiredMap := make(map[string]int)
	for i, node := range xiter.Enumerate(required.Iter()) {
		requiredMap[node] = i
	}

	allVisitedMaskValue := allVisitedMask(required.Len())

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
		for neighbor := range neighbors {
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

	required := collections.NewSet("dac", "fft")

	count := countPathsWithRequiredNodes(graph, required, "out", "svr")
	return strconv.Itoa(count), nil
}
