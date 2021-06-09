package invoice

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/kushaljain/go-invoicing/cache"
	"github.com/kushaljain/go-invoicing/customer"
	"github.com/kushaljain/go-invoicing/inventory"
	"github.com/kushaljain/go-invoicing/taxes"
	"github.com/kushaljain/go-invoicing/utilities"
)

type Invoice struct {
	customerName     string
	items            []invoiceItem
	unavailableItems []string
	discount         float64
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

// var outputDirectory = ""

func GenerateInvoice(customer customer.Customer, discounts map[string]float64, mutex *sync.Mutex, Wg *sync.WaitGroup, cache *cache.RedisCache) error {
	tax, taxErr := taxes.GetTaxes(cache)
	inv, invErr := inventory.GetInventory(cache)
	if err := utilities.CheckErrors(taxErr, invErr); err != nil {
		return err
	}

	invoice := Invoice{
		customerName:     customer.Name,
		unavailableItems: make([]string, 0),
		discount:         0,
	}

	discount, ok := discounts[customer.PaymentMode]
	if ok {
		if discount <= 0 {
			return errors.New("Invalid discount amount for customer - " + invoice.customerName)
		}
		invoice.discount = discount
	}

	cart := customer.Cart
	sgst, ok := tax.SGSTList[customer.State]
	if !ok {
		return fmt.Errorf("state code '%s' for customer '%s' does not exist", customer.State, customer.Name)
	}

	mutex.Lock()
	for _, product := range cart.Items {
		item := invoiceItem{
			productName: product.Name,
		}

		if product.Quantity > inv.GetStock(item.productName) {
			invoice.unavailableItems = append(invoice.unavailableItems, item.productName)
			continue
		}
		item.quantity = product.Quantity
		inv.UpdateProductStock(cache, item.productName, (-1)*item.quantity)
		item.price = inv.GetPrice(item.productName)
		item.totalBeforeTax = float64(item.quantity) * item.price
		item.cgst = inv.GetCgst(item.productName)
		item.cgstValue = item.totalBeforeTax * (float64(item.cgst) / 100)
		item.sgst = sgst
		item.sgstValue = item.totalBeforeTax * (float64(item.sgst) / 100)
		item.totalAfterTax = item.totalBeforeTax + item.sgstValue + item.cgstValue

		invoice.items = append(invoice.items, item)
	}
	mutex.Unlock()
	err := invoice.Print()
	if utilities.IsError(err) {
		return err
	}
	Wg.Done()
	return nil
}

func (inv Invoice) Print() error {
	var totalCartValue float64

	file, err := os.Create("output/invoices/" + inv.customerName + "_invoice.txt")
	if utilities.IsError(err) {
		return errors.New("Error in creating invoice file of - " + inv.customerName)
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
	if inv.discount > 0 {
		discountAmount := totalCartValue * (inv.discount / 100)
		totalAfterDiscount := totalCartValue - discountAmount
		invoiceContent += fmt.Sprintf("\nTotal Cart Value: %.2f", totalCartValue)
		invoiceContent += fmt.Sprintf("\nDiscount (%.1f%%): %.2f", inv.discount, discountAmount)
		invoiceContent += fmt.Sprintf("\nTotal: %.f", totalAfterDiscount)
	} else {
		invoiceContent += fmt.Sprintf("\nTotal: %.2f", totalCartValue)
	}
	_, err = file.WriteString(invoiceContent)
	if utilities.IsError(err) {
		return errors.New("Error in writing to invoice file of - " + inv.customerName)
	}
	return nil
}
