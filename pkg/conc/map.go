// Package conc provides parallel evaluation of functions.
package conc

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// Map returns f applied to each element of in, in input order. It calls f
// concurrently from up to GOMAXPROCS goroutines, so f must be safe for
// concurrent use.
func Map[In, Out any](f func(In) Out, in []In) []Out {
	res := make([]Out, len(in))
	var idx atomic.Int64

	var wg sync.WaitGroup
	for range min(runtime.GOMAXPROCS(0), len(in)) {
		wg.Go(func() {

			for {
				inIdx := int(idx.Add(1) - 1)
				if inIdx >= len(in) {
					return
				}

				res[inIdx] = f(in[inIdx])
			}
		})
	}
	wg.Wait()
	return res
}
