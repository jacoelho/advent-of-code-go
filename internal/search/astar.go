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

func minHeap[T comparable]() *collections.Heap[minItem[T]] {
	return collections.NewHeap[minItem[T]](func(m minItem[T], m2 minItem[T]) bool {
		return m.priority < m2.priority
	})
}

// AStar performs the A* search algorithm.
func AStar[T comparable](
	start T,
	neighbours func(T) []T,
	heuristic func(T) int,
	stepCost func(T, T) int,
) (int, []T, bool) {
	priorityQueue := minHeap[T]()

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
	priorityQueue := minHeap[T]()
	priorityQueue.Push(minItem[T]{
		item:     start,
		priority: heuristic(start),
	})

	paths := map[T][][]T{start: {{start}}}
	pathCost := map[T]int{start: 0}
	var goalPaths [][]T
	lowestGoalCost := -1

	for current := range priorityQueue.PopSeq() {
		currItem := current.item
		currCost := pathCost[currItem]

		if heuristic(currItem) == 0 {
			if lowestGoalCost == -1 || currCost <= lowestGoalCost {
				if lowestGoalCost != currCost {
					lowestGoalCost = currCost
					goalPaths = nil
				}
				goalPaths = append(goalPaths, paths[currItem]...)
				continue
			}
			break
		}

		for _, neighbor := range neighbours(currItem) {
			newCost := currCost + stepCost(currItem, neighbor)
			oldCost, seen := pathCost[neighbor]

			if !seen || newCost < oldCost {
				// update path cost and priority queue for a better path
				pathCost[neighbor] = newCost
				priorityQueue.Push(minItem[T]{item: neighbor, priority: newCost + heuristic(neighbor)})
				paths[neighbor] = nil
			}

			if !seen || newCost == pathCost[neighbor] {
				paths[neighbor] = append(paths[neighbor], appendPaths(paths[currItem], neighbor)...)
			}
		}
	}

	if len(goalPaths) > 0 {
		return lowestGoalCost, goalPaths, true
	}

	return 0, nil, false
}

func appendPaths[T comparable](paths [][]T, node T) [][]T {
	newPaths := make([][]T, len(paths))
	for i, path := range paths {
		// Pre-allocate the new path with enough capacity for the additional node
		newPath := make([]T, len(path)+1)
		copy(newPath, path)       // Copy the existing path
		newPath[len(path)] = node // Append the new node directly
		newPaths[i] = newPath
	}
	return newPaths
}
