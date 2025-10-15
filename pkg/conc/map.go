package conc

import (
	"runtime"
	"sync"
	"sync/atomic"
)

func Map[In, Out any](f func(In) Out, in []In) []Out {
	res := make([]Out, len(in))
	var idx atomic.Int64

	var wg sync.WaitGroup
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
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
