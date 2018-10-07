package schema

import (
	"context"

	lru "github.com/hashicorp/golang-lru"
	dataloader "gopkg.in/nicksrandall/dataloader.v5"
)

// Cache implements the dataloader.Cache interface
type Cache struct {
	*lru.ARCCache
}

// Get gets an item from the cache
func (c *Cache) Get(_ context.Context, key dataloader.Key) (dataloader.Thunk, bool) {
	v, ok := c.ARCCache.Get(key)
	if ok {
		return v.(dataloader.Thunk), ok
	}
	return nil, ok
}

// Set sets an item in the cache
func (c *Cache) Set(_ context.Context, key dataloader.Key, value dataloader.Thunk) {
	c.ARCCache.Add(key, value)
}

// Delete deletes an item in the cache
func (c *Cache) Delete(_ context.Context, key dataloader.Key) bool {
	if c.ARCCache.Contains(key) {
		c.ARCCache.Remove(key)
		return true
	}
	return false
}

// Clear cleasrs the cache
func (c *Cache) Clear() {
	c.ARCCache.Purge()
}

var c, _ = lru.NewARC(100)
var cache = &Cache{c}

// Loader is a cached loader
var Loader = dataloader.NewBatchedLoader(batchFunc, dataloader.WithCache(cache))

func batchFunc(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
	var results []*dataloader.Result
	// do some pretend work to resolve keys
	for _, key := range keys {
		results = append(results, &dataloader.Result{key.String(), nil})
	}
	return results
}
