package cacheinfra

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key        string
	value      interface{}
	expiration int64
}

type InMemoryCache struct {
	mu         sync.Mutex
	maxEntries int
	ll         *list.List
	cache      map[string]*list.Element
}

func NewInMemoryCache(maxEntries int) *InMemoryCache {
	return &InMemoryCache{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

// Set inserts or updates a value with TTL
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		ee.Value.(*entry).expiration = expirationTime(ttl)
		return
	}

	e := &entry{
		key:        key,
		value:      value,
		expiration: expirationTime(ttl),
	}
	ele := c.ll.PushFront(e)
	c.cache[key] = ele

	if c.maxEntries > 0 && c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
}

// Get returns a value if present and not expired
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ele, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	ent := ele.Value.(*entry)

	if ent.expiration > 0 && time.Now().UnixNano() > ent.expiration {
		c.removeElement(ele)
		return nil, false
	}

	c.ll.MoveToFront(ele)
	return ent.value, true
}

// removeOldest evicts the LRU entry
func (c *InMemoryCache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *InMemoryCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	ent := e.Value.(*entry)
	delete(c.cache, ent.key)
}

func expirationTime(ttl time.Duration) int64 {
	if ttl > 0 {
		return time.Now().Add(ttl).UnixNano()
	}
	return 0
}
