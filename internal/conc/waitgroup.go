package conc

import "sync"

type WaitGroup struct {
	wg sync.WaitGroup
}

func (w *WaitGroup) Go(f func()) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		f()
	}()
}

func (w *WaitGroup) Wait() {
	w.wg.Wait()
}
