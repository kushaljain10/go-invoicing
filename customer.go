package main

import (
	"encoding/csv"
	"errors"
	"os"
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

func fetchCustomers() ([]Customer, error) {

	var customers []Customer

	r, err := getCSVReaderWithoutHeader("customers.csv")
	if isError(err) {
		return nil, err
	}

	for {
		c, err := r.Read()
		if isEOF(err) {
			break
		}
		if isError(err) {
			return nil, err
		}

		customerName := c[0]
		if isEmptyString(customerName) {
			return nil, errors.New("error while fetching customer details - no customer name found")
		}

		customerState := c[1]
		if isEmptyString(customerState) {
			return nil, errors.New("error while fetching customer details - no customer state found")
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

	file, err := os.Open("cart/" + customerName + ".csv")
	if isError(err) {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	_, err = r.Read()
	if isError(err) {
		return nil, err
	}

	for {
		c, err := r.Read()
		if isEOF(err) {
			break
		}
		if isError(err) {
			return nil, err
		}

		productName := c[0]
		if isEmptyString(productName) {
			return nil, errors.New("error while fetching cart details - no product name found")
		}

		quantity, err := strconv.Atoi(c[1])
		if isError(err) {
			return nil, err
		}
		if !isPositiveInt(quantity) {
			return nil, errors.New("invalid product quantity for " + productName + " for customer - " + customerName)
		}

		cart.items = append(cart.items, cartItem{
			name:     productName,
			quantity: quantity,
		})
	}
	return cart, nil
}
