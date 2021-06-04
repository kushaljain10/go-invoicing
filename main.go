package main

import (
	"log"
	"sync"
)

var WorkQueue = make(chan Customer, 100)
var InventoryData *Inventory
var TaxesData *Taxes
var Mu = &sync.Mutex{}
var Wg = sync.WaitGroup{}

func main() {
	var err error
	InventoryData, err = getInventory()
	if isError(err) {
		log.Fatalln(err)
	}

	TaxesData = NewTaxes()
	err = TaxesData.getSGSTList()
	if isError(err) {
		log.Fatalln(err)
	}

	customers, err := fetchCustomers()
	if isError(err) {
		log.Fatalln(err)
	}

	StartDispatcher(3)
	for _, customer := range customers {
		Wg.Add(1)
		WorkQueue <- customer
	}

	Wg.Wait()

}
