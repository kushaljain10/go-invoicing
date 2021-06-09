package inventory

import (
	"errors"
	"strconv"

	"github.com/kushaljain/go-invoicing/cache"
	"github.com/kushaljain/go-invoicing/utilities"
)

type Inventory struct {
	products map[string]ProductValues
}

type ProductValues struct {
	Price float64
	Cgst  int
	Stock int
}

func NewInventory() *Inventory {
	return &Inventory{
		products: make(map[string]ProductValues),
	}
}

func GetInventory(cache *cache.RedisCache) (*Inventory, error) {
	inventory := GetInventoryFromCache(cache, "inventory")
	if inventory != nil {
		return inventory, nil
	}

	inventory = NewInventory()
	err := inventory.fetchProducts()
	if err != nil {
		return nil, err
	}

	SetInventoryInCache(cache, "inventory", inventory)

	return inventory, nil
}

func (inventory *Inventory) fetchProducts() error {
	reader, err := utilities.GetCSVReader("input/products.csv")
	if utilities.IsError(err) {
		return err
	}
	headerIndices, err := utilities.GetCSVHeaderIndices(reader)
	if utilities.IsError(err) {
		return errors.New("error while fetching inventory details - " + err.Error())
	}

	for {
		currentProduct, err := reader.Read()
		if utilities.IsEOF(err) {
			break
		}
		if utilities.IsError(err) {
			return err
		}

		productName := currentProduct[headerIndices["ProductName"]]
		if !utilities.MatchRegex(productName, utilities.ALPHANUMERIC_REGEX) {
			return errors.New("Invalid product name - " + productName)
		}

		price, err := strconv.ParseFloat(currentProduct[headerIndices["ProductPrice"]], 64)
		if utilities.IsError(err) || !utilities.IsPositiveFloat(price) {
			return errors.New("Invalid price for the product - " + currentProduct[0])
		}

		cgst, err := strconv.Atoi(currentProduct[headerIndices["CGST"]])
		if utilities.IsError(err) || !utilities.IsPositiveInt(cgst) {
			return errors.New("Invalid CGST for the product - " + currentProduct[0])
		}

		stock, err := strconv.Atoi(currentProduct[headerIndices["Quantity"]])
		if utilities.IsError(err) || !utilities.IsPositiveInt(stock) {
			return errors.New("Invalid stock quantity for the product - " + currentProduct[0])
		}

		inventory.products[productName] = ProductValues{
			Price: price,
			Cgst:  cgst,
			Stock: stock,
		}
	}
	return nil
}

func (i *Inventory) UpdateProductStock(cache *cache.RedisCache, productName string, updateBy int) {
	i.products[productName] = ProductValues{
		Price: i.products[productName].Price,
		Cgst:  i.products[productName].Cgst,
		Stock: i.products[productName].Stock + updateBy,
	}
	updateProductInCache(cache, "inventory", productName, i.products[productName])
}

func (i *Inventory) GetStock(productName string) int {
	return i.products[productName].Stock
}

func (i *Inventory) GetPrice(productName string) float64 {
	return i.products[productName].Price
}

func (i *Inventory) GetCgst(productName string) int {
	return i.products[productName].Cgst
}
