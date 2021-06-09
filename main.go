package main

import (
	"log"
	"sync"

	"github.com/kushaljain/go-invoicing/cache"
	"github.com/kushaljain/go-invoicing/customer"
	"github.com/kushaljain/go-invoicing/invoice"
	"github.com/kushaljain/go-invoicing/utilities"
	"github.com/kushaljain/go-invoicing/workerpool"
)

var (
	// inventoryData *inventory.Inventory
	// taxesData     *taxes.Taxes
	discounts map[string]float64
	customers []customer.Customer
	mu        = &sync.Mutex{}
	wg        = sync.WaitGroup{}
	invCache  *cache.RedisCache
)

func init() {
	invCache = cache.NewRedisCache("localhost:6379", 0, 30)

	var err error
	// inventoryData, err = inventory.GetInventory(invCache)
	// if utilities.IsError(err) {
	// 	log.Fatalln(err)
	// }

	// taxesData, err = taxes.GetTaxes(invCache)
	// if utilities.IsError(err) {
	// 	log.Fatalln(err)
	// }

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
	err := invoice.GenerateInvoice(w.customer, discounts, mu, &wg, invCache)
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
