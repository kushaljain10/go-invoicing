package main

import "log"

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

	for _, customer := range customers {
		invoice, err := generateInvoice(customer, inventory, taxes)
		if isError(err) {
			log.Fatalln(err)
		}
		err = invoice.Print()
		if isError(err) {
			log.Fatalln(err)
		}
	}

}
