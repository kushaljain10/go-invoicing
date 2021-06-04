package main

import "fmt"

type Worker struct {
	ID          int
	Work        chan Customer
	WorkerQueue chan chan Customer
	QuitChan    chan bool
}

func NewWorker(id int, workerQueue chan chan Customer) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan Customer),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
	return worker
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				generateInvoice(InventoryData, TaxesData, work, Mu)
				fmt.Println("invoicing of", work.name, "done by", w.ID)
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
