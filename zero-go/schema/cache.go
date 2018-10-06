package schema

import (
	"context"

	"github.com/graph-gophers/dataloader"
	lru "github.com/hashicorp/golang-lru"
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

var Loader = dataloader.NewBatchedLoader(batchedFunc, dataloader.WithCache(Cache))
