package inventory

import (
	"errors"
	"strconv"

	"github.com/kushaljain/go-invoicing/utilities"
)

type Inventory struct {
	products map[string]ProductValues
}

type ProductValues struct {
	price float64
	cgst  int
	stock int
}

func GetInventory() (*Inventory, error) {
	inventory := &Inventory{
		products: make(map[string]ProductValues),
	}

	reader, err := utilities.GetCSVReader("input/products.csv")
	if utilities.IsError(err) {
		return nil, err
	}
	headerIndices, err := utilities.GetCSVHeaderIndices(reader)
	if utilities.IsError(err) {
		return nil, errors.New("error while fetching inventory details - " + err.Error())
	}

	for {
		currentProduct, err := reader.Read()
		if utilities.IsEOF(err) {
			break
		}
		if utilities.IsError(err) {
			return nil, err
		}

		productName := currentProduct[headerIndices["ProductName"]]
		if !utilities.MatchRegex(productName, utilities.ALPHANUMERIC_REGEX) {
			return nil, errors.New("Invalid product name - " + productName)
		}

		price, err := strconv.ParseFloat(currentProduct[headerIndices["ProductPrice"]], 64)
		if utilities.IsError(err) || !utilities.IsPositiveFloat(price) {
			return nil, errors.New("Invalid price for the product - " + currentProduct[0])
		}

		cgst, err := strconv.Atoi(currentProduct[headerIndices["CGST"]])
		if utilities.IsError(err) || !utilities.IsPositiveInt(cgst) {
			return nil, errors.New("Invalid CGST for the product - " + currentProduct[0])
		}

		stock, err := strconv.Atoi(currentProduct[headerIndices["Quantity"]])
		if utilities.IsError(err) || !utilities.IsPositiveInt(stock) {
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

func (i *Inventory) UpdateProductStock(productName string, updateBy int) {
	i.products[productName] = ProductValues{
		price: i.products[productName].price,
		cgst:  i.products[productName].cgst,
		stock: i.products[productName].stock + updateBy,
	}
}

func (i *Inventory) GetStock(productName string) int {
	return i.products[productName].stock
}

func (i *Inventory) GetPrice(productName string) float64 {
	return i.products[productName].price
}

func (i *Inventory) GetCgst(productName string) int {
	return i.products[productName].cgst
}
