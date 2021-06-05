package inventory

import (
	"errors"
	"strconv"

	"github.com/kushaljain/go-invoicing/utilities"
)

type Inventory struct {
	Products map[string]ProductValues
}

type ProductValues struct {
	Price float64
	Cgst  int
	Stock int
}

func GetInventory() (*Inventory, error) {
	inventory := &Inventory{
		Products: make(map[string]ProductValues),
	}

	reader, err := utilities.GetCSVReaderWithoutHeader("input/products.csv")
	if utilities.IsError(err) {
		return nil, err
	}

	for {
		currentProduct, err := reader.Read()
		if utilities.IsEOF(err) {
			break
		}
		if utilities.IsError(err) {
			return nil, err
		}

		productName := currentProduct[0]
		if !utilities.IsAlphaNumeric(productName) {
			return nil, errors.New("Invalid product name - " + productName)
		}

		price, err := strconv.ParseFloat(currentProduct[1], 64)
		if utilities.IsError(err) || !utilities.IsPositiveFloat(price) {
			return nil, errors.New("Invalid price for the product - " + currentProduct[0])
		}

		cgst, err := strconv.Atoi(currentProduct[2])
		if utilities.IsError(err) || !utilities.IsPositiveInt(cgst) {
			return nil, errors.New("Invalid CGST for the product - " + currentProduct[0])
		}

		stock, err := strconv.Atoi(currentProduct[3])
		if utilities.IsError(err) || !utilities.IsPositiveInt(stock) {
			return nil, errors.New("Invalid stock quantity for the product - " + currentProduct[0])
		}

		inventory.Products[productName] = ProductValues{
			Price: price,
			Cgst:  cgst,
			Stock: stock,
		}
	}
	return inventory, nil
}
