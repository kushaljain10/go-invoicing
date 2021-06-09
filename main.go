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
	discounts                 map[string]float64
	customers                 []customer.Customer
	mu                        = &sync.Mutex{}
	wg                        = sync.WaitGroup{}
	invCache                  *cache.RedisCache
	CACHE_DURATION_IN_SECONDS = 30
)

func init() {
	invCache = cache.NewRedisCache("localhost:6379", 0, 30)

	var err error

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
