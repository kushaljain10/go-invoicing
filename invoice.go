package main

import (
	"fmt"
	"os"
)

type Invoice struct {
	customerName string
	items        []invoiceItem
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

func generateInvoice(inventory *Inventory, taxes *Taxes, customersChannel <-chan Customer, invoicesChannel chan<- Invoice) {
	for customer := range customersChannel {
		invoice := Invoice{
			customerName: customer.name,
		}
		cart := customer.cart
		sgst := taxes.SGSTList[customer.state]

		for _, product := range cart.items {
			item := invoiceItem{
				productName: product.name,
			}

			item.quantity = product.quantity
			item.price = inventory.products[item.productName].price
			item.totalBeforeTax = float64(item.quantity) * item.price
			item.cgst = inventory.products[item.productName].cgst
			item.cgstValue = item.totalBeforeTax * (float64(item.cgst) / 100)
			item.sgst = sgst
			item.sgstValue = item.totalBeforeTax * (float64(item.sgst) / 100)
			item.totalAfterTax = item.totalBeforeTax + item.sgstValue + item.cgstValue

			invoice.items = append(invoice.items, item)
		}
		invoicesChannel <- invoice
	}
}

func (inv Invoice) Print() error {
	var totalCartValue float64

	file, err := os.Create("invoices/" + inv.customerName + "_invoice.txt")
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
	invoiceContent += fmt.Sprintf("\nTotal: %.2f", totalCartValue)
	_, err = file.WriteString(invoiceContent)
	if isError(err) {
		return err
	}
	return nil
}
