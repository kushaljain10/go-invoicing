package workerpool

import (
	"sync"
)

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int, Wg *sync.WaitGroup) {
	// First initialize the channel we are going to put the workers' work into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// now create all of our workers
	for i := 0; i < nworkers; i++ {
		worker := NewWorker(i+1, WorkerQueue, Wg)
		worker.Start()
	}

	go func() {
		for {
			work := <-WorkQueue
			go func() {
				worker := <-WorkerQueue
				worker <- work
			}()
		}
	}()
}
