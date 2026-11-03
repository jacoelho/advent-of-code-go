package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

type caveSystem struct {
	neighbours []collections.Set[int]
	small      []uint64
	start, end int
}

func parseCaves(r io.Reader) (caveSystem, error) {
	var caves caveSystem
	ids := make(map[string]int)
	var names []string
	id := func(name string) int {
		if existing, ok := ids[name]; ok {
			return existing
		}
		index := len(names)
		ids[name] = index
		names = append(names, name)
		caves.neighbours = append(caves.neighbours, collections.NewSet[int]())
		return index
	}
	s := bufio.NewScanner(r)
	for s.Scan() {
		left, right, ok := strings.Cut(s.Text(), "-")
		if !ok || left == "" || right == "" {
			return caves, fmt.Errorf("invalid cave connection %q", s.Text())
		}
		a, b := id(left), id(right)
		caves.neighbours[a].Add(b)
		caves.neighbours[b].Add(a)
	}
	if err := s.Err(); err != nil {
		return caves, err
	}
	var startOK, endOK bool
	caves.start, startOK = ids["start"]
	caves.end, endOK = ids["end"]
	if !startOK || !endOK {
		return caves, fmt.Errorf("cave system requires start and end")
	}
	caves.small = make([]uint64, len(names))
	bits := 0
	for i, name := range names {
		if name != strings.ToLower(name) && name != strings.ToUpper(name) {
			return caves, fmt.Errorf("mixed-case cave name %q", name)
		}
		if i != caves.start && i != caves.end && name == strings.ToLower(name) {
			if bits == 64 {
				return caves, fmt.Errorf("cave system exceeds 64 small caves")
			}
			caves.small[i] = uint64(1) << bits
			bits++
		}
	}
	return caves, nil
}

type caveVisit struct {
	cave    int
	visited uint64
	doubled bool
}

func cavePaths(r io.Reader, smallVisitLimit int) (uint64, error) {
	caves, err := parseCaves(r)
	if err != nil {
		return 0, err
	}
	initial := caveVisit{cave: caves.start, doubled: smallVisitLimit == 1}
	states := []caveVisit{initial}
	ids := map[caveVisit]int{initial: 0}
	parents := [][]int{nil}
	var endings []int
	edges := 0
	for index := 0; index < len(states); index++ {
		current := states[index]
		if current.cave == caves.end {
			endings = append(endings, index)
			continue
		}
		for next := range caves.neighbours[current.cave] {
			if next == caves.start {
				continue
			}
			visit := current
			visit.cave = next
			if bit := caves.small[next]; bit != 0 {
				if visit.visited&bit != 0 {
					if visit.doubled {
						continue
					}
					visit.doubled = true
				} else {
					visit.visited |= bit
				}
			}
			target, known := ids[visit]
			if !known {
				if len(states) == 1_000_000 {
					return 0, fmt.Errorf("cave paths exceed state budget")
				}
				target = len(states)
				ids[visit] = target
				states = append(states, visit)
				parents = append(parents, nil)
			}
			if edges == 8_000_000 {
				return 0, fmt.Errorf("cave paths exceed transition budget")
			}
			parents[target] = append(parents[target], index)
			edges++
		}
	}

	// Only cycles that can still reach end create infinitely many complete routes.
	useful := make([]bool, len(states))
	queue := append([]int(nil), endings...)
	for _, ending := range endings {
		useful[ending] = true
	}
	for head := 0; head < len(queue); head++ {
		for _, parent := range parents[queue[head]] {
			if !useful[parent] {
				useful[parent] = true
				queue = append(queue, parent)
			}
		}
	}
	usefulCount := len(queue)
	pending := make([]int, len(states))
	for target, sources := range parents {
		if useful[target] {
			for _, source := range sources {
				if useful[source] {
					pending[source]++
				}
			}
		}
	}
	paths := make([]uint64, len(states))
	queue = queue[:0]
	for _, ending := range endings {
		paths[ending] = 1
		queue = append(queue, ending)
	}
	for head := 0; head < len(queue); head++ {
		child := queue[head]
		for _, parent := range parents[child] {
			if !useful[parent] {
				continue
			}
			if paths[parent] > math.MaxUint64-paths[child] {
				return 0, fmt.Errorf("cave path count overflows uint64")
			}
			paths[parent] += paths[child]
			pending[parent]--
			if pending[parent] == 0 {
				queue = append(queue, parent)
			}
		}
	}
	if len(queue) != usefulCount {
		return 0, fmt.Errorf("cave system has infinitely many routes to end")
	}
	return paths[0], nil
}

func day12p01(r io.Reader) (string, error) {
	paths, err := cavePaths(r, 1)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(paths, 10), nil
}

func day12p02(r io.Reader) (string, error) {
	paths, err := cavePaths(r, 2)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(paths, 10), nil
}
