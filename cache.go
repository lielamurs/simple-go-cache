package sgcache

import (
	"sync"
	"time"

	"github.com/DmitriyVTitov/size"
)

type MemoryCache interface {
	Get(key string) (entry interface{}, found bool)
	Set(key string, data interface{}, ttl time.Duration)
	Delete(key string)
	Close()
}

type Cache struct {
	sync.RWMutex
	m                  map[string]*Item
	mapSize, sizeLimit uint64
	ticker             *time.Ticker
	tickerStop         chan bool
}

// Return a new cache with a cleanup interval and size limit in bytes.
// Automatic cleanup is started to delete expired items.
func New(cleanupInterval time.Duration, sizeLimit uint64) *Cache {
	interval := cleanupInterval
	if interval < time.Second {
		interval = time.Second
	}
	cache := &Cache{
		m:          make(map[string]*Item),
		sizeLimit:  sizeLimit,
		ticker:     time.NewTicker(interval),
		tickerStop: make(chan bool),
	}
	cache.mapSize = uint64(size.Of(cache.m))
	cache.startCleanupTimer()
	return cache
}

// Stop the cleanup process
func (c *Cache) Close() {
	c.ticker.Stop()
	c.tickerStop <- true
}

// Get returns a cache entry for the provided key.
func (c *Cache) Get(key string) (entry interface{}, found bool) {
	c.RLock()
	defer c.RUnlock()
	e, ok := c.m[key]
	if !ok {
		return nil, false
	}
	return e.data, true
}

// Set adds new item to cache with a ttl duration. If adding the item
// causes cache to exceed defined memory limit no action is performed.
func (c *Cache) Set(key string, data interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	item := Item{data: data, ttl: time.Now().Add(ttl)}
	itemSize := uint64(size.Of(item))
	if itemSize+c.mapSize <= c.sizeLimit {
		c.m[key] = &item
		c.mapSize += itemSize
	}
}

// Delete removes item from cache.
func (c *Cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	c.removeItemSize(c.m[key])
	delete(c.m, key)
}

// Cleanup cecks cached items and removes expired ones.
func (c *Cache) cleanup() {
	c.Lock()
	defer c.Unlock()
	for key, item := range c.m {
		if item.expired() {
			c.removeItemSize(c.m[key])
			delete(c.m, key)
		}
	}
}

// Start the cleanup timer
func (c *Cache) startCleanupTimer() {
	go (func() {
		for {
			select {
			case <-c.tickerStop:
				return
			case <-c.ticker.C:
				c.cleanup()
			}
		}
	})()
}

// Removes item size from map size
func (c *Cache) removeItemSize(item *Item) {
	itemSize := uint64(size.Of(item))
	c.mapSize -= itemSize
}
