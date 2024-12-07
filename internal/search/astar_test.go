package search

import (
	"reflect"
	"testing"
)

func TestAStar(t *testing.T) {
	start := "A"
	goal := "F"

	neighbours := func(state string) []string {
		graph := map[string][]string{
			"A": {"B", "C"},
			"B": {"A", "D", "E"},
			"C": {"A", "F"},
			"D": {"B"},
			"E": {"B"},
			"F": {"C"},
		}
		return graph[state]
	}

	heuristic := func(state string) int {
		if state == goal {
			return 0
		}
		return 1 // simple heuristic
	}

	stepCost := func(from, to string) int {
		return 1 // constant step cost
	}

	cost, path, found := AStar(start, neighbours, heuristic, stepCost)
	if !found {
		t.Error("expected to find a path")
		return
	}

	if !reflect.DeepEqual(path, []string{"A", "C", "F"}) {
		t.Errorf("expected to find a path, got: %v", path)
		return
	}

	if cost != 2 {
		t.Errorf("expected cost to be 2, got %d", cost)
	}
}
