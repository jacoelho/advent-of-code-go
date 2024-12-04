package search

import (
	"iter"
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
