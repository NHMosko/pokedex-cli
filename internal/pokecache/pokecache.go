package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	Data map[string]cacheEntry
	Mut *sync.RWMutex
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		cutoff := time.Now().Add(-interval)
		c.Mut.Lock()
		for key, entry := range c.Data {
			if entry.createdAt.Before(cutoff) {
				delete(c.Data, key)
			}
		}
		c.Mut.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.Mut.Lock()
	c.Data[key] = cacheEntry{createdAt: time.Now(), val: val}
	c.Mut.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mut.RLock()
	out, ok:= c.Data[key]
	c.Mut.RUnlock()
	if ok {
		return out.val, true
	}
	return nil, false
}



func NewCache(interval time.Duration) *Cache {
	cache := &Cache{map[string]cacheEntry{}, &sync.RWMutex{}} 
	go cache.reapLoop(interval)
	return cache
}
