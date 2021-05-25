package main

import (
	"encoding/csv"
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

		cart, err := fetchCart(customerName)
		if isError(err) {
			return nil, err
		}

		customers = append(customers, Customer{
			name:  c[0],
			state: c[1],
			cart:  *cart,
		})
	}
	return customers, nil
}

func fetchCart(customerName string) (*Cart, error) {
	var cart *Cart

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

		quantity, err := strconv.Atoi(c[1])
		if isError(err) {
			return nil, err
		}

		cart.items = append(cart.items, cartItem{
			name:     c[0],
			quantity: quantity,
		})
	}
	return cart, nil
}
