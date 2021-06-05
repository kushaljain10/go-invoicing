package customer

import (
	"errors"
	"strconv"

	"github.com/kushaljain/go-invoicing/utilities"
)

type Customer struct {
	Name        string
	State       string
	Cart        Cart
	PaymentMode string
}

type Cart struct {
	Items []cartItem
}

type cartItem struct {
	Name     string
	Quantity int
}

// Directory Names
var inputDirectory = "input"
var cartDirectory = inputDirectory + "/cart"

func FetchCustomers() ([]Customer, error) {

	var customers []Customer

	reader, err := utilities.GetCSVReaderWithoutHeader(inputDirectory + "/customers.csv")
	if utilities.IsError(err) {
		return nil, err
	}

	for {
		customer, err := reader.Read()
		if utilities.IsEOF(err) {
			break
		}
		if utilities.IsError(err) {
			return nil, err
		}

		customerName := customer[0]
		if utilities.IsEmptyString(customerName) {
			return nil, errors.New("error while fetching customer details - customer name not found")
		}
		if !utilities.MatchRegex(customerName, "^[a-zA-Z0-9_ ]*$") {
			return nil, errors.New("Invalid customer name - " + customerName)
		}

		customerState := customer[1]
		if utilities.IsEmptyString(customerState) {
			return nil, errors.New("error while fetching customer details - customer state not found")
		}
		if !utilities.MatchRegex(customerState, "^[A-Z]*$") {
			return nil, errors.New("Invalid state code - " + customerState)
		}

		paymentMode := customer[2]
		if utilities.IsEmptyString(paymentMode) {
			return nil, errors.New("error while fetching customer details - payment mode not found")
		}
		if !utilities.MatchRegex(paymentMode, "^[a-zA-Z]*$") {
			return nil, errors.New("Invalid Payment Mode - " + paymentMode)
		}

		cart, err := FetchCart(customerName)
		if utilities.IsError(err) {
			return nil, err
		}

		customers = append(customers, Customer{
			Name:        customerName,
			State:       customerState,
			Cart:        *cart,
			PaymentMode: paymentMode,
		})
	}
	return customers, nil
}

func FetchCart(customerName string) (*Cart, error) {
	cart := &Cart{
		Items: make([]cartItem, 0),
	}

	reader, err := utilities.GetCSVReaderWithoutHeader(cartDirectory + "/" + customerName + ".csv")
	if utilities.IsError(err) {
		return nil, errors.New("Error while fetching the cart for customer - " + customerName)
	}

	for {
		item, err := reader.Read()
		if utilities.IsEOF(err) {
			break
		}
		if utilities.IsError(err) {
			return nil, err
		}

		productName := item[0]
		if utilities.IsEmptyString(productName) {
			return nil, errors.New("error while fetching cart details - no product name found")
		}
		if !utilities.IsAlphaNumeric(productName) {
			return nil, errors.New("Invalid product name - " + productName + " for customer - " + customerName)
		}

		quantity, err := strconv.Atoi(item[1])
		if utilities.IsError(err) || !utilities.IsPositiveInt(quantity) {
			return nil, errors.New("Invalid quantity for the product - " + productName + " for customer - " + customerName)
		}

		cart.Items = append(cart.Items, cartItem{
			Name:     productName,
			Quantity: quantity,
		})
	}
	return cart, nil
}
