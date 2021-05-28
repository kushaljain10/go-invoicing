package main

import "strconv"

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

	r, err := getCSVReaderWithoutHeader("products.csv")
	if isError(err) {
		return nil, err
	}

	for {
		p, err := r.Read()
		if isEOF(err) {
			break
		}
		if isError(err) {
			return nil, err
		}

		price, err := strconv.ParseFloat(p[1], 64)
		if isError(err) {
			return nil, err
		}

		cgst, err := strconv.Atoi(p[2])
		if isError(err) {
			return nil, err
		}

		stock, err := strconv.Atoi(p[3])
		if isError(err) {
			return nil, err
		}

		inventory.products[p[0]] = ProductValues{
			price: price,
			cgst:  cgst,
			stock: stock,
		}
	}
	return inventory, nil
}
