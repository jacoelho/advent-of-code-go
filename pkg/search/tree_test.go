package search

import (
	"errors"
	"iter"
	"slices"
	"testing"
)

func treeTraversals[T any]() map[string]func(T, func(T) ([]T, error)) iter.Seq2[T, error] {
	return map[string]func(T, func(T) ([]T, error)) iter.Seq2[T, error]{
		"DFS": DFSTree[T],
		"BFS": BFSTree[T],
	}
}

func TestTreeTraversal(t *testing.T) {
	type node []int
	for name, traverse := range treeTraversals[node]() {
		t.Run(name, func(t *testing.T) {
			expansions := 0
			children := func(current node) ([]node, error) {
				expansions++
				switch current[0] {
				case 1:
					return []node{{2}, {3}}, nil
				case 2, 3:
					return []node{{4}}, nil
				default:
					return nil, nil
				}
			}
			sequence := traverse(node{1}, children)
			if expansions != 0 {
				t.Fatal("constructing a traversal expanded its root")
			}
			want := map[string][]int{"DFS": {1, 3, 4, 2, 4}, "BFS": {1, 2, 3, 4, 4}}[name]
			for pass := range 2 {
				var got []int
				for current, err := range sequence {
					if err != nil {
						t.Fatal(err)
					}
					got = append(got, current[0])
				}
				if !slices.Equal(got, want) {
					t.Fatalf("pass %d: got %v, want %v", pass, got, want)
				}
			}
		})
	}
}

func TestTreeEarlyStop(t *testing.T) {
	for name, traverse := range treeTraversals[int]() {
		t.Run(name, func(t *testing.T) {
			visited := false
			for node, err := range traverse(7, func(int) ([]int, error) {
				t.Fatal("expanded a node after the consumer stopped")
				return nil, nil
			}) {
				if node != 7 || err != nil {
					t.Fatalf("got %d, %v; want root 7", node, err)
				}
				visited = true
				break
			}
			if !visited {
				t.Fatal("root was not yielded")
			}
		})
	}
}

func TestTreeExpansionError(t *testing.T) {
	failure := errors.New("expansion failed")
	for name, traverse := range treeTraversals[int]() {
		t.Run(name, func(t *testing.T) {
			first := map[string]int{"DFS": 2, "BFS": 1}[name]
			var got []int
			failures := 0
			for node, err := range traverse(0, func(node int) ([]int, error) {
				if node == 0 {
					return []int{1, 2}, nil
				}
				if node != first {
					t.Errorf("expanded %d after failure", node)
				}
				return []int{3}, failure
			}) {
				if err != nil {
					if !errors.Is(err, failure) {
						t.Fatalf("got error %v, want %v", err, failure)
					}
					failures++
					continue
				}
				if failures != 0 {
					t.Fatalf("yielded node %d after failure", node)
				}
				got = append(got, node)
			}
			if !slices.Equal(got, []int{0, first}) || failures != 1 {
				t.Fatalf("nodes %v, errors %d; want [0 %d] and one error", got, failures, first)
			}
		})
	}
}

func TestTreeCopiesChildValues(t *testing.T) {
	for name, traverse := range treeTraversals[string]() {
		t.Run(name, func(t *testing.T) {
			children := []string{"a", "b"}
			var got []string
			for node, err := range traverse("root", func(node string) ([]string, error) {
				if node == "root" {
					return children, nil
				}
				clear(children)
				return nil, nil
			}) {
				if err != nil {
					t.Fatal(err)
				}
				got = append(got, node)
			}
			want := map[string][]string{"DFS": {"root", "b", "a"}, "BFS": {"root", "a", "b"}}[name]
			if !slices.Equal(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
}
