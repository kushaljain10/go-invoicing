package main

import (
	"log"
	"sync"
)

func main() {
	inventory, err := getInventory()
	if isError(err) {
		log.Fatalln(err)
	}

	taxes := NewTaxes()
	err = taxes.getSGSTList()
	if isError(err) {
		log.Fatalln(err)
	}

	customers, err := fetchCustomers()
	if isError(err) {
		log.Fatalln(err)
	}

	mutex := &sync.Mutex{}

	numOfCustomers := len(customers)
	customersChannel := make(chan Customer, numOfCustomers)
	invoicesChannel := make(chan Invoice, numOfCustomers)

	go generateInvoice(inventory, taxes, customersChannel, invoicesChannel, mutex)
	go generateInvoice(inventory, taxes, customersChannel, invoicesChannel, mutex)
	go generateInvoice(inventory, taxes, customersChannel, invoicesChannel, mutex)

	for _, c := range customers {
		customersChannel <- c
	}
	close(customersChannel)

	for i := 0; i < numOfCustomers; i++ {
		invoice := <-invoicesChannel
		err = invoice.Print()
		if isError(err) {
			log.Fatalln(err)
		}
	}

}
