package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Invoice struct {
	customerName     string
	items            []invoiceItem
	unavailableItems []string
}

type invoiceItem struct {
	productName    string
	quantity       int
	price          float64
	sgst           int
	sgstValue      float64
	cgst           int
	cgstValue      float64
	totalBeforeTax float64
	totalAfterTax  float64
}

func generateInvoice(inventory *Inventory, taxes *Taxes, customer Customer, mutex *sync.Mutex, wg *sync.WaitGroup) {
	// for customer := range customersChannel {
	invoice := Invoice{
		customerName:     customer.name,
		unavailableItems: make([]string, 0),
	}
	cart := customer.cart
	sgst := taxes.SGSTList[customer.state]

	mutex.Lock()
	for _, product := range cart.items {
		item := invoiceItem{
			productName: product.name,
		}

		if product.quantity > inventory.products[item.productName].stock {
			invoice.unavailableItems = append(invoice.unavailableItems, item.productName)
			continue
		}
		item.quantity = product.quantity
		inventory.products[item.productName] = ProductValues{
			price: inventory.products[item.productName].price,
			cgst:  inventory.products[item.productName].cgst,
			stock: inventory.products[item.productName].stock - item.quantity,
		}
		item.price = inventory.products[item.productName].price
		item.totalBeforeTax = float64(item.quantity) * item.price
		item.cgst = inventory.products[item.productName].cgst
		item.cgstValue = item.totalBeforeTax * (float64(item.cgst) / 100)
		item.sgst = sgst
		item.sgstValue = item.totalBeforeTax * (float64(item.sgst) / 100)
		item.totalAfterTax = item.totalBeforeTax + item.sgstValue + item.cgstValue

		invoice.items = append(invoice.items, item)
	}
	mutex.Unlock()
	err := invoice.Print()
	if isError(err) {
		log.Fatalln(err)
	}
	wg.Done()
	// invoicesChannel <- invoice
	// }
}

func (inv Invoice) Print() error {
	var totalCartValue float64

	file, err := os.Create("output/invoices/" + inv.customerName + "_invoice.txt")
	if isError(err) {
		return err
	}
	defer file.Close()

	invoiceContent := "Customer Name: " + inv.customerName + "\n\n"
	for i, item := range inv.items {
		invoiceContent +=
			fmt.Sprintf("%d. %s | ", i+1, item.productName) +
				fmt.Sprintf("Quantity: %d | ", item.quantity) +
				fmt.Sprintf("Price: %.2f | ", item.price) +
				fmt.Sprintf("SGST: %.2f | ", item.sgstValue) +
				fmt.Sprintf("CGST: %.2f | ", item.cgstValue) +
				fmt.Sprintf("Total: %.2f\n", item.totalAfterTax)
		totalCartValue += item.totalAfterTax
	}
	if len(inv.unavailableItems) > 0 {
		invoiceContent += "\nFollowing items were not in stock:\n"
		for i, item := range inv.unavailableItems {
			invoiceContent += fmt.Sprintf("%d. %s\n", i+1, item)
		}
	}
	invoiceContent += fmt.Sprintf("\nTotal: %.2f", totalCartValue)
	_, err = file.WriteString(invoiceContent)
	if isError(err) {
		return err
	}
	return nil
}
