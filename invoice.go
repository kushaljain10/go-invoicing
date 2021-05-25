package main

type Invoice struct {
	customerName  string
	customerState string
	customerSGST  int
	items         []invoiceItem
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

func generateInvoice(customer Customer, inventory Inventory, taxes Taxes) Invoice, error {
	var err error
	var invoice Invoice
	cart := customer.cart
	sgst := taxes.SGSTList[customer.state]
	
	for _, product := range cart {
		item := & invoiceItem{
			productName: product.name
		}

		item.quantity, err = strconv.Atoi(product.quantity)
		if isError(err) {
			return nil, err
		}

		item.price = inventory.products[item.productName].price
		item.totalBeforeTax = float(item.quantity) * item.price
		item.cgst = inventory.products[item.productName].cgst
		item.cgstValue = item.totalBeforeTax * (float(item.cgst) / 100)
		item.sgst = sgst
		item.sgstValue = item.totalBeforeTax * (float(item.sgst) / 100)
		item.totalAfterTax = item.totalBeforeTax + item.sgstValue + item.cgstValue
		
		invoice.items = append(invoice.items, item)
	}
}

func (inv Invoice) Print() err {
	var totalCartValue float64

	file, err := os.Create("invoices/" + inv.customerName + "_invoice.txt")
	if isError(err) {
		return err
	}
	defer file.Close()

	invoiceContent := "Customer Name: " + inv.customerName + "\n\n"
	for i, item := range inv.items {
		invoiceContent +=
			fmt.Sprintf("%d. %s | ", i+1, item,productName) +
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

}
