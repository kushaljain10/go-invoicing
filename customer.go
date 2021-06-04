package main

import (
	"errors"
	"strconv"
)

type Customer struct {
	name  string
	state string
	cart  Cart
}

type Cart struct {
	items []cartItem
}

type cartItem struct {
	name     string
	quantity int
}

// Directory Names
var inputDirectory = "input"
var cartDirectory = inputDirectory + "/cart"

func fetchCustomers() ([]Customer, error) {

	var customers []Customer

	reader, err := getCSVReaderWithoutHeader(inputDirectory + "/customers.csv")
	if isError(err) {
		return nil, err
	}

	for {
		customer, err := reader.Read()
		if isEOF(err) {
			break
		}
		if isError(err) {
			return nil, err
		}

		customerName := customer[0]
		if isEmptyString(customerName) {
			return nil, errors.New("error while fetching customer details - no customer name found")
		}
		if !matchRegex(customerName, "^[a-zA-Z0-9_ ]*$") {
			return nil, errors.New("Invalid customer name - " + customerName)
		}

		customerState := customer[1]
		if isEmptyString(customerState) {
			return nil, errors.New("error while fetching customer details - no customer state found")
		}
		if !matchRegex(customerState, "^[A-Z]*$") {
			return nil, errors.New("Invalid state code - " + customerState)
		}

		cart, err := fetchCart(customerName)
		if isError(err) {
			return nil, err
		}

		customers = append(customers, Customer{
			name:  customerName,
			state: customerState,
			cart:  *cart,
		})
	}
	return customers, nil
}

func fetchCart(customerName string) (*Cart, error) {
	cart := &Cart{
		items: make([]cartItem, 0),
	}

	reader, err := getCSVReaderWithoutHeader(cartDirectory + "/" + customerName + ".csv")
	if isError(err) {
		return nil, errors.New("Error while fetching the cart for customer - " + customerName)
	}

	for {
		item, err := reader.Read()
		if isEOF(err) {
			break
		}
		if isError(err) {
			return nil, err
		}

		productName := item[0]
		if isEmptyString(productName) {
			return nil, errors.New("error while fetching cart details - no product name found")
		}
		if !isAlphaNumeric(productName) {
			return nil, errors.New("Invalid product name - " + productName + " for customer - " + customerName)
		}

		quantity, err := strconv.Atoi(item[1])
		if isError(err) || !isPositiveInt(quantity) {
			return nil, errors.New("Invalid quantity for the product - " + productName + " for customer - " + customerName)
		}

		cart.items = append(cart.items, cartItem{
			name:     productName,
			quantity: quantity,
		})
	}
	return cart, nil
}
