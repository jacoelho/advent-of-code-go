package aoc2024

import (
	"bufio"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/graph"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
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
		fields := strings.Split(s.Text(), "-")
		a, b := fields[0], fields[1]
		addEdge(a, b)
		addEdge(b, a)
	}
	return result, s.Err()
}

func day23p01(r io.Reader) (string, error) {
	m := aoc.Must(parseNetworkMap(r))

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
	m := aoc.Must(parseNetworkMap(r))

	cliques := graph.MaximalCliques(m)
	longest := xiter.MaxBy(func(a, b collections.Set[string]) bool {
		return b.Len() > a.Len()
	}, cliques)

	result := slices.Collect(longest.Iter())
	slices.Sort(result)
	return strings.Join(result, ","), nil
}
