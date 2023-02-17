package main

import (
	"errors"
	"sync"
	"time"
)

type Cacher interface {
	Update(car Entity) error
	Read(id string) (Entity, error)
	Purge()
}

type cachedEntity struct {
	data      Entity
	createdAt int64
	lastUsed  int64
	freq      int64
}

// O(n)
type InMemoryCache struct {
	data         map[string]cachedEntity
	evictionAlgo evictionAlgo
	maxCapacity  int
	cond         sync.Cond
}

func NewInMemoryCache(ev evictionAlgo, maxCap int, initVals []Entity) Cacher {
	cache := &InMemoryCache{
		data:         make(map[string]cachedEntity),
		evictionAlgo: ev,
		maxCapacity:  maxCap,
	}

	if len(initVals) > 0 {
		cache.warm(initVals)
	}

	return cache
}

func (c *InMemoryCache) Update(ent Entity) error {
	now := time.Now().UTC()

	if ent == nil {
		return errors.New("Entity is incomplete")
	}

	c.cond.L.Lock()
	if c.maxCapacity == len(c.data) {
		c.evict()
		c.cond.Wait()
	}
	c.data[ent.Id()] = cachedEntity{
		data:      ent,
		createdAt: now.Unix(),
		lastUsed:  now.Unix(),
	}
	c.cond.L.Unlock()

	return nil
}

func (c *InMemoryCache) Read(id string) (Entity, error) {
	c.cond.L.Lock()
	if c.maxCapacity == len(c.data) {
		c.cond.Wait()
	}

	e, exists := c.data[id]
	if !exists {
		return nil, errors.New("object not in cache")
	}

	e.lastUsed = time.Now().Unix()
	e.freq += 1
	c.cond.L.Unlock()

	return e.data, nil
}

func (c *InMemoryCache) warm(ents []Entity) {
	now := time.Now().UTC()

	for _, ent := range ents {
		c.data[ent.Id()] = cachedEntity{
			data:      ent,
			createdAt: now.Unix(),
			lastUsed:  now.Unix(),
		}
	}
}

func (c *InMemoryCache) delete(id string) {
	c.cond.L.Lock()

	if c.maxCapacity == len(c.data) {
		c.cond.Wait()
	}
	delete(c.data, id)
	c.cond.L.Unlock()
}

func (c *InMemoryCache) evict() {
	c.evictionAlgo.evict(c)
	c.cond.Broadcast()
}

func (c *InMemoryCache) Purge() {
	c.cond.L.Lock()

	if c.maxCapacity == len(c.data) {
		c.cond.Wait()
	}
	c.data = make(map[string]cachedEntity)
	c.cond.L.Unlock()
}
