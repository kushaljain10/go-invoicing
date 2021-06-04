package main

import (
	"errors"
	"strconv"
)

type Inventory struct {
	products map[string]ProductValues
}

type ProductValues struct {
	price float64
	cgst  int
	stock int
}

func getInventory() (*Inventory, error) {
	inventory := &Inventory{
		products: make(map[string]ProductValues),
	}

	reader, err := getCSVReaderWithoutHeader("input/products.csv")
	if isError(err) {
		return nil, err
	}

	for {
		currentProduct, err := reader.Read()
		if isEOF(err) {
			break
		}
		if isError(err) {
			return nil, err
		}

		productName := currentProduct[0]
		if !isAlphaNumeric(productName) {
			return nil, errors.New("Invalid product name - " + productName)
		}

		price, err := strconv.ParseFloat(currentProduct[1], 64)
		if isError(err) || !isPositiveFloat(price) {
			return nil, errors.New("Invalid price for the product - " + currentProduct[0])
		}

		cgst, err := strconv.Atoi(currentProduct[2])
		if isError(err) || !isPositiveInt(cgst) {
			return nil, errors.New("Invalid CGST for the product - " + currentProduct[0])
		}

		stock, err := strconv.Atoi(currentProduct[3])
		if isError(err) || !isPositiveInt(stock) {
			return nil, errors.New("Invalid stock quantity for the product - " + currentProduct[0])
		}

		inventory.products[productName] = ProductValues{
			price: price,
			cgst:  cgst,
			stock: stock,
		}
	}
	return inventory, nil
}
