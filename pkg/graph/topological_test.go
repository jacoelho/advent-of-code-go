package graph

import (
	"iter"
	"slices"
	"testing"
)

func TestTopologicalOrder(t *testing.T) {
	for _, tt := range []struct {
		name  string
		nodes []string
		edges map[string][]string
		want  []string
	}{
		{name: "empty"},
		{name: "input priority", nodes: []string{"c", "b", "a"}, want: []string{"c", "b", "a"}},
		{name: "newly ready priority", nodes: []string{"c", "b", "a", "d"}, edges: map[string][]string{"a": {"c"}}, want: []string{"b", "a", "c", "d"}},
		{name: "diamond", nodes: []string{"d", "c", "b", "a"}, edges: map[string][]string{"a": {"b", "c"}, "b": {"d"}, "c": {"d"}}, want: []string{"a", "c", "b", "d"}},
		{name: "duplicate edges", nodes: []string{"b", "a"}, edges: map[string][]string{"a": {"b", "b"}}, want: []string{"a", "b"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			expanded := make(map[string]int)
			got, err := TopologicalOrder(tt.nodes, func(node string) iter.Seq[string] {
				expanded[node]++
				return slices.Values(tt.edges[node])
			})
			if err != nil || !slices.Equal(got, tt.want) {
				t.Fatalf("order = %v, %v; want %v", got, err, tt.want)
			}
			for _, node := range tt.nodes {
				if expanded[node] != 1 {
					t.Fatalf("successors(%q) called %d times", node, expanded[node])
				}
			}
		})
	}
}

func TestTopologicalOrderErrors(t *testing.T) {
	for _, tt := range []struct {
		name  string
		nodes []int
		edges map[int][]int
	}{
		{name: "duplicate node", nodes: []int{1, 1}},
		{name: "unknown endpoint", nodes: []int{1}, edges: map[int][]int{1: {2}}},
		{name: "self cycle", nodes: []int{1}, edges: map[int][]int{1: {1}}},
		{name: "disconnected cycle", nodes: []int{0, 1, 2}, edges: map[int][]int{1: {2}, 2: {1}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TopologicalOrder(tt.nodes, func(node int) iter.Seq[int] { return slices.Values(tt.edges[node]) })
			if err == nil || got != nil {
				t.Fatalf("order = %v, %v; want error without partial order", got, err)
			}
		})
	}
}
