package pokecache

import (
	"time"
)

func NewCache(interval time.Duration) *Cache {
	var newCache Cache
	newCache.cache = make(map[string]cacheEntry)
	go newCache.reapLoop(interval)
	return &newCache
}

func (c *Cache) Add(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var entry cacheEntry
	entry.createdAt = time.Now()
	entry.val = value
	c.cache[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.cache[key]

	return entry.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mutex.Lock()
		cutoff := time.Now().Add(-interval)
		for key, entry := range c.cache {
			if entry.createdAt.Before(cutoff) {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}
