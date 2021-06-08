package main

import (
	"log"
	"sync"

	"github.com/kushaljain/go-invoicing/cache"
	"github.com/kushaljain/go-invoicing/customer"
	"github.com/kushaljain/go-invoicing/inventory"
	"github.com/kushaljain/go-invoicing/invoice"
	"github.com/kushaljain/go-invoicing/taxes"
	"github.com/kushaljain/go-invoicing/utilities"
	"github.com/kushaljain/go-invoicing/workerpool"
)

var (
	inventoryData *inventory.Inventory
	taxesData     *taxes.Taxes
	discounts     map[string]float64
	customers     []customer.Customer
	mu            = &sync.Mutex{}
	wg            = sync.WaitGroup{}
	invCache      cache.Cache
)

func init() {
	invCache = cache.NewRedisCache("localhost:6379", 0, 30)

	var err error
	inventoryData, err = inventory.GetInventory()
	if utilities.IsError(err) {
		log.Fatalln(err)
	}

	taxesData = invCache.Get("taxes")
	if taxesData == nil {
		taxesData = taxes.NewTaxes()
		err = taxesData.GetSGSTList()
		if utilities.IsError(err) {
			log.Fatalln(err)
		}
		invCache.Set("taxes", taxesData)
	}

	discounts = map[string]float64{"UPI": 5}

	customers, err = customer.FetchCustomers()
	if utilities.IsError(err) {
		log.Fatalln(err)
	}
}

type generateInvoiceWork struct {
	customer customer.Customer
}

func (w generateInvoiceWork) Process() {
	err := invoice.GenerateInvoice(inventoryData, taxesData, w.customer, discounts, mu, &wg)
	if utilities.IsError(err) {
		log.Fatalln(err)
	}
}

func main() {

	workerpool.StartDispatcher(3, &wg)

	for _, customer := range customers {
		wg.Add(1)
		workerpool.Collector(&generateInvoiceWork{customer: customer})
	}

	wg.Wait()
}
