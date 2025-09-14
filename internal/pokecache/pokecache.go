package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheItem struct {
	CreatedAt time.Time
	Data      []byte
}

type cache struct {
	items      map[string]cacheItem
	lock       sync.RWMutex
	expiration time.Duration
}

func (c *cache) Add(key string, value []byte) {
	c.lock.Lock()

	c.items[key] = cacheItem{
		CreatedAt: time.Now(),
		Data:      value,
	}

	c.lock.Unlock()
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	result, exists := c.items[key]
	if !exists {
		fmt.Println("Cache miss")
		return []byte{}, false
	}

	fmt.Println("Cache hit")
	return result.Data, true
}

func (c *cache) reapLoop() {
	for {
		time.Sleep(c.expiration)

		now := time.Now()
		for key, value := range c.items {
			if now.Sub(value.CreatedAt) > c.expiration {
				c.lock.Lock()
				delete(c.items, key)
				c.lock.Unlock()
			}
		}
	}
}

var Cache = cache{
	items:      map[string]cacheItem{},
	lock:       sync.RWMutex{},
	expiration: 30 * time.Second,
}

func init() {
	go Cache.reapLoop()
}
