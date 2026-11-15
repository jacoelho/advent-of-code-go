package aoc2019

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type orbit struct {
	parent string
	child  string
}

func parseOrbit(line []byte) (orbit, error) {
	parent, child, ok := bytes.Cut(line, []byte(")"))
	if !ok || len(parent) == 0 || len(child) == 0 || bytes.ContainsAny(child, ")") {
		return orbit{}, fmt.Errorf("invalid orbit %q", line)
	}
	return orbit{parent: string(parent), child: string(child)}, nil
}

// parseOrbitMap returns each object's parent. Every object orbits at most one
// parent and the orbits form no cycle.
func parseOrbitMap(r io.Reader) (map[string]string, error) {
	s := scanner.NewScanner(r, parseOrbit)
	parents := make(map[string]string)
	for orb := range s.Values() {
		if _, exists := parents[orb.child]; exists {
			return nil, fmt.Errorf("%s orbits more than one object", orb.child)
		}
		parents[orb.child] = orb.parent
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(parents) == 0 {
		return nil, fmt.Errorf("empty orbit map")
	}
	return parents, nil
}

// orbitDepths returns each object's number of direct and indirect orbits.
func orbitDepths(parents map[string]string) (map[string]int, error) {
	depths := make(map[string]int, len(parents)+1)
	var chain []string
	for object := range parents {
		chain = chain[:0]
		current := object
		for {
			if _, known := depths[current]; known {
				break
			}
			parent, ok := parents[current]
			if !ok {
				depths[current] = 0
				break
			}
			chain = append(chain, current)
			if len(chain) > len(parents) {
				return nil, fmt.Errorf("orbit cycle through %s", object)
			}
			current = parent
		}
		for i := len(chain) - 1; i >= 0; i-- {
			depths[chain[i]] = depths[parents[chain[i]]] + 1
		}
	}
	return depths, nil
}

func day06p01(r io.Reader) (string, error) {
	parents, err := parseOrbitMap(r)
	if err != nil {
		return "", err
	}
	depths, err := orbitDepths(parents)
	if err != nil {
		return "", err
	}
	total := 0
	for _, depth := range depths {
		total += depth
	}
	return strconv.Itoa(total), nil
}

func day06p02(r io.Reader) (string, error) {
	parents, err := parseOrbitMap(r)
	if err != nil {
		return "", err
	}
	depths, err := orbitDepths(parents)
	if err != nil {
		return "", err
	}
	you, ok := parents["YOU"]
	if !ok {
		return "", fmt.Errorf("YOU is not in orbit")
	}
	san, ok := parents["SAN"]
	if !ok {
		return "", fmt.Errorf("SAN is not in orbit")
	}

	transfers := 0
	for you != san {
		if depths[you] >= depths[san] {
			next, ok := parents[you]
			if !ok {
				return "", fmt.Errorf("YOU and SAN share no orbit")
			}
			you = next
		} else {
			san = parents[san]
		}
		transfers++
	}
	return strconv.Itoa(transfers), nil
}
