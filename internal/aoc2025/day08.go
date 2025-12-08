package aoc2025

import (
	"cmp"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseJunctionBoxes(r io.Reader) ([]grid.Position3D[int], error) {
	s := scanner.NewScanner(r, func(line []byte) (grid.Position3D[int], error) {
		nums := convert.ExtractDigits[int](string(line))
		if len(nums) != 3 {
			return grid.Position3D[int]{}, fmt.Errorf("expected 3 integers, got %d", len(nums))
		}
		return grid.NewPosition3D(nums[0], nums[1], nums[2]), nil
	})

	return slices.Collect(s.Values()), s.Err()
}

// unionFind is a union-find (disjoint set) data structure.
type unionFind struct {
	parent, size []int
	distinct     int
}

func newUnionFind(n int) *unionFind {
	parent := make([]int, n)
	size := make([]int, n)
	for i := range parent {
		parent[i] = i
		size[i] = 1
	}
	return &unionFind{parent: parent, size: size, distinct: n}
}

// find returns the root of the set containing x, applying path compression.
func (u *unionFind) find(x int) int {
	root := x
	for u.parent[root] != root {
		root = u.parent[root]
	}
	for u.parent[x] != x {
		next := u.parent[x]
		u.parent[x] = root
		x = next
	}
	return root
}

// union merges the sets containing x and y, returning true if merged.
func (u *unionFind) union(x, y int) bool {
	rootX, rootY := u.find(x), u.find(y)
	if rootX == rootY {
		return false
	}
	if u.size[rootX] < u.size[rootY] {
		rootX, rootY = rootY, rootX
	}
	u.parent[rootY] = rootX
	u.size[rootX] += u.size[rootY]
	u.distinct--
	return true
}

func (u *unionFind) sizes() map[int]int {
	sizes := make(map[int]int, u.distinct)
	for i := range u.parent {
		root := u.find(i)
		if i == root {
			sizes[root] = u.size[root]
		}
	}
	return sizes
}

func (u *unionFind) count() int {
	return u.distinct
}

type boxPair struct {
	indexA, indexB int
	dist           int
}

func buildBoxPairs(boxes []grid.Position3D[int]) []boxPair {
	n := len(boxes)
	boxPairs := make([]boxPair, 0, n*(n-1)/2)
	for i := range boxes {
		for j := i + 1; j < n; j++ {
			pair := boxPair{
				indexA: i,
				indexB: j,
				dist:   boxes[i].EuclideanDistanceSquared(boxes[j]),
			}
			boxPairs = append(boxPairs, pair)
		}
	}
	slices.SortFunc(boxPairs, func(a, b boxPair) int {
		return cmp.Compare(a.dist, b.dist)
	})
	return boxPairs
}

func day08p01(connections int) func(r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		boxes, err := parseJunctionBoxes(r)
		if err != nil {
			return "", err
		}

		boxPairs := buildBoxPairs(boxes)

		uf := newUnionFind(len(boxes))
		for _, pair := range boxPairs[:connections] {
			uf.union(pair.indexA, pair.indexB)
		}

		circuits := slices.Collect(maps.Values(uf.sizes()))
		slices.SortFunc(circuits, func(a, b int) int { return cmp.Compare(b, a) })

		if len(circuits) < 3 {
			return "", fmt.Errorf("expected at least 3 circuits, got %d", len(circuits))
		}

		result := xslices.Product(circuits[:3])
		return strconv.Itoa(result), nil
	}
}

func day08p02(r io.Reader) (string, error) {
	boxes, err := parseJunctionBoxes(r)
	if err != nil {
		return "", err
	}

	boxPairs := buildBoxPairs(boxes)

	uf := newUnionFind(len(boxes))
	var lastPair boxPair
	for _, pair := range boxPairs {
		if uf.union(pair.indexA, pair.indexB) {
			lastPair = pair
			if uf.count() == 1 {
				break
			}
		}
	}

	result := boxes[lastPair.indexA].X * boxes[lastPair.indexB].X
	return strconv.Itoa(result), nil
}
