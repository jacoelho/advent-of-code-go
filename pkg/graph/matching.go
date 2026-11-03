package graph

import (
	"errors"
	"fmt"
)

var (
	ErrNoPerfectMatching = errors.New("graph has no perfect matching")
	ErrAmbiguousMatching = errors.New("graph has multiple perfect matchings")
)

// UniquePerfectMatching assigns each left vertex a distinct right vertex.
// Both partitions have len(candidates) vertices; candidates[left] contains the
// permitted right indices. Duplicate candidates are allowed. Inputs are not
// modified. The empty graph has one empty matching.
//
// The result is indexed by left vertex. Missing or multiple perfect matchings
// return ErrNoPerfectMatching or ErrAmbiguousMatching, respectively, with no
// partial result. Invalid right indices return an error.
//
// Augmenting paths establish a matching, then excluding each matched edge checks
// uniqueness. This takes O(V²E) time and O(V) additional storage.
func UniquePerfectMatching(candidates [][]int) ([]int, error) {
	n := len(candidates)
	for left, rights := range candidates {
		for _, right := range rights {
			if right < 0 || right >= n {
				return nil, fmt.Errorf("matching edge %d -> %d is outside %d vertices", left, right, n)
			}
		}
	}

	owners, seen := make([]int, n), make([]int, n)
	match := func(excludedLeft, excludedRight int) bool {
		for right := range owners {
			owners[right] = -1
		}
		clear(seen)
		var assign func(int, int) bool
		assign = func(left, generation int) bool {
			for _, right := range candidates[left] {
				if seen[right] == generation || left == excludedLeft && right == excludedRight {
					continue
				}
				seen[right] = generation
				if owners[right] < 0 || assign(owners[right], generation) {
					owners[right] = left
					return true
				}
			}
			return false
		}
		for left := range candidates {
			if !assign(left, left+1) {
				return false
			}
		}
		return true
	}
	if !match(-1, -1) {
		return nil, ErrNoPerfectMatching
	}
	assignment := make([]int, n)
	for right, left := range owners {
		assignment[left] = right
	}
	// Every alternative matching omits an edge of the first one.
	for left, right := range assignment {
		if match(left, right) {
			return nil, ErrAmbiguousMatching
		}
	}
	return assignment, nil
}
