package conc

import (
	"iter"
	"math"
	"runtime"
)

// Blocks yields the outputs of scan over consecutive index blocks
// [start, start+size), [start+size, start+2*size), … in index order, ending
// at math.MaxInt. It scans a batch of blocks concurrently with Map before
// yielding any of them, so a consumer that stops early wastes at most one
// batch, and no goroutine outlives the iteration. scan must be safe for
// concurrent use. Blocks panics if size is not positive.
func Blocks[T any](start, size int, scan func(start, end int) []T) iter.Seq[T] {
	if size <= 0 {
		panic("conc: Blocks size must be positive")
	}
	return func(yield func(T) bool) {
		type block struct{ start, end int }
		batch := make([]block, 0, 4*runtime.GOMAXPROCS(0))
		scanBlock := func(b block) []T { return scan(b.start, b.end) }
		for next, done := start, false; !done; {
			batch = batch[:0]
			for len(batch) < cap(batch) && !done {
				end := math.MaxInt
				if next <= math.MaxInt-size {
					end = next + size
				}
				batch = append(batch, block{next, end})
				next, done = end, end == math.MaxInt
			}
			for _, values := range Map(scanBlock, batch) {
				for _, value := range values {
					if !yield(value) {
						return
					}
				}
			}
		}
	}
}
