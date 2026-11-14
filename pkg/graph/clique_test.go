package graph

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

func TestMaximalCliquesMatchesBruteForce(t *testing.T) {
	random := rand.New(rand.NewPCG(5, 6))
	for trial := range 100 {
		n := 1 + random.IntN(9)
		g := make(map[int]collections.Set[int], n)
		for v := range n {
			g[v] = collections.NewSet[int]()
		}
		for a := range n {
			for b := a + 1; b < n; b++ {
				if random.IntN(2) == 0 {
					g[a].Add(b)
					g[b].Add(a)
				}
			}
		}
		var got []string
		for clique := range MaximalCliques(g) {
			got = append(got, cliqueKey(clique))
		}
		slices.Sort(got)
		if want := bruteForceMaximalCliques(g, n); !slices.Equal(got, want) {
			t.Fatalf("trial %d: cliques = %v; want %v", trial, got, want)
		}
	}
}

func bruteForceMaximalCliques(g map[int]collections.Set[int], n int) []string {
	isClique := func(mask int) bool {
		for a := range n {
			for b := a + 1; b < n; b++ {
				if mask&(1<<a) != 0 && mask&(1<<b) != 0 && !g[a].Contains(b) {
					return false
				}
			}
		}
		return true
	}
	var result []string
	for mask := 1; mask < 1<<n; mask++ {
		if !isClique(mask) {
			continue
		}
		maximal := true
		for v := range n {
			if mask&(1<<v) == 0 && isClique(mask|1<<v) {
				maximal = false
				break
			}
		}
		if maximal {
			clique := collections.NewSet[int]()
			for v := range n {
				if mask&(1<<v) != 0 {
					clique.Add(v)
				}
			}
			result = append(result, cliqueKey(clique))
		}
	}
	slices.Sort(result)
	return result
}

func cliqueKey(clique collections.Set[int]) string {
	members := slices.Sorted(clique.Iter())
	var b strings.Builder
	for _, m := range members {
		b.WriteByte(byte('0' + m))
	}
	return b.String()
}
