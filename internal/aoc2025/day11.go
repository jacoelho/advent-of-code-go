package aoc2025

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

const deviceExit = "out"

type deviceNetwork struct {
	ids     map[string]int
	outputs [][]int
}

func (n *deviceNetwork) id(name string) int {
	id, exists := n.ids[name]
	if !exists {
		id = len(n.outputs)
		n.ids[strings.Clone(name)] = id
		n.outputs = append(n.outputs, nil)
	}
	return id
}

func isDeviceSeparator(c rune) bool { return c == ':' || c == ' ' }

// parseDevicesList reads one record per device. A later record for the same
// device replaces its outputs; devices without a record have no outputs.
func parseDevicesList(r io.Reader) (deviceNetwork, error) {
	network := deviceNetwork{ids: make(map[string]int)}
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	records := 0
	for line := range s.Values() {
		fields := strings.FieldsFunc(line, isDeviceSeparator)
		if len(fields) == 0 {
			return deviceNetwork{}, fmt.Errorf("empty device record")
		}
		device := network.id(fields[0])
		outputs := make([]int, 0, len(fields)-1)
		for _, name := range fields[1:] {
			outputs = append(outputs, network.id(name))
		}
		network.outputs[device] = outputs
		records++
	}
	if err := s.Err(); err != nil {
		return deviceNetwork{}, err
	}
	if records == 0 {
		return deviceNetwork{}, fmt.Errorf("empty device list")
	}
	return network, nil
}

// paths counts the paths from one device to another. Every path ends at the
// exit device, and a reachable cycle is an error.
func (n *deviceNetwork) paths(from, to string) (int, error) {
	start, exists := n.ids[from]
	if !exists {
		return 0, nil
	}
	target, exists := n.ids[to]
	if !exists {
		target = -1
	}
	exit, exists := n.ids[deviceExit]
	if !exists {
		exit = -1
	}
	children := func(device int, yield func(int) bool) {
		for _, next := range n.outputs[device] {
			if !yield(next) {
				return
			}
		}
	}
	terminal := func(device int) (int, bool) {
		switch device {
		case target:
			return 1, true
		case exit:
			return 0, true
		}
		return 0, false
	}
	return graph.CountDAGPaths([]int{start}, children, terminal)
}

// pathsThrough counts paths from start to the exit that visit every device in
// via, in that order.
func (n *deviceNetwork) pathsThrough(start string, via ...string) (int, error) {
	total := 1
	from := start
	for _, to := range append(via, deviceExit) {
		count, err := n.paths(from, to)
		if err != nil {
			return 0, err
		}
		product, ok := xmath.CheckedMul(total, count)
		if !ok {
			return 0, fmt.Errorf("device path count overflows int")
		}
		total = product
		from = to
	}
	return total, nil
}

func day11p01(r io.Reader) (string, error) {
	network, err := parseDevicesList(r)
	if err != nil {
		return "", err
	}
	count, err := network.pathsThrough("you")
	if err != nil {
		return "", fmt.Errorf("device paths: %w", err)
	}
	return strconv.Itoa(count), nil
}

// day11p02 counts paths from svr visiting both dac and fft. The graph is
// acyclic, so each such path visits them in exactly one of the two orders.
func day11p02(r io.Reader) (string, error) {
	network, err := parseDevicesList(r)
	if err != nil {
		return "", err
	}
	dacFirst, err := network.pathsThrough("svr", "dac", "fft")
	if err != nil {
		return "", fmt.Errorf("device paths: %w", err)
	}
	fftFirst, err := network.pathsThrough("svr", "fft", "dac")
	if err != nil {
		return "", fmt.Errorf("device paths: %w", err)
	}
	total, ok := xmath.CheckedAdd(dacFirst, fftFirst)
	if !ok {
		return "", fmt.Errorf("device paths: device path count overflows int")
	}
	return strconv.Itoa(total), nil
}
