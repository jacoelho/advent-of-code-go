package aoc2025

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseDevicesList(r io.Reader) (map[string][]string, error) {
	s := scanner.NewScanner(r, func(line []byte) ([]string, error) {
		fields := strings.FieldsFunc(string(line), func(r rune) bool {
			return r == ':' || r == ' '
		})
		if len(fields) == 0 {
			return nil, fmt.Errorf("empty device record")
		}
		return fields, nil
	})

	result := make(map[string][]string)
	for entry := range s.Values() {
		result[entry[0]] = entry[1:]
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("empty device list")
	}
	return result, nil
}

type deviceState struct {
	device string
	state  int
}

func countPaths(
	devices map[string][]string,
	required []string,
	target string,
	device string,
) (int, error) {
	requiredMap := make(map[string]int, len(required))
	for i, node := range required {
		requiredMap[node] = i
	}
	allRequired := (1 << len(required)) - 1
	children := func(ds deviceState, yield func(deviceState) bool) {
		for _, neighbor := range devices[ds.device] {
			next := deviceState{device: neighbor, state: ds.state}
			if bit, required := requiredMap[neighbor]; required {
				next.state |= 1 << bit
			}
			if !yield(next) {
				return
			}
		}
	}
	terminal := func(ds deviceState) (int, bool) {
		if ds.device != target {
			return 0, false
		}
		if ds.state == allRequired {
			return 1, true
		}
		return 0, true
	}
	return graph.CountDAGPaths([]deviceState{{device: device}}, children, terminal)
}

func day11p01(r io.Reader) (string, error) {
	graph, err := parseDevicesList(r)
	if err != nil {
		return "", err
	}

	count, err := countPaths(graph, nil, "out", "you")
	if err != nil {
		return "", fmt.Errorf("device paths: %w", err)
	}
	return strconv.Itoa(count), nil
}

func day11p02(r io.Reader) (string, error) {
	graph, err := parseDevicesList(r)
	if err != nil {
		return "", err
	}

	required := []string{"dac", "fft"}

	count, err := countPaths(graph, required, "out", "svr")
	if err != nil {
		return "", fmt.Errorf("device paths: %w", err)
	}
	return strconv.Itoa(count), nil
}
