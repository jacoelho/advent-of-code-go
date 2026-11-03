package search

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

type testEdge struct {
	to   string
	cost int
}

func graphEdges(graph map[string][]testEdge) Edges[string] {
	return func(from string, yield func(string, int) bool) error {
		for _, edge := range graph[from] {
			if !yield(edge.to, edge.cost) {
				break
			}
		}
		return nil
	}
}

func TestAStar(t *testing.T) {
	for _, test := range []struct {
		name      string
		graph     map[string][]testEdge
		heuristic map[string]int
		cost      int
		path      []string
	}{
		{
			name: "stale entry",
			graph: map[string][]testEdge{
				"s": {{"a", 10}, {"b", 1}}, "b": {{"a", 1}}, "a": {{"g", 100}},
			},
			cost: 102, path: []string{"s", "b", "a", "g"},
		},
		{
			name: "inconsistent admissible heuristic reopens",
			graph: map[string][]testEdge{
				"s": {{"a", 2}, {"b", 1}}, "b": {{"a", 0}, {"g", 100}}, "a": {{"g", 2}},
			},
			heuristic: map[string]int{"s": 3, "a": 0, "b": 2, "g": 0},
			cost:      3, path: []string{"s", "b", "a", "g"},
		},
		{
			name: "zero cost cycle",
			graph: map[string][]testEdge{
				"s": {{"a", 0}}, "a": {{"s", 0}, {"g", 2}},
			},
			cost: 2, path: []string{"s", "a", "g"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			goal := func(node string) bool { return node == "g" }
			var heuristic func(string) int
			if test.heuristic != nil {
				heuristic = func(node string) int { return test.heuristic[node] }
			}
			expanded := make(map[string]int)
			edges := func(node string, yield func(string, int) bool) error {
				expanded[node]++
				return graphEdges(test.graph)(node, yield)
			}
			cost, path, found, err := AStar("s", goal, edges, heuristic, Limits{})
			if err != nil || !found || cost != test.cost || !slices.Equal(path, test.path) {
				t.Fatalf("got (%d, %v, %v, %v), want (%d, %v, true, nil)", cost, path, found, err, test.cost, test.path)
			}
			if test.name == "stale entry" && expanded["a"] != 1 {
				t.Fatalf("expanded a %d times", expanded["a"])
			}
			if cost, found, err := AStarDistance("s", goal, graphEdges(test.graph), heuristic, Limits{}); err != nil || !found || cost != test.cost {
				t.Fatalf("distance = (%d, %v, %v), want %d", cost, found, err, test.cost)
			}
		})
	}
}

func TestWeightedSearchTerminalAndFailure(t *testing.T) {
	noEdges := func(string, func(string, int) bool) error { t.Fatal("expanded goal"); return nil }
	cost, path, found, err := AStar("s", func(string) bool { return true }, noEdges, nil, Limits{})
	if err != nil || !found || cost != 0 || !slices.Equal(path, []string{"s"}) {
		t.Fatalf("start goal = (%d, %v, %v, %v)", cost, path, found, err)
	}
	goal := func(node string) bool { return node == "g" }
	if cost, found, err := AStarDistance("s", goal, graphEdges(nil), nil, Limits{}); err != nil || found || cost != 0 {
		t.Fatalf("unreachable = (%d, %v, %v)", cost, found, err)
	}
	failure := errors.New("expansion failed")
	cost, found, err = AStarDistance("s", goal, func(_ string, yield func(string, int) bool) error {
		yield("g", 1)
		return failure
	}, nil, Limits{})
	if !errors.Is(err, failure) || found || cost != 0 {
		t.Fatalf("expansion failure = (%d, %v, %v)", cost, found, err)
	}
}

func TestWeightedSearchInvalidCosts(t *testing.T) {
	goal := func(node string) bool { return node == "g" }
	for _, test := range []struct {
		name      string
		graph     map[string][]testEdge
		heuristic func(string) int
	}{
		{"negative edge", map[string][]testEdge{"s": {{"g", -1}}}, nil},
		{"path overflow", map[string][]testEdge{"s": {{"a", math.MaxInt}}, "a": {{"g", 1}}}, nil},
		{"priority overflow", map[string][]testEdge{"s": {{"a", math.MaxInt}}}, func(string) int { return 1 }},
		{"negative heuristic", nil, func(string) int { return -1 }},
	} {
		t.Run(test.name, func(t *testing.T) {
			if cost, found, err := AStarDistance("s", goal, graphEdges(test.graph), test.heuristic, Limits{}); err == nil || found || cost != 0 {
				t.Fatalf("invalid search = (%d, %v, %v)", cost, found, err)
			}
		})
	}
	stopped := false
	_, _, err := AStarDistance("s", goal, func(_ string, yield func(string, int) bool) error {
		stopped = !yield("g", -1)
		return nil
	}, nil, Limits{})
	if err == nil || !stopped {
		t.Fatalf("invalid edge did not stop producer: stopped=%v, err=%v", stopped, err)
	}
	if cost, found, err := AStarDistance("s", goal, graphEdges(map[string][]testEdge{"s": {{"g", math.MaxInt}}}), nil, Limits{}); err != nil || !found || cost != math.MaxInt {
		t.Fatalf("maximum valid cost = (%d, %v, %v)", cost, found, err)
	}
}

