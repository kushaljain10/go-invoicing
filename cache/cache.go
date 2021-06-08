package cache

import (
	"github.com/kushaljain/go-invoicing/taxes"
)

type Cache interface {
	Set(key string, value *taxes.Taxes)
	Get(key string) *taxes.Taxes
}
