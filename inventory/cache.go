package inventory

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kushaljain/go-invoicing/cache"
)

func SetInventoryInCache(cache *cache.RedisCache, key string, inventory *Inventory) {
	ctx := context.TODO()
	client := cache.GetClient()

	for name, product := range inventory.products {
		json, err := json.Marshal(product)
		if err != nil {
			panic(err)
		}

		client.HSet(ctx, key, name, json)
		client.Expire(ctx, key, cache.Expires*time.Second)
	}
}

func GetInventoryFromCache(cache *cache.RedisCache, key string) *Inventory {
	ctx := context.TODO()
	client := cache.GetClient()

	val, err := client.HGetAll(ctx, key).Result()
	if err != nil || len(val) == 0 {
		return nil
	}

	inv := NewInventory()
	var currProduct ProductValues
	for name, product := range val {
		err = json.Unmarshal([]byte(product), &currProduct)
		if err != nil {
			return nil
		}
		inv.products[name] = currProduct
	}
	return inv
}

func updateProductInCache(cache *cache.RedisCache, key string, name string, product ProductValues) {
	ctx := context.TODO()
	client := cache.GetClient()

	json, err := json.Marshal(product)
	if err != nil {
		panic(err)
	}

	client.HSet(ctx, key, name, json)
}
