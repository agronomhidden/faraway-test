package cache

import "sync"

type Cachier interface {
	Set(key string)
	Exist(key string) bool
	Delete(key string)
}

func NewCache() Cachier {
	return &requestsCache{
		m: make(map[string]struct{}),
	}
}

type requestsCache struct {
	mx sync.RWMutex
	m  map[string]struct{}
}

func (c *requestsCache) Set(key string) {
	c.mx.Lock()
	c.m[key] = struct{}{}
	c.mx.Unlock()
}
func (c *requestsCache) Exist(key string) bool {
	c.mx.RLock()
	_, ok := c.m[key]
	c.mx.RUnlock()
	return ok
}

func (c *requestsCache) Delete(key string) {
	c.mx.Lock()
	delete(c.m, key)
	c.mx.Unlock()
}
