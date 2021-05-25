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

func generateInvoice(customer Customer, inventory Inventory) {

}
