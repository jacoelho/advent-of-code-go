package aoc2019

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/scanner"
)

type orbit struct {
	parent string
	child  string
}

func parseOrbit(line []byte) (orbit, error) {
	parts := strings.Split(string(line), ")")
	if len(parts) != 2 {
		return orbit{}, fmt.Errorf("invalid orbit format: %s", line)
	}
	return orbit{parent: parts[0], child: parts[1]}, nil
}

func parseOrbitMap(r io.Reader) (map[string]string, error) {
	s := scanner.NewScanner(r, parseOrbit)
	parents := make(map[string]string)
	for orb := range s.Values() {
		parents[orb.child] = orb.parent
	}
	return parents, s.Err()
}

func getAncestors(parents map[string]string, node string) []string {
	var ancestors []string
	current := node
	for {
		parent, exists := parents[current]
		if !exists {
			break
		}
		ancestors = append(ancestors, parent)
		current = parent
	}
	return ancestors
}

func day6p01(r io.Reader) (string, error) {
	parents, err := parseOrbitMap(r)
	if err != nil {
		return "", err
	}

	totalOrbits := 0
	for child := range parents {
		ancestors := getAncestors(parents, child)
		totalOrbits += len(ancestors)
	}

	return strconv.Itoa(totalOrbits), nil
}

func day6p02(r io.Reader) (string, error) {
	parents, err := parseOrbitMap(r)
	if err != nil {
		return "", err
	}

	youAncestors := getAncestors(parents, "YOU")
	sanAncestors := getAncestors(parents, "SAN")

	sanSet := make(map[string]int)
	for i, ancestor := range sanAncestors {
		sanSet[ancestor] = i
	}

	var transfers int
	for youDist, ancestor := range youAncestors {
		if sanDist, found := sanSet[ancestor]; found {
			transfers = youDist + sanDist
			break
		}
	}

	return strconv.Itoa(transfers), nil
}
