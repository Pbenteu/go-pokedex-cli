package pokecache

import (
	"sync"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	Cache.Add("teste", []byte("teste"))

	value, exists := Cache.Get("teste")
	if !exists {
		t.Errorf("Key doesnt exist")
		return
	}

	if string(value) != "teste" {
		t.Errorf("expected: 'teste', received: '%s'", string(value))
	}
}

func TestCacheReapLoop(t *testing.T) {
	fastCache := cache{
		expiration: 5 * time.Millisecond,
		items:      map[string]cacheItem{},
		lock:       sync.RWMutex{},
	}
	go fastCache.reapLoop()

	fastCache.Add("teste", []byte("teste"))

	time.Sleep(10 * time.Millisecond)

	_, exists := fastCache.Get("teste")
	if exists {
		t.Errorf("Key should have been expired")
		return
	}
}
