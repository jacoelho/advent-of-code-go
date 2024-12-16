package search

import (
	"slices"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
)

func ConstantStepCost[T any](_, _ T) int {
	return 1
}

type minItem[T any] struct {
	item     T
	priority int
}

// AStar performs the A* search algorithm.
func AStar[T comparable](
	start T,
	neighbours func(T) []T,
	heuristic func(T) int,
	stepCost func(T, T) int,
) (int, []T, bool) {
	priorityQueue := collections.NewHeap[minItem[T]](func(m minItem[T], m2 minItem[T]) bool {
		return m.priority < m2.priority
	})

	priorityQueue.Push(minItem[T]{
		item:     start,
		priority: heuristic(start),
	})

	previous := make(map[T]*T)
	pathCost := make(map[T]int)
	pathCost[start] = 0

	for current := range priorityQueue.PopSeq() {
		// Check if the current state is the goal.
		if heuristic(current.item) == 0 {
			var path []T
			for cur := &current.item; cur != nil; cur = previous[*cur] {
				path = append(path, *cur)
			}
			slices.Reverse(path)
			return pathCost[current.item], path, true
		}

		for _, neighbor := range neighbours(current.item) {
			newCost := pathCost[current.item] + stepCost(current.item, neighbor)
			if oldCost, ok := pathCost[neighbor]; !ok || newCost < oldCost {
				pathCost[neighbor] = newCost
				priorityQueue.Push(minItem[T]{
					item:     neighbor,
					priority: newCost + heuristic(neighbor),
				})
				previous[neighbor] = &current.item
			}
		}
	}

	return 0, nil, false
}

// AStarBag visits all paths with the lowest cost
func AStarBag[T comparable](
	start T,
	neighbours func(T) []T,
	heuristic func(T) int,
	stepCost func(T, T) int,
) (int, [][]T, bool) {
	type minItem[T any] struct {
		item     T
		priority int
	}

	priorityQueue := collections.NewHeap[minItem[T]](func(m minItem[T], m2 minItem[T]) bool {
		return m.priority < m2.priority
	})

	priorityQueue.Push(minItem[T]{
		item:     start,
		priority: heuristic(start),
	})

	paths := make(map[T][][]T)
	paths[start] = [][]T{{start}}

	pathCost := make(map[T]int)
	pathCost[start] = 0

	var goalPaths [][]T
	lowestGoalCost := -1

	for current := range priorityQueue.PopSeq() {
		if heuristic(current.item) == 0 {
			currentCost := pathCost[current.item]
			if lowestGoalCost == -1 || currentCost == lowestGoalCost {
				lowestGoalCost = currentCost
				goalPaths = append(goalPaths, paths[current.item]...)
			} else if currentCost > lowestGoalCost {
				break
			}
			continue
		}

		for _, neighbor := range neighbours(current.item) {
			newCost := pathCost[current.item] + stepCost(current.item, neighbor)
			if oldCost, ok := pathCost[neighbor]; !ok || newCost < oldCost {
				// Found a better path; update costs and paths
				pathCost[neighbor] = newCost
				priorityQueue.Push(minItem[T]{
					item:     neighbor,
					priority: newCost + heuristic(neighbor),
				})

				var newPaths [][]T
				for _, path := range paths[current.item] {
					newPath := append([]T{}, path...)
					newPath = append(newPath, neighbor)
					newPaths = append(newPaths, newPath)
				}
				paths[neighbor] = newPaths
			} else if newCost == oldCost {
				for _, path := range paths[current.item] {
					newPath := append([]T{}, path...)
					newPath = append(newPath, neighbor)
					paths[neighbor] = append(paths[neighbor], newPath)
				}
			}
		}
	}

	if len(goalPaths) > 0 {
		return lowestGoalCost, goalPaths, true
	}

	return 0, nil, false
}
