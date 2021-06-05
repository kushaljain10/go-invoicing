package main

import (
	"sync"

	"github.com/kushaljain/go-invoicing/customer"
)

var WorkerQueue chan chan customer.Customer

func StartDispatcher(nworkers int, Wg *sync.WaitGroup) {
	// First initialize the channel we are going to put the workers' work into.
	WorkerQueue = make(chan chan customer.Customer, nworkers)

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
