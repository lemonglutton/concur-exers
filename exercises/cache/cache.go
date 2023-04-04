package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Cacher interface {
	Set(e Entity) error
	Read(id int) (Entity, error)
	Purge()
	Print()
}

type CachedEntity struct {
	data      Entity
	createdAt int64
	lastUsed  int64
	freq      int64
}

func (ce CachedEntity) Sprintf() string {
	return fmt.Sprintf("Data: %v, createdAt: %v, lastUsed: %v, freq:%d\n", ce.data, ce.createdAt, ce.createdAt, ce.freq)
}

// O(n)
type InMemoryCache struct {
	data         map[int]CachedEntity
	evictionAlgo evictionAlgo
	maxCapacity  int
	mu           sync.RWMutex
}

func NewInMemoryCache(ev evictionAlgo, maxCap int, initVals []Entity) Cacher {
	cache := &InMemoryCache{
		data:         make(map[int]CachedEntity),
		evictionAlgo: ev,
		maxCapacity:  maxCap,
		mu:           sync.RWMutex{},
	}

	if len(initVals) > cache.maxCapacity {
		panic("exceeding data capacity during warmup")
	}

	if len(initVals) > 0 {
		cache.warm(initVals)
	}

	return cache
}

func (c *InMemoryCache) Set(e Entity) error {
	now := time.Now().UTC()

	if c == nil {
		return errors.New("car entity is incomplete")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	log.Printf("Adding/modifying with: %v\n", e)
	if _, exists := c.data[e.Id()]; c.maxCapacity <= len(c.data) && !exists {
		// log.Printf("Starting cleanup...\n")
		c.evict()
	}
	c.data[e.Id()] = CachedEntity{
		data:      e,
		createdAt: now.UnixMicro(),
		lastUsed:  now.UnixMicro(),
	}

	return nil
}

func (c *InMemoryCache) Read(id int) (Entity, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, exists := c.data[id]
	if !exists {
		return Car{}, errors.New("object not in cache")
	}

	e.lastUsed = time.Now().UnixMicro()
	e.freq += 1
	return e.data, nil
}

func (c *InMemoryCache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[int]CachedEntity)
}

func (c *InMemoryCache) warm(ents []Entity) {
	now := time.Now().UTC()

	for _, ent := range ents {
		c.data[ent.Id()] = CachedEntity{
			data:      ent,
			createdAt: now.UnixMicro(),
			lastUsed:  now.UnixMicro(),
		}
	}
}

func (c *InMemoryCache) Print() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, val := range c.data {
		fmt.Printf("%v", val.Sprintf())
	}

}

func (c *InMemoryCache) delete(id int) {
	delete(c.data, id)
}

func (c *InMemoryCache) evict() {
	c.evictionAlgo.evict(c)
}
