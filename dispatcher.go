package main

import "fmt"

var WorkerQueue chan chan Customer

func StartDispatcher(nworkers int) {
	// First initialize the channel we are going to put the workers' work into.
	WorkerQueue = make(chan chan Customer, nworkers)

	// now create all of our workers
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
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
