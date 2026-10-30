package pokecache

import (
	"time"
	"fmt"
	"sync"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	entries map[string]cacheEntry
	interval time.Duration
	mu sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	entries := make(map[string]cacheEntry)
	cache := &Cache{
		entries: entries,
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	cache.entries[key] = entry
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry, ok := cache.entries[key]
	return entry.val, ok
}

func (cache *Cache) reapLoop() {
	ticker := time.NewTicker(cache.interval)
	defer ticker.Stop()
	for range ticker.C {
		cache.mu.Lock()
		for key, entry := range cache.entries {
			if time.Since(entry.createdAt) > cache.interval {
				delete(cache.entries, key)
			}
		}
		cache.mu.Unlock()
	}
}

func (cache *Cache) Show() {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	fmt.Printf("Showing cache:\n")
	for key, val := range cache.entries {
		fmt.Printf("Key: %v, Value: %v\n", key, val)
	}
	return
}

func main() {
	cache := NewCache(1 * time.Second)
	cache.Add("testKey", []byte("testValue"))
	value, _ := cache.Get("testKey")
	fmt.Println(string(value))
	for {
		time.Sleep(1 * time.Second)
	}
}