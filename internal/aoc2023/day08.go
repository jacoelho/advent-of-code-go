package aoc2023

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type network map[string][2]string

type networkEntry struct {
	node string
	dest [2]string
}

func parseNetworkEntry(line []byte) (networkEntry, error) {
	parts := strings.Split(string(line), " = ")
	if len(parts) != 2 {
		return networkEntry{}, fmt.Errorf("invalid line: %s", line)
	}

	node := parts[0]
	destinations := strings.Trim(parts[1], "()")
	destParts := strings.Split(destinations, ", ")
	if len(destParts) != 2 {
		return networkEntry{}, fmt.Errorf("invalid destinations: %s", destinations)
	}

	return networkEntry{
		node: node,
		dest: [2]string{destParts[0], destParts[1]},
	}, nil
}

func parseNetwork(r io.Reader) (string, network, error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte("\n\n")), func(b []byte) (string, error) {
		return string(b), nil
	})

	sections := slices.Collect(s.Values())
	if s.Err() != nil {
		return "", nil, s.Err()
	}

	if len(sections) < 2 {
		return "", nil, fmt.Errorf("invalid input format")
	}

	instructions := strings.TrimSpace(sections[0])

	entryScanner := scanner.NewScanner(strings.NewReader(sections[1]), parseNetworkEntry)
	entries := slices.Collect(entryScanner.Values())
	if entryScanner.Err() != nil {
		return "", nil, entryScanner.Err()
	}

	net := make(network)
	for _, entry := range entries {
		net[entry.node] = entry.dest
	}

	return instructions, net, nil
}

func navigate(instructions string, net network, start string, isEnd func(string) bool) int {
	current := start
	steps := 0

	for !isEnd(current) {
		instruction := instructions[steps%len(instructions)]

		if instruction == 'L' {
			current = net[current][0]
		} else {
			current = net[current][1]
		}

		steps++
	}

	return steps
}

func day08p01(r io.Reader) (string, error) {
	instructions, net, err := parseNetwork(r)
	if err != nil {
		return "", err
	}

	steps := navigate(instructions, net, "AAA", func(node string) bool {
		return node == "ZZZ"
	})

	return strconv.Itoa(steps), nil
}

func day08p02(r io.Reader) (string, error) {
	instructions, net, err := parseNetwork(r)
	if err != nil {
		return "", err
	}

	nodes := slices.Collect(maps.Keys(net))
	slices.Sort(nodes)

	starts := xslices.Filter(func(node string) bool {
		return strings.HasSuffix(node, "A")
	}, nodes)

	cycles := xslices.Map(func(start string) int {
		return navigate(instructions, net, start, func(node string) bool {
			return strings.HasSuffix(node, "Z")
		})
	}, starts)

	result := xmath.LCM(cycles[0], cycles[1:]...)

	return strconv.Itoa(result), nil
}
