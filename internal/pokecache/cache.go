package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu           *sync.RWMutex
	mp           map[string]cacheEntry
	TimeInterval time.Duration
	stopCh       chan struct{}
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	cache := Cache{
		mu:           &sync.RWMutex{},
		mp:           make(map[string]cacheEntry),
		TimeInterval: interval,
		stopCh:       make(chan struct{}),
	}
	ticker := time.NewTicker(cache.TimeInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.reapLoop()
			case <-cache.stopCh:
				return
			}
		}
	}()
	return cache
}

func (c *Cache) Close() {
	select {
	case <-c.stopCh:
		// already closed
	default:
		close(c.stopCh)
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mp[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	v, ok := c.mp[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}

	if time.Since(v.createdAt) > c.TimeInterval {
		c.mu.Lock()
		// re-check under write lock to avoid races
		if cur, ok := c.mp[key]; ok && time.Since(cur.createdAt) > c.TimeInterval {
			delete(c.mp, key)
		}
		c.mu.Unlock()
		return nil, false
	}

	return v.val, true
}

func (c *Cache) reapLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.mp {
		if time.Since(v.createdAt) > c.TimeInterval {
			delete(c.mp, k)
		}
	}
}
