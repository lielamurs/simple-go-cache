package sgcache

import (
	"sync"
	"time"
)

type MemoryCache interface {
	Get(key string) (entry interface{}, found bool)
	Set(key string, data interface{}, ttl time.Duration)
	Delete(key string)
}

type Cache struct {
	sync.RWMutex
	ct time.Duration
	m  map[string]*Item
}

func New(cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		ct: cleanupInterval,
		m:  make(map[string]*Item),
	}
	cache.startCleanupTimer()
	return cache
}

func (c *Cache) Get(key string) (entry interface{}, found bool) {
	c.RLock()
	defer c.RUnlock()
	if _, ok := c.m[key]; !ok {
		return nil, false
	}
	e := c.m[key]
	return e.data, true
}

func (c *Cache) Set(key string, data interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	c.m[key] = &Item{data: data, ttl: time.Now().Add(ttl)}
}

func (c *Cache) Delete(key string) {
	delete(c.m, key)
}

func (c *Cache) cleanup() {
	c.Lock()
	for key, item := range c.m {
		if item.expired() {
			delete(c.m, key)
		}
	}
	c.Unlock()
}

func (c *Cache) startCleanupTimer() {
	interval := c.ct
	if interval < time.Second {
		interval = time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go (func() {
		for range ticker.C {
			c.cleanup()
		}
	})()
}
