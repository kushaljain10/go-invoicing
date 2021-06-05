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

	reader, err := utilities.GetCSVReader(inputDirectory + "/customers.csv")
	if utilities.IsError(err) {
		return nil, errors.New("Error while fetching the customer details - " + err.Error())
	}
	headerIndices, err := utilities.GetCSVHeaderIndices(reader)
	if utilities.IsError(err) {
		return nil, errors.New("error while fetching customer details - " + err.Error())
	}

	for {
		customer, err := reader.Read()
		if utilities.IsEOF(err) {
			break
		}
		if utilities.IsError(err) {
			return nil, errors.New("error in reading from CSV")
		}

		customerName, nameErr := utilities.ValidateString(customer[headerIndices["CustomerName"]], "^[a-zA-Z0-9_ ]*$")
		customerState, stateErr := utilities.ValidateString(customer[headerIndices["CustomerState"]], "^[A-Z]*$")
		paymentMode, paymentErr := utilities.ValidateString(customer[headerIndices["PaymentMode"]], "^[a-zA-Z]*$")

		if err := utilities.CheckErrors(nameErr, stateErr, paymentErr); err != nil {
			return nil, err
		}

		cart, err := fetchCart(customerName)
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

func fetchCart(customerName string) (*Cart, error) {
	cart := &Cart{
		Items: make([]cartItem, 0),
	}

	reader, err := utilities.GetCSVReader(cartDirectory + "/" + customerName + ".csv")
	if utilities.IsError(err) {
		return nil, errors.New("Error while fetching the cart for - " + customerName + ": " + err.Error())
	}
	headerIndices, err := utilities.GetCSVHeaderIndices(reader)
	if utilities.IsError(err) {
		return nil, errors.New("error while fetching the cart for - " + customerName + ": " + err.Error())
	}

	for {
		item, err := reader.Read()
		if utilities.IsEOF(err) {
			break
		}
		if utilities.IsError(err) {
			return nil, err
		}

		productName, nameErr := utilities.ValidateString(item[headerIndices["ProductName"]], utilities.ALPHANUMERIC_REGEX)
		if utilities.IsError(nameErr) {
			return nil, errors.New("error while fetching product name - " + err.Error())
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
