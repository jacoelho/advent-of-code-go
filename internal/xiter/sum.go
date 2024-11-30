package xiter

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/internal/xconstraints"
)

func Sum[V xconstraints.Number](seq iter.Seq[V]) V {
	var total V
	for v := range seq {
		total += v
	}
	return total
}
