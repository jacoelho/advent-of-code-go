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
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
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

	addItem := func(k1, k2, k3 string) {
		el := [3]string{k1, k2, k3}
		slices.Sort(el[:])
		result.Add(el)
	}

	for k1 := range m {
		for k2 := range m[k1].Iter() {
			for k3 := range m[k2].Iter() {
				if k1 != k3 && m[k3].Contains(k1) {
					addItem(k1, k2, k3)
				}
			}
		}
	}

	total := xiter.CountBy(func(v [3]string) bool {
		return xslices.Any(func(s string) bool { return strings.HasPrefix(s, "t") }, v[:])
	}, result.Iter())

	return strconv.Itoa(total), nil
}

func day23p02(r io.Reader) (string, error) {
	m := aoc.Must(parseNetworkMap(r))

	cliques := graph.MaximalCliques(m)
	longest := xslices.MaxBy(func(a, b collections.Set[string]) bool {
		return b.Len() > a.Len()
	}, cliques)

	result := slices.Collect(longest.Iter())
	slices.Sort(result)
	return strings.Join(result, ","), nil
}
