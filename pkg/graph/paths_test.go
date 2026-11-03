package graph

import (
	"math"
	"testing"
)

func TestCountDAGPaths(t *testing.T) {
	for _, test := range []struct {
		name      string
		roots     []string
		edges     map[string][]string
		terminals map[string]int
		want      int
		wantError bool
	}{
		{
			name: "empty",
		},
		{
			name:      "diamond",
			roots:     []string{"s"},
			edges:     map[string][]string{"s": {"a", "b"}, "a": {"t"}, "b": {"t"}},
			terminals: map[string]int{"t": 1},
			want:      2,
		},
		{
			name:      "terminal stops",
			roots:     []string{"t"},
			edges:     map[string][]string{"t": {"t"}},
			terminals: map[string]int{"t": 3},
			want:      3,
		},
		{
			name:      "zero terminal stops",
			roots:     []string{"t"},
			edges:     map[string][]string{"t": {"t"}},
			terminals: map[string]int{"t": 0},
		},
		{
			name:      "unreachable",
			roots:     []string{"s"},
			edges:     map[string][]string{"s": {"a"}},
			terminals: map[string]int{"t": 1},
		},
		{
			name:      "parallel edges and roots",
			roots:     []string{"s", "s"},
			edges:     map[string][]string{"s": {"t", "t"}},
			terminals: map[string]int{"t": 1},
			want:      4,
		},
		{
			name:      "reachable cycle",
			roots:     []string{"s"},
			edges:     map[string][]string{"s": {"a"}, "a": {"s"}},
			wantError: true,
		},
		{
			name:  "unreachable cycle",
			roots: []string{"s"},
			edges: map[string][]string{"a": {"a"}},
		},
		{
			name:      "edge sum overflow",
			roots:     []string{"s"},
			edges:     map[string][]string{"s": {"a", "b"}},
			terminals: map[string]int{"a": math.MaxInt, "b": 1},
			wantError: true,
		},
		{
			name:      "root sum overflow",
			roots:     []string{"t", "t"},
			terminals: map[string]int{"t": math.MaxInt},
			wantError: true,
		},
		{
			name:      "negative terminal",
			roots:     []string{"t"},
			terminals: map[string]int{"t": -1},
			wantError: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			expanded := make(map[string]bool)
			children := func(node string, yield func(string) bool) {
				if expanded[node] {
					t.Fatalf("expanded %q twice", node)
				}
				expanded[node] = true
				for _, child := range test.edges[node] {
					if !yield(child) {
						return
					}
				}
			}
			terminal := func(node string) (int, bool) {
				count, found := test.terminals[node]
				return count, found
			}
			got, err := CountDAGPaths(test.roots, children, terminal)
			if (err != nil) != test.wantError || got != test.want {
				t.Fatalf("count = %d, %v; want %d, error %v", got, err, test.want, test.wantError)
			}
			for node := range test.terminals {
				if expanded[node] {
					t.Fatalf("expanded terminal %q", node)
				}
			}
		})
	}
}

func TestCountDAGPathsDeepChain(t *testing.T) {
	const length = 100_000
	children := func(node int, yield func(int) bool) { yield(node + 1) }
	terminal := func(node int) (int, bool) { return 1, node == length }
	got, err := CountDAGPaths([]int{0}, children, terminal)
	if err != nil || got != 1 {
		t.Fatalf("deep chain count = %d, %v; want 1", got, err)
	}
}
