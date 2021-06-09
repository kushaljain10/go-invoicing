package taxes

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kushaljain/go-invoicing/cache"
)

func setTaxInCache(cache *cache.RedisCache, key string, taxes *Taxes) {
	ctx := context.TODO()
	client := cache.GetClient()

	json, err := json.Marshal(taxes.SGSTList)
	if err != nil {
		panic(err)
	}

	client.Set(ctx, key, json, cache.Expires*time.Second)
}

func getTaxFromCache(cache *cache.RedisCache, key string) *Taxes {
	ctx := context.TODO()
	client := cache.GetClient()

	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return nil
	}

	var tax = NewTaxes()
	err = json.Unmarshal([]byte(val), &tax.SGSTList)
	if err != nil {
		return nil
	}
	return tax
}
