package search

import (
	"iter"
	"slices"
	"testing"
)

func treeTraversals[T any]() map[string]func(T, func(T) []T) iter.Seq[T] {
	return map[string]func(T, func(T) []T) iter.Seq[T]{
		"DFS": DFSTree[T],
		"BFS": BFSTree[T],
	}
}

func TestTreeTraversal(t *testing.T) {
	type node []int
	for name, traverse := range treeTraversals[node]() {
		t.Run(name, func(t *testing.T) {
			expansions := 0
			children := func(current node) []node {
				expansions++
				switch current[0] {
				case 1:
					return []node{{2}, {3}}
				case 2, 3:
					return []node{{4}}
				default:
					return nil
				}
			}
			sequence := traverse(node{1}, children)
			if expansions != 0 {
				t.Fatal("constructing a traversal expanded its root")
			}
			want := map[string][]int{"DFS": {1, 3, 4, 2, 4}, "BFS": {1, 2, 3, 4, 4}}[name]
			for pass := range 2 {
				var got []int
				for current := range sequence {
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
			for node := range traverse(7, func(int) []int {
				t.Fatal("expanded a node after the consumer stopped")
				return nil
			}) {
				if node != 7 {
					t.Fatalf("got %d; want root 7", node)
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

func TestTreeCopiesChildValues(t *testing.T) {
	for name, traverse := range treeTraversals[string]() {
		t.Run(name, func(t *testing.T) {
			children := []string{"a", "b"}
			var got []string
			for node := range traverse("root", func(node string) []string {
				if node == "root" {
					return children
				}
				clear(children)
				return nil
			}) {
				got = append(got, node)
			}
			want := map[string][]string{"DFS": {"root", "b", "a"}, "BFS": {"root", "a", "b"}}[name]
			if !slices.Equal(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
}
