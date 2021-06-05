package main

import (
	"log"
	"sync"

	"github.com/kushaljain/go-invoicing/customer"
	"github.com/kushaljain/go-invoicing/inventory"
	"github.com/kushaljain/go-invoicing/invoice"
	"github.com/kushaljain/go-invoicing/taxes"
	"github.com/kushaljain/go-invoicing/utilities"
	"github.com/kushaljain/go-invoicing/workerpool"
)

var (
	InventoryData *inventory.Inventory
	TaxesData     *taxes.Taxes
	Discounts     map[string]float64
	customers     []customer.Customer
	Mu            = &sync.Mutex{}
	Wg            = sync.WaitGroup{}
)

func init() {
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

	customers, err = customer.FetchCustomers()
	if utilities.IsError(err) {
		log.Fatalln(err)
	}
}

type generateInvoiceWork struct {
	customer customer.Customer
}

func (w generateInvoiceWork) Process() {
	err := invoice.GenerateInvoice(InventoryData, TaxesData, w.customer, Discounts, Mu, &Wg)
	if utilities.IsError(err) {
		log.Fatalln(err)
	}
}

func main() {

	workerpool.StartDispatcher(3, &Wg)

	for _, customer := range customers {
		Wg.Add(1)
		workerpool.Collector(&generateInvoiceWork{customer: customer})
	}

	Wg.Wait()
}