func TestWeightedSearchLimits(t *testing.T) {
	goal := func(node string) bool { return node == "g" }
	graph := map[string][]testEdge{"s": {{"a", 1}, {"a", 1}, {"b", 2}}, "a": {{"g", 3}}}
	for _, test := range []struct {
		name   string
		limits Limits
		ok     bool
	}{
		{"unlimited", Limits{}, true},
		{"expanded at boundary", Limits{MaxExpanded: 3}, true},
		{"expanded exceeded", Limits{MaxExpanded: 2}, false},
		{"states at boundary", Limits{MaxStates: 4}, true},
		{"states exceeded", Limits{MaxStates: 3}, false},
		{"frontier at boundary", Limits{MaxFrontier: 2}, true},
		{"frontier exceeded", Limits{MaxFrontier: 1}, false},
		{"negative", Limits{MaxStates: -1}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			cost, found, err := AStarDistance("s", goal, graphEdges(graph), nil, test.limits)
			if test.ok {
				if err != nil || !found || cost != 4 {
					t.Fatalf("bounded search = (%d, %v, %v)", cost, found, err)
				}
			} else if err == nil || found || cost != 0 {
				t.Fatalf("expected limit error, got (%d, %v, %v)", cost, found, err)
			}
		})
	}
	// A superseded queued entry still consumes frontier space until it is popped.
	stale := map[string][]testEdge{"s": {{"a", 10}, {"b", 1}}, "b": {{"a", 1}, {"g", 100}}}
	if _, _, err := AStarDistance("s", goal, graphEdges(stale), nil, Limits{MaxFrontier: 2}); err == nil {
		t.Fatal("stale frontier entry escaped the limit")
	}
}

func TestDijkstraOptimalStates(t *testing.T) {
	graph := map[string][]testEdge{
		"s": {{"a", 1}, {"b", 1}, {"x", 9}}, "a": {{"c", 1}}, "b": {{"c", 1}},
		"c": {{"g1", 1}, {"g2", 1}}, "g1": {{"beyond", 1}},
	}
	goal := func(node string) bool { return node == "g1" || node == "g2" }
	cost, states, found, err := DijkstraOptimalStates("s", goal, graphEdges(graph), Limits{})
	slices.Sort(states)
	if err != nil || !found || cost != 3 || !slices.Equal(states, []string{"a", "b", "c", "g1", "g2", "s"}) {
		t.Fatalf("optimal states = (%d, %v, %v, %v)", cost, states, found, err)
	}
	for _, cost := range []int{-1, 0} {
		if _, _, _, err := DijkstraOptimalStates("s", goal, graphEdges(map[string][]testEdge{"s": {{"g1", cost}}}), Limits{}); err == nil {
			t.Fatalf("accepted edge cost %d", cost)
		}
	}
}

func TestWeightedSearchAgainstDistances(t *testing.T) {
	// Floyd-Warshall supplies an independent oracle for costs and optimal states.
	random := rand.New(rand.NewPCG(1, 2))
	const infinity = 1_000_000
	for trial := range 80 {
		n := 2 + random.IntN(5)
		weights, distances := make([][]int, n), make([][]int, n)
		for i := range n {
			weights[i], distances[i] = make([]int, n), make([]int, n)
			for j := range n {
				distances[i][j] = infinity
				if i != n-1 && i != j && random.IntN(3) == 0 {
					weights[i][j] = 1 + random.IntN(5)
					distances[i][j] = weights[i][j]
				}
			}
			distances[i][i] = 0
		}
		for k := range n {
			for i := range n {
				for j := range n {
					distances[i][j] = min(distances[i][j], distances[i][k]+distances[k][j])
				}
			}
		}
		edges := func(from int, yield func(int, int) bool) error {
			for to, cost := range weights[from] {
				if cost > 0 && !yield(to, cost) {
					break
				}
			}
			return nil
		}
		goal := func(node int) bool { return node == n-1 }
		want, reachable := distances[0][n-1], distances[0][n-1] < infinity
		cost, path, found, err := AStar(0, goal, edges, nil, Limits{})
		if err != nil || found != reachable || found && cost != want {
			t.Fatalf("trial %d: path (%d, %v, %v), want (%d, %v)", trial, cost, found, err, want, reachable)
		}
		if found {
			pathCost := 0
			if len(path) == 0 || path[0] != 0 || path[len(path)-1] != n-1 {
				t.Fatalf("trial %d: invalid endpoints %v", trial, path)
			}
			for i := 1; i < len(path); i++ {
				weight := weights[path[i-1]][path[i]]
				if weight == 0 {
					t.Fatalf("trial %d: nonexistent path edge %v", trial, path)
				}
				pathCost += weight
			}
			if pathCost != want {
				t.Fatalf("trial %d: path cost %d, want %d", trial, pathCost, want)
			}
		}
		cost, states, found, err := DijkstraOptimalStates(0, goal, edges, Limits{})
		var wantStates []int
		if reachable {
			for node := range n {
				if distances[0][node]+distances[node][n-1] == want {
					wantStates = append(wantStates, node)
				}
			}
		}
		slices.Sort(states)
		if err != nil || found != reachable || found && cost != want || !slices.Equal(states, wantStates) {
			t.Fatalf("trial %d: states (%d, %v, %v, %v), want (%d, %v, %v)", trial, cost, states, found, err, want, wantStates, reachable)
		}
	}
}
