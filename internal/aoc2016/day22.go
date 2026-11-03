package aoc2016

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type storageNode struct {
	position   grid.Position2D[int]
	size, used int
}
type storageGrid struct {
	nodes        []storageNode
	adjacent     [][]int
	origin, goal int
}

func readStorage(r io.Reader) (storageGrid, error) {
	var result storageGrid
	indices := make(map[grid.Position2D[int]]int)
	pattern := regexp.MustCompile(`^/dev/grid/node-x(\d+)-y(\d+)\s+(\d+)T\s+(\d+)T\s+(\d+)T\s+\d+%$`)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" || strings.HasPrefix(line, "root@") || strings.HasPrefix(line, "Filesystem") {
			continue
		}
		match := pattern.FindStringSubmatch(line)
		if match == nil {
			return storageGrid{}, fmt.Errorf("invalid storage node %q", line)
		}
		var values [5]int
		for i := range values {
			value, err := convert.ScanNumber[int]([]byte(match[i+1]))
			if err != nil {
				return storageGrid{}, fmt.Errorf("storage node: %w", err)
			}
			values[i] = value
		}
		position := grid.Position2D[int]{X: values[0], Y: values[1]}
		if _, exists := indices[position]; exists {
			return storageGrid{}, fmt.Errorf("duplicate storage node %v", position)
		}
		if values[2] <= 0 || values[3] > values[2] || values[4] != values[2]-values[3] {
			return storageGrid{}, fmt.Errorf("inconsistent storage capacity at %v", position)
		}
		indices[position] = len(result.nodes)
		result.nodes = append(result.nodes, storageNode{position, values[2], values[3]})
	}
	if err := s.Err(); err != nil {
		return storageGrid{}, err
	}
	if len(result.nodes) == 0 {
		return storageGrid{}, fmt.Errorf("no storage nodes")
	}
	result.adjacent = make([][]int, len(result.nodes))
	for i, node := range result.nodes {
		for neighbour := range grid.Neighbours4(node.position) {
			if j, exists := indices[neighbour]; exists {
				result.adjacent[i] = append(result.adjacent[i], j)
			}
		}
	}
	var exists bool
	result.origin, exists = indices[grid.Position2D[int]{}]
	if !exists {
		return storageGrid{}, fmt.Errorf("missing access node at origin")
	}
	result.goal = result.origin
	for i, node := range result.nodes {
		if node.position.Y == 0 && node.position.X > result.nodes[result.goal].position.X {
			result.goal = i
		}
	}
	return result, nil
}

func day22p01(r io.Reader) (string, error) {
	g, err := readStorage(r)
	if err != nil {
		return "", err
	}
	pairs := 0
	for i, source := range g.nodes {
		if source.used == 0 {
			continue
		}
		for j, destination := range g.nodes {
			if i != j && source.used <= destination.size-destination.used {
				pairs++
			}
		}
	}
	return strconv.Itoa(pairs), nil
}

func (g storageGrid) movableStorage() ([]bool, int, bool) {
	empty := -1
	for i, node := range g.nodes {
		if node.used == 0 {
			if empty >= 0 {
				return nil, 0, false
			}
			empty = i
		}
	}
	if empty < 0 {
		return nil, 0, false
	}
	movable := make([]bool, len(g.nodes))
	minSize, minUsed, maxSize, maxUsed := math.MaxInt, math.MaxInt, 0, 0
	for i, node := range g.nodes {
		if node.used > g.nodes[empty].size {
			continue
		}
		movable[i] = true
		minSize, maxSize = min(minSize, node.size), max(maxSize, node.size)
		if node.used > 0 {
			minUsed, maxUsed = min(minUsed, node.used), max(maxUsed, node.used)
		}
	}
	if !movable[g.goal] || !movable[g.origin] || maxUsed > minSize || minUsed <= maxSize-minUsed {
		return nil, 0, false
	}
	// These inequalities rule out merges and freeze every node outside the movable set.
	for i, node := range g.nodes {
		if movable[i] {
			continue
		}
		if node.used <= maxSize || node.size-node.used >= minUsed {
			return nil, 0, false
		}
	}
	return movable, empty, true
}

func (g storageGrid) moveCompactStorage(movable []bool, empty int) (int, error) {
	type state struct{ goal, empty int }
	heuristic := func(s state) int {
		if s.goal == g.origin {
			return 0
		}
		goal := g.nodes[s.goal].position
		return goal.Distance(g.nodes[g.origin].position) + goal.Distance(g.nodes[s.empty].position) - 1
	}
	edges := func(current state, yield func(state, int) bool) error {
		for _, from := range g.adjacent[current.empty] {
			if !movable[from] {
				continue
			}
			candidate := state{current.goal, from}
			if from == candidate.goal {
				candidate.goal = current.empty
			}
			if !yield(candidate, 1) {
				return nil
			}
		}
		return nil
	}
	cost, found, err := search.AStarDistance(state{g.goal, empty}, func(s state) bool {
		return s.goal == g.origin
	}, edges, heuristic, search.Limits{MaxExpanded: 1_000_000})
	if err != nil {
		return 0, fmt.Errorf("compact storage search: %w", err)
	}
	if !found {
		return 0, fmt.Errorf("goal data cannot reach the access node")
	}
	return cost, nil
}

func (g storageGrid) moveExactStorage() (int, error) {
	type key struct {
		goal int
		used string
	}
	type state struct {
		goal, cost int
		used       []int
	}
	encode := func(goal int, used []int) key {
		data := make([]byte, 0, 8*len(used))
		for _, value := range used {
			data = binary.LittleEndian.AppendUint64(data, uint64(value))
		}
		return key{goal, string(data)}
	}
	start := state{goal: g.goal, used: make([]int, len(g.nodes))}
	for i, node := range g.nodes {
		start.used[i] = node.used
	}
	queue := collections.NewDeque[state](32)
	queue.PushBack(start)
	seen := collections.NewSet(encode(start.goal, start.used))
	stateLimit := (128 << 20) / (16*len(g.nodes) + 128)
	for queue.Size() > 0 {
		current, _ := queue.PopFront()
		if current.goal == g.origin {
			return current.cost, nil
		}
		for from, payload := range current.used {
			if payload == 0 {
				continue
			}
			for _, to := range g.adjacent[from] {
				if payload > g.nodes[to].size-current.used[to] {
					continue
				}
				next := state{current.goal, current.cost + 1, slices.Clone(current.used)}
				next.used[from], next.used[to] = 0, next.used[to]+payload
				if from == next.goal {
					next.goal = to
				}
				key := encode(next.goal, next.used)
				if seen.Contains(key) {
					continue
				}
				if seen.Len() >= stateLimit {
					return 0, fmt.Errorf("exact storage search exceeded %d states for %d nodes", stateLimit, len(g.nodes))
				}
				seen.Add(key)
				queue.PushBack(next)
			}
		}
	}
	return 0, fmt.Errorf("goal data cannot reach the access node")
}

func day22p02(r io.Reader) (string, error) {
	g, err := readStorage(r)
	if err != nil {
		return "", err
	}
	if g.nodes[g.goal].used == 0 {
		return "", fmt.Errorf("goal node contains no data")
	}
	var moves int
	if movable, empty, ok := g.movableStorage(); ok {
		moves, err = g.moveCompactStorage(movable, empty)
	} else {
		moves, err = g.moveExactStorage()
	}
	if err != nil {
		return "", err
	}
	return strconv.Itoa(moves), nil
}
