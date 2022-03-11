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

func (c *Cache) Close() {
	c.ticker.Stop()
	c.tickerStop <- true
}

func (c *Cache) Get(key string) (entry interface{}, found bool) {
	c.RLock()
	defer c.RUnlock()
	e, ok := c.m[key]
	if !ok {
		return nil, false
	}
	return e.data, true
}

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

func (c *Cache) Delete(key string) {
	delete(c.m, key)
}

func (c *Cache) cleanup() {
	c.Lock()
	for key, item := range c.m {
		if item.expired() {
			itemSize := uint64(size.Of(item))
			delete(c.m, key)
			c.mapSize -= itemSize
		}
	}
	c.Unlock()
}

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
