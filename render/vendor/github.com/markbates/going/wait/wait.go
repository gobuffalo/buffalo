package wait

import "sync"

// Wait cleans up the pattern around using sync.WaitGroup
func Wait(length int, block func(index int)) {
	var w sync.WaitGroup
	w.Add(length)
	for i := 0; i < length; i++ {
		go func(w *sync.WaitGroup, index int) {
			block(index)
			w.Done()
		}(&w, i)
	}
	w.Wait()
}
