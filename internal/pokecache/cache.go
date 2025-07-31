package pokecache

import (
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	mu sync.Mutex
	location map[string]cacheEntry
}

func NewCache(customDuration time.Duration) *Cache {
	c := &Cache{
		location: make(map[string]cacheEntry),
	}

	go c.reapLoop(customDuration)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.location[key] = cacheEntry{
		val: val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	entry, found := c.location[key]

	if !found {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	for {
		time.Sleep(interval)
		c.mu.Lock()
		for k, v := range c.location {
			if time.Since(v.createdAt) > interval {
				delete(c.location, k)
			}
		}
		c.mu.Unlock()
	}
}