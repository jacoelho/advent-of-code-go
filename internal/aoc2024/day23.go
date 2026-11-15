package aoc2024

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parseNetworkMap(r io.Reader) (map[string]collections.Set[string], error) {
	result := make(map[string]collections.Set[string])

	addEdge := func(key, value string) {
		if _, exists := result[key]; !exists {
			result[key] = collections.NewSet[string]()
		}
		result[key].Add(value)
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		a, b, ok := strings.Cut(s.Text(), "-")
		if !ok || !day23ValidComputer(a) || !day23ValidComputer(b) || a == b {
			return nil, fmt.Errorf("invalid network connection %q: want two distinct computers as aa-bb", s.Text())
		}
		addEdge(a, b)
		addEdge(b, a)
	}
	return result, s.Err()
}

func day23ValidComputer(name string) bool {
	return len(name) == 2 && name[0] >= 'a' && name[0] <= 'z' && name[1] >= 'a' && name[1] <= 'z'
}

func day23p01(r io.Reader) (string, error) {
	m, err := parseNetworkMap(r)
	if err != nil {
		return "", err
	}

	result := collections.NewSet[[3]string]()

	addTriplet := func(k1, k2, k3 string) {
		el := [3]string{k1, k2, k3}
		slices.Sort(el[:])
		result.Add(el)
	}

	for a, connections := range m {
		if !strings.HasPrefix(a, "t") {
			continue
		}
		for b := range connections {
			for c := range connections.Intersect(m[b]) {
				addTriplet(a, b, c)
			}
		}
	}

	return strconv.Itoa(result.Len()), nil
}

func day23p02(r io.Reader) (string, error) {
	m, err := parseNetworkMap(r)
	if err != nil {
		return "", err
	}

	cliques := graph.MaximalCliques(m)
	longest, ok := xiter.MaxBy(func(a, b collections.Set[string]) bool {
		return b.Len() > a.Len()
	}, cliques)

	if !ok {
		return "", fmt.Errorf("network has no clique")
	}

	result := slices.Collect(longest.Iter())
	slices.Sort(result)
	return strings.Join(result, ","), nil
}
