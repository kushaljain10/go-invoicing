package taxes

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kushaljain/go-invoicing/cache"
)

func setTax(cache *cache.RedisCache, key string, value *Taxes) {
	ctx := context.TODO()
	client := cache.GetClient()

	tax, err := GetSGSTList(value)
	if err != nil {
		panic(err)
	}
	json, err := json.Marshal(tax)
	if err != nil {
		panic(err)
	}

	client.Set(ctx, key, json, cache.Expires*time.Second)
}

func getTax(cache *cache.RedisCache, key string) *Taxes {
	ctx := context.TODO()
	client := cache.GetClient()

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	if val == string(redis.Nil) {
		return nil
	}
	var tax = NewTaxes()
	err = json.Unmarshal([]byte(val), &tax.SGSTList)
	if err != nil {
		return nil
	}
	return tax
}
