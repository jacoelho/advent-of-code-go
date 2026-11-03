package aoc2022

import (
	"fmt"
	"io"
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

func valveScores(r io.Reader, minutes int) ([]int64, error) {
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
	if len(useful) > 20 {
		return nil, fmt.Errorf("pressure search supports at most 20 positive-flow valves")
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
		opened uint32
		from   uint8
		time   uint8
	}
	best := make([]int64, 1<<len(useful))
	seen := make(map[state]int64)
	work := 0
	var explore func(state, int64) error
	explore = func(current state, pressure int64) error {
		work++
		if work > 20_000_000 || len(seen) >= 2_000_000 {
			return fmt.Errorf("pressure search exceeded its state or work limit")
		}
		if previous, exists := seen[current]; exists && previous >= pressure {
			return nil
		}
		seen[current] = pressure
		best[current.opened] = max(best[current.opened], pressure)
		for next, valve := range useful {
			bit := uint32(1) << next
			distance := distances[current.from][next]
			remaining := int(current.time) - distance - 1
			if current.opened&bit != 0 || distance < 0 || remaining <= 0 {
				continue
			}
			child := state{current.opened | bit, uint8(next + 1), uint8(remaining)}
			if err := explore(child, pressure+int64(remaining)*valves[valve].rate); err != nil {
				return err
			}
		}
		return nil
	}
	if err := explore(state{time: uint8(minutes)}, 0); err != nil {
		return nil, err
	}
	return best, nil
}

func day16p01(r io.Reader) (string, error) {
	scores, err := valveScores(r, 30)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(slices.Max(scores), 10), nil
}

func day16p02(r io.Reader) (string, error) {
	scores, err := valveScores(r, 26)
	if err != nil {
		return "", err
	}
	subsets := slices.Clone(scores)
	for bit := 1; bit < len(subsets); bit <<= 1 {
		for mask := range subsets {
			if mask&bit != 0 {
				subsets[mask] = max(subsets[mask], subsets[mask^bit])
			}
		}
	}
	var best int64
	for mask, score := range scores {
		best = max(best, score+subsets[len(scores)-1^mask])
	}
	return strconv.FormatInt(best, 10), nil
}
