package search

import (
	"fmt"
	"math"
	"slices"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

// Edges emits a state's outgoing edges and their costs. It must stop when yield
// returns false and must not retain yield or call it concurrently. An error
// aborts the search, including any partial expansion.
// State identity, terminal states, and finite graph generation belong to callers.
type Edges[T comparable] func(T, func(T, int) bool) error

// Limits bounds search storage and expansion. Zero means unlimited; negative
// limits are invalid. MaxStates bounds distinct discovered states, including
// start. MaxFrontier bounds queued entries, including stale entries.
// MaxExpanded bounds non-goal, non-stale entries whose edges are requested,
// including reopenings.
type Limits struct {
	MaxExpanded int
	MaxStates   int
	MaxFrontier int
}

// AStarDistance finds the minimum cost to a terminal goal without storing paths.
// Edge costs must be nonnegative. A nil heuristic means zero (Dijkstra); otherwise
// it must be nonnegative, admissible, and zero at goals. Improved states reopen,
// so consistency is not required. Validation applies to encountered edges and
// estimates. A failed search returns (0, false, nil); errors return no result.
func AStarDistance[T comparable](start T, goal func(T) bool, edges Edges[T], heuristic func(T) int, limits Limits) (int, bool, error) {
	result, err := weightedSearch(start, goal, edges, heuristic, limits, distanceOnly)
	return result.cost, result.found, err
}

// AStar has the same contract as AStarDistance and also returns one optimal path,
// including start and goal. Ties between optimal paths are unspecified.
func AStar[T comparable](start T, goal func(T) bool, edges Edges[T], heuristic func(T) int, limits Limits) (int, []T, bool, error) {
	result, err := weightedSearch(start, goal, edges, heuristic, limits, onePath)
	if err != nil || !result.found {
		return 0, nil, false, err
	}
	path := []T{result.goals[0]}
	for path[len(path)-1] != start {
		path = append(path, result.previous[path[len(path)-1]])
	}
	slices.Reverse(path)
	return result.cost, path, true, nil
}

// DijkstraOptimalStates returns the union of states on all minimum-cost routes
// to any terminal goal, in unspecified order. Costs must be strictly positive:
// predecessors then form a DAG, avoiding ambiguous zero-cost detours. It stores
// predecessor edges, never complete paths. Failure and limits follow AStarDistance.
func DijkstraOptimalStates[T comparable](start T, goal func(T) bool, edges Edges[T], limits Limits) (int, []T, bool, error) {
	result, err := weightedSearch(start, goal, edges, nil, limits, optimalStates)
	if err != nil || !result.found {
		return 0, nil, false, err
	}
	seen := collections.NewSet[T]()
	pending := result.goals
	var states []T
	for len(pending) > 0 {
		current := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		if seen.Contains(current) {
			continue
		}
		seen.Add(current)
		states = append(states, current)
		pending = append(pending, result.parents[current]...)
	}
	return result.cost, states, true, nil
}

type retention uint8

const (
	distanceOnly retention = iota
	onePath
	optimalStates
)

type weightedEntry[T comparable] struct {
	state       T
	cost, score int
}

type weightedResult[T comparable] struct {
	cost     int
	found    bool
	goals    []T
	previous map[T]T
	parents  map[T][]T
}

func weightedSearch[T comparable](start T, goal func(T) bool, edges Edges[T], heuristic func(T) int, limits Limits, keep retention) (weightedResult[T], error) {
	var result weightedResult[T]
	fail := func(err error) (weightedResult[T], error) { return weightedResult[T]{}, err }
	if limits.MaxExpanded < 0 || limits.MaxStates < 0 || limits.MaxFrontier < 0 {
		return fail(fmt.Errorf("search limits must be nonnegative"))
	}
	estimate := func(state T) (int, error) {
		if heuristic == nil {
			return 0, nil
		}
		value := heuristic(state)
		if value < 0 {
			return 0, fmt.Errorf("search heuristic must be nonnegative")
		}
		return value, nil
	}
	score, err := estimate(start)
	if err != nil {
		return fail(err)
	}
	frontier := collections.NewHeap(func(a, b weightedEntry[T]) bool { return a.score < b.score })
	frontier.Push(weightedEntry[T]{state: start, score: score})
	distances := map[T]int{start: 0}
	switch keep {
	case onePath:
		result.previous = make(map[T]T)
	case optimalStates:
		result.parents = make(map[T][]T)
	}
	expanded := 0
	var current weightedEntry[T]
	var edgeErr error
	yieldEdge := func(next T, cost int) bool {
		if edgeErr != nil {
			return false
		}
		if cost < 0 || keep == optimalStates && cost == 0 {
			edgeErr = fmt.Errorf("invalid search edge cost %d", cost)
			return false
		}
		if current.cost > math.MaxInt-cost {
			edgeErr = fmt.Errorf("search path cost overflows int")
			return false
		}
		distance := current.cost + cost
		previous, exists := distances[next]
		if exists && distance >= previous {
			if keep == optimalStates && distance == previous {
				result.parents[next] = append(result.parents[next], current.state)
			}
			return true
		}
		heuristic, err := estimate(next)
		if err != nil {
			edgeErr = err
			return false
		}
		if distance > math.MaxInt-heuristic {
			edgeErr = fmt.Errorf("search priority overflows int")
			return false
		}
		if !exists && limits.MaxStates > 0 && len(distances) == limits.MaxStates {
			edgeErr = fmt.Errorf("search exceeds %d discovered states", limits.MaxStates)
			return false
		}
		if limits.MaxFrontier > 0 && frontier.Len() == limits.MaxFrontier {
			edgeErr = fmt.Errorf("search exceeds %d queued entries", limits.MaxFrontier)
			return false
		}
		distances[next] = distance
		frontier.Push(weightedEntry[T]{state: next, cost: distance, score: distance + heuristic})
		switch keep {
		case onePath:
			result.previous[next] = current.state
		case optimalStates:
			result.parents[next] = []T{current.state}
		}
		return true
	}
	for entry := range frontier.PopSeq() {
		current = entry
		if current.cost != distances[current.state] {
			continue
		}
		if result.found && current.cost > result.cost {
			break // Only the Dijkstra all-optimal-states search continues past a goal.
		}
		if goal(current.state) {
			result.cost, result.found = current.cost, true
			if keep != distanceOnly {
				result.goals = append(result.goals, current.state)
			}
			if keep != optimalStates {
				return result, nil
			}
			continue
		}
		if limits.MaxExpanded > 0 && expanded == limits.MaxExpanded {
			return fail(fmt.Errorf("search exceeds %d expanded states", limits.MaxExpanded))
		}
		expanded++
		edgeErr = nil
		err := edges(current.state, yieldEdge)
		if err != nil {
			return fail(err)
		}
		if edgeErr != nil {
			return fail(edgeErr)
		}
	}
	return result, nil
}
