package graph

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestUniquePerfectMatching(t *testing.T) {
	for _, tt := range []struct {
		name       string
		candidates [][]int
		want       []int
		err        error
	}{
		{name: "empty", want: []int{}},
		{name: "singleton", candidates: [][]int{{0}}, want: []int{0}},
		{name: "reassignment", candidates: [][]int{{0, 1}, {0}}, want: []int{1, 0}},
		{name: "duplicates", candidates: [][]int{{0, 0}}, want: []int{0}},
		{name: "ambiguous", candidates: [][]int{{0, 1}, {0, 1}}, err: ErrAmbiguousMatching},
		{name: "Hall deficiency", candidates: [][]int{{0}, {0}, {1, 2}}, err: ErrNoPerfectMatching},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UniquePerfectMatching(tt.candidates)
			if !errors.Is(err, tt.err) || !slices.Equal(got, tt.want) {
				t.Fatalf("matching = %v, %v; want %v, %v", got, err, tt.want, tt.err)
			}
		})
	}
	for _, candidates := range [][][]int{{{-1}}, {{1}}} {
		if got, err := UniquePerfectMatching(candidates); err == nil || got != nil {
			t.Fatalf("invalid candidates %v: got %v, %v", candidates, got, err)
		}
	}
}

func TestUniquePerfectMatchingAgainstPermutations(t *testing.T) {
	random := rand.New(rand.NewPCG(16, 2020))
	for n := range 5 {
		for range 64 {
			candidates := make([][]int, n)
			for left := range candidates {
				for right := range n {
					if random.IntN(2) == 0 {
						candidates[left] = append(candidates[left], right)
					}
				}
			}
			var solutions [][]int
			var enumerate func([]int)
			enumerate = func(prefix []int) {
				if len(prefix) == n {
					solutions = append(solutions, slices.Clone(prefix))
					return
				}
				for right := range n {
					if !slices.Contains(prefix, right) && slices.Contains(candidates[len(prefix)], right) {
						enumerate(append(prefix, right))
					}
				}
			}
			enumerate(nil)
			got, err := UniquePerfectMatching(candidates)
			switch len(solutions) {
			case 0:
				if !errors.Is(err, ErrNoPerfectMatching) {
					t.Fatalf("%v: got %v, %v; want no matching", candidates, got, err)
				}
			case 1:
				if err != nil || !slices.Equal(got, solutions[0]) {
					t.Fatalf("%v: got %v, %v; want %v", candidates, got, err, solutions[0])
				}
			default:
				if !errors.Is(err, ErrAmbiguousMatching) {
					t.Fatalf("%v: got %v, %v; want ambiguous matching", candidates, got, err)
				}
			}
		}
	}
}
