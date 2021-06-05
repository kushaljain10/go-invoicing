package main

import (
	"log"
	"sync"

	"github.com/kushaljain/go-invoicing/customer"
	"github.com/kushaljain/go-invoicing/invoice"
)

type Worker struct {
	ID          int
	Work        chan customer.Customer
	WorkerQueue chan chan customer.Customer
	QuitChan    chan bool
	WaitGroup   *sync.WaitGroup
}

func NewWorker(id int, workerQueue chan chan customer.Customer, Wg *sync.WaitGroup) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan customer.Customer),
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
				err := invoice.GenerateInvoice(InventoryData, TaxesData, work, Discounts, Mu, &Wg)
				if err != nil {
					log.Fatalln(err)
				}
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
