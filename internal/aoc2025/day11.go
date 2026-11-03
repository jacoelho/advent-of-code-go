package aoc2025

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type deviceNetwork struct {
	ids     map[string]int
	targets []string
}

func isDeviceSeparator(c rune) bool { return c == ':' || c == ' ' }

func parseDevicesList(r io.Reader) (deviceNetwork, error) {
	network := deviceNetwork{ids: make(map[string]int)}
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	for line := range s.Values() {
		line = strings.TrimLeft(line, ": ")
		if line == "" {
			return deviceNetwork{}, fmt.Errorf("empty device record")
		}
		end := strings.IndexAny(line, ": ")
		if end < 0 {
			end = len(line)
		}
		name := line[:end]
		id, exists := network.ids[name]
		if !exists {
			id = len(network.targets)
			network.targets = append(network.targets, "")
		} else {
			// Replace the key too, so it cannot retain the superseded line.
			delete(network.ids, name)
		}
		network.ids[name] = id
		network.targets[id] = line[end:]
	}
	if err := s.Err(); err != nil {
		return deviceNetwork{}, err
	}
	if len(network.targets) == 0 {
		return deviceNetwork{}, fmt.Errorf("empty device list")
	}
	return network, nil
}

func countPaths(devices deviceNetwork, required []string, target, start string) (int, error) {
	startID, exists := devices.ids[start]
	if !exists {
		if start == target && len(required) == 0 {
			return 1, nil
		}
		return 0, nil
	}
	targetID, exists := devices.ids[target]
	if !exists {
		targetID = len(devices.targets)
	}
	maskBits := len(required)
	if maskBits >= strconv.IntSize || len(devices.targets) > math.MaxInt>>maskBits {
		return 0, fmt.Errorf("device path state does not fit in int")
	}
	allRequired := (1 << maskBits) - 1
	requiredBits := make([]int, len(devices.targets)+1)
	for bit, name := range required {
		if id, exists := devices.ids[name]; exists {
			requiredBits[id] = 1 << bit
		} else if name == target {
			requiredBits[targetID] = 1 << bit
		}
	}
	// Low bits record required visits; higher bits identify the device.
	children := func(state int, yield func(int) bool) {
		mask := state & allRequired
		for name := range strings.FieldsFuncSeq(devices.targets[state>>maskBits], isDeviceSeparator) {
			next, exists := devices.ids[name]
			if !exists {
				if name != target {
					// Undefined nonterminal devices are leaves contributing no paths.
					continue
				}
				next = targetID
			}
			if !yield(next<<maskBits | mask | requiredBits[next]) {
				return
			}
		}
	}
	terminal := func(state int) (int, bool) {
		if state>>maskBits != targetID {
			return 0, false
		}
		if state&allRequired == allRequired {
			return 1, true
		}
		return 0, true
	}
	return graph.CountDAGPaths([]int{startID << maskBits}, children, terminal)
}

func day11p01(r io.Reader) (string, error) {
	network, err := parseDevicesList(r)
	if err != nil {
		return "", err
	}
	count, err := countPaths(network, nil, "out", "you")
	if err != nil {
		return "", fmt.Errorf("device paths: %w", err)
	}
	return strconv.Itoa(count), nil
}

func day11p02(r io.Reader) (string, error) {
	network, err := parseDevicesList(r)
	if err != nil {
		return "", err
	}
	count, err := countPaths(network, []string{"dac", "fft"}, "out", "svr")
	if err != nil {
		return "", fmt.Errorf("device paths: %w", err)
	}
	return strconv.Itoa(count), nil
}
