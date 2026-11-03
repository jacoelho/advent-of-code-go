package graph

import (
	"fmt"
	"math"
)

// CountDAGPaths sums terminal contributions along every path from roots. A
// terminal contribution must be nonnegative; even a zero contribution
// suppresses expansion. Repeated roots and parallel edges retain their
// multiplicity.
//
// Counts are memoized across roots. Reachable cycles and int overflow return an
// error with no partial count. The children callback is called once per
// nonterminal node and must invoke yield synchronously, without retaining it,
// stopping if it returns false. Callers must provide a finite graph with stable
// node identities. Traversal uses an explicit stack, retaining pending children
// only for the active ancestry.
func CountDAGPaths[T comparable](roots []T, children func(T, func(T) bool), terminal func(T) (int, bool)) (int, error) {
	type frame struct {
		node             T
		start, next, end int
		count            int
		entered          bool
	}
	// Negative entries identify active ancestors; completed counts are nonnegative.
	counts := make(map[T]int)
	var stack []frame
	var pending []T
	appendChild := func(child T) bool {
		pending = append(pending, child)
		return true
	}
	total := 0
	for _, root := range roots {
		if count, known := counts[root]; known {
			if total > math.MaxInt-count {
				return 0, fmt.Errorf("DAG path count overflows int")
			}
			total += count
			continue
		}
		stack = append(stack, frame{node: root})
		for len(stack) > 0 {
			current := &stack[len(stack)-1]
			if !current.entered {
				current.start = len(pending)
				if count, done := terminal(current.node); done {
					if count < 0 {
						return 0, fmt.Errorf("negative DAG terminal contribution")
					}
					current.count = count
				} else {
					counts[current.node] = -1
					children(current.node, appendChild)
				}
				current.next, current.end = current.start, len(pending)
				current.entered = true
			}
			if current.next < current.end {
				child := pending[current.next]
				current.next++
				count, known := counts[child]
				if !known {
					stack = append(stack, frame{node: child})
					continue
				}
				if count < 0 {
					return 0, fmt.Errorf("cycle in DAG at %v", child)
				}
				if current.count > math.MaxInt-count {
					return 0, fmt.Errorf("DAG path count overflows int")
				}
				current.count += count
				continue
			}

			count := current.count
			counts[current.node] = count
			clear(pending[current.start:current.end])
			pending = pending[:current.start]
			clear(stack[len(stack)-1:])
			stack = stack[:len(stack)-1]
			sum := &total
			if len(stack) > 0 {
				sum = &stack[len(stack)-1].count
			}
			if *sum > math.MaxInt-count {
				return 0, fmt.Errorf("DAG path count overflows int")
			}
			*sum += count
		}
	}
	return total, nil
}
