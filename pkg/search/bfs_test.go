package search

import (
	"iter"
	"maps"
	"reflect"
	"slices"
	"testing"
)

func TestBFS(t *testing.T) {
	neighbours := func(node int) iter.Seq[int] {
		graph := map[int][]int{
			1: {2, 3},
			2: {4},
			3: {4, 5},
			4: {6},
			5: {6},
			6: {1},
		}
		return slices.Values(graph[node])
	}

	got := slices.Collect(BFS(1, neighbours))
	want := []int{1, 2, 3, 4, 5, 6}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTraversalReplay(t *testing.T) {
	neighbours := func(node int) iter.Seq[int] {
		return slices.Values(map[int][]int{1: {2, 3}, 2: {4}, 3: {4}, 4: {1}}[node])
	}
	for _, test := range []struct {
		name string
		walk iter.Seq[int]
		want []int
	}{
		{"breadth first", BFS(1, neighbours), []int{1, 2, 3, 4}},
		{"depth first", DFS(1, neighbours), []int{1, 3, 4, 2}},
	} {
		t.Run(test.name, func(t *testing.T) {
			for node := range test.walk {
				if node != 1 {
					break
				}
			}
			for range 2 {
				if got := slices.Collect(test.walk); !slices.Equal(got, test.want) {
					t.Fatalf("replayed traversal = %v, want %v", got, test.want)
				}
			}
		})
	}
}

func TestBFSDistances(t *testing.T) {
	for _, operation := range []string{"all", "maximum", "target", "unreachable", "start"} {
		t.Run(operation, func(t *testing.T) {
			expanded := make(map[int]bool)
			neighbours := func(node int) iter.Seq[int] {
				if expanded[node] {
					t.Fatalf("node %d expanded more than once", node)
				}
				expanded[node] = true
				return slices.Values(map[int][]int{1: {2, 3}, 2: {4}, 3: {4}, 4: {1, 5}}[node])
			}
			switch operation {
			case "all":
				if got := BFSDistances(1, neighbours); !maps.Equal(got, map[int]int{1: 0, 2: 1, 3: 1, 4: 2, 5: 3}) {
					t.Fatalf("distances = %v", got)
				}
			case "maximum":
				if got := BFSMaxDistance(1, neighbours); got != 3 {
					t.Fatalf("maximum distance = %d, want 3", got)
				}
			case "target":
				if got := BFSDistanceTo(1, 4, neighbours); got != 2 {
					t.Fatalf("target distance = %d, want 2", got)
				}
				if expanded[4] || expanded[5] {
					t.Fatal("expanded target or nodes beyond it")
				}
			case "unreachable":
				if got := BFSDistanceTo(1, 6, neighbours); got != -1 {
					t.Fatalf("unreachable distance = %d, want -1", got)
				}
			case "start":
				if got := BFSDistanceTo(1, 1, neighbours); got != 0 || len(expanded) != 0 {
					t.Fatalf("start distance = %d; expanded %v", got, expanded)
				}
			}
		})
	}
}
