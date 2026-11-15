package aoc2022

import (
	"cmp"
	"fmt"
	"io"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type pressureValve struct {
	name    string
	rate    int64
	tunnels []string
}

func valveScores(r io.Reader, minutes int) (map[uint64]int64, error) {
	pattern := regexp.MustCompile(`^Valve ([A-Z]{2}) has flow rate=(\d+); (?:tunnels lead to valves|tunnel leads to valve) ([A-Z]{2}(?:, [A-Z]{2})*)$`)
	s := scanner.NewScanner(r, func(line []byte) (pressureValve, error) {
		match := pattern.FindSubmatch(line)
		if match == nil {
			return pressureValve{}, fmt.Errorf("invalid valve %q", line)
		}
		rate, err := convert.ScanNumber[int32](match[2])
		if err != nil {
			return pressureValve{}, fmt.Errorf("valve flow rate: %w", err)
		}
		return pressureValve{string(match[1]), int64(rate), strings.Split(string(match[3]), ", ")}, nil
	})
	var valves []pressureValve
	ids := make(map[string]int)
	for valve := range s.Values() {
		if _, exists := ids[valve.name]; exists {
			return nil, fmt.Errorf("duplicate valve %s", valve.name)
		}
		ids[valve.name] = len(valves)
		valves = append(valves, valve)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	start, exists := ids["AA"]
	if !exists {
		return nil, fmt.Errorf("missing start valve AA")
	}
	var useful []int
	edges := make([][]int, len(valves))
	for i, valve := range valves {
		if valve.rate > 0 {
			useful = append(useful, i)
		}
		for _, name := range valve.tunnels {
			id, exists := ids[name]
			if !exists {
				return nil, fmt.Errorf("valve %s references missing valve %s", valve.name, name)
			}
			edges[i] = append(edges[i], id)
		}
	}
	if len(useful) > 64 {
		return nil, fmt.Errorf("%d positive-flow valves exceed the 64-bit opened-valve set", len(useful))
	}

	sources := append([]int{start}, useful...)
	distances := make([][]int, len(sources))
	for i, source := range sources {
		distance := make([]int, len(valves))
		for j := range distance {
			distance[j] = -1
		}
		distance[source] = 0
		queue := collections.NewDeque[int](len(valves))
		queue.PushBack(source)
		for queue.Size() > 0 {
			current, _ := queue.PopFront()
			for _, next := range edges[current] {
				if distance[next] < 0 {
					distance[next] = distance[current] + 1
					queue.PushBack(next)
				}
			}
		}
		for _, destination := range useful {
			distances[i] = append(distances[i], distance[destination])
		}
	}

	type state struct {
		opened uint64
		from   int
		time   int
	}
	best := make(map[uint64]int64)
	seen := make(map[state]int64)
	var explore func(state, int64)
	explore = func(current state, pressure int64) {
		if previous, exists := seen[current]; exists && previous >= pressure {
			return
		}
		seen[current] = pressure
		best[current.opened] = max(best[current.opened], pressure)
		for next, valve := range useful {
			bit := uint64(1) << next
			distance := distances[current.from][next]
			remaining := int(current.time) - distance - 1
			if current.opened&bit != 0 || distance < 0 || remaining <= 0 {
				continue
			}
			child := state{current.opened | bit, next + 1, remaining}
			explore(child, pressure+int64(remaining)*valves[valve].rate)
		}
	}
	explore(state{time: minutes}, 0)
	return best, nil
}

func day16p01(r io.Reader) (string, error) {
	scores, err := valveScores(r, 30)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(slices.Max(slices.Collect(maps.Values(scores))), 10), nil
}

func day16p02(r io.Reader) (string, error) {
	scores, err := valveScores(r, 26)
	if err != nil {
		return "", err
	}
	type route struct {
		opened uint64
		score  int64
	}
	routes := make([]route, 0, len(scores))
	for opened, score := range scores {
		routes = append(routes, route{opened, score})
	}
	slices.SortFunc(routes, func(a, b route) int { return cmp.Compare(b.score, a.score) })
	var best int64
	for i, a := range routes {
		if 2*a.score <= best {
			break
		}
		for _, b := range routes[i:] {
			if a.score+b.score <= best {
				break
			}
			if a.opened&b.opened == 0 {
				best = a.score + b.score
				break
			}
		}
	}
	return strconv.FormatInt(best, 10), nil
}
