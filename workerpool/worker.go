package workerpool

import (
	"sync"
)

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
	WaitGroup   *sync.WaitGroup
}

type WorkRequest interface {
	Process()
}

func NewWorker(id int, workerQueue chan chan WorkRequest, Wg *sync.WaitGroup) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		WaitGroup:   Wg,
	}
	return worker
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				work.Process()
			case <-w.QuitChan:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
