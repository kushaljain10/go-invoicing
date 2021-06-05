package main

import (
	"log"
	"sync"

	"github.com/kushaljain/go-invoicing/customer"
	"github.com/kushaljain/go-invoicing/inventory"
	"github.com/kushaljain/go-invoicing/taxes"
	"github.com/kushaljain/go-invoicing/utilities"
)

var WorkQueue = make(chan customer.Customer, 100)
var InventoryData *inventory.Inventory
var TaxesData *taxes.Taxes
var Discounts map[string]float64
var Mu = &sync.Mutex{}
var Wg = sync.WaitGroup{}

func main() {
	var err error
	InventoryData, err = inventory.GetInventory()
	if utilities.IsError(err) {
		log.Fatalln(err)
	}

	TaxesData = taxes.NewTaxes()
	err = TaxesData.GetSGSTList()
	if utilities.IsError(err) {
		log.Fatalln(err)
	}

	Discounts = map[string]float64{"UPI": 5}

	customers, err := customer.FetchCustomers()
	if utilities.IsError(err) {
		log.Fatalln(err)
	}

	StartDispatcher(3, &Wg)

	for _, customer := range customers {
		Wg.Add(1)
		WorkQueue <- customer
	}

	Wg.Wait()

}
