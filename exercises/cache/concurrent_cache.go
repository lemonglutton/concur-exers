package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Cacher interface {
	Update(car Entity) error
	Read(id string) (Entity, error)
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
	return fmt.Sprintf("Data: %v, createdAt: %v, lastUsed: %v, freq:%d\n", ce.data, time.Unix(ce.createdAt, 0), time.Unix(ce.createdAt, 0), ce.freq)
}

// O(n)
type InMemoryCache struct {
	data         map[string]CachedEntity
	evictionAlgo evictionAlgo
	maxCapacity  int
	mu           sync.Mutex
}

func NewInMemoryCache(ev evictionAlgo, maxCap int, initVals []Entity) Cacher {
	cache := &InMemoryCache{
		data:         make(map[string]CachedEntity),
		evictionAlgo: ev,
		maxCapacity:  maxCap,
		mu:           sync.Mutex{},
	}

	if len(initVals) > cache.maxCapacity {
		panic("exceeding data capacity during warmup")
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

	c.mu.Lock()
	log.Printf("Adding/modifying with: %v\n", ent)
	if _, exists := c.data[ent.Id()]; c.maxCapacity <= len(c.data) && !exists {
		log.Printf("Starting cleanup...")
		c.evict()
	}
	c.data[ent.Id()] = CachedEntity{
		data:      ent,
		createdAt: now.Unix(),
		lastUsed:  now.Unix(),
	}
	c.mu.Unlock()

	return nil
}

func (c *InMemoryCache) Read(id string) (Entity, error) {
	c.mu.Lock()
	e, exists := c.data[id]
	if !exists {
		c.mu.Unlock()
		return nil, errors.New("object not in cache")
	}
	e.lastUsed = time.Now().Unix()
	e.freq += 1
	c.mu.Unlock()

	return e.data, nil
}

func (c *InMemoryCache) Purge() {
	c.mu.Lock()
	c.data = make(map[string]CachedEntity)
	c.mu.Unlock()
}

func (c *InMemoryCache) warm(ents []Entity) {
	now := time.Now().UTC()

	for _, ent := range ents {
		c.data[ent.Id()] = CachedEntity{
			data:      ent,
			createdAt: now.Unix(),
			lastUsed:  now.Unix(),
		}
	}
}

func (c *InMemoryCache) Print() {
	for _, val := range c.data {
		fmt.Printf("%v", val.Sprintf())
	}
}

func (c *InMemoryCache) delete(id string) {
	delete(c.data, id)
}

func (c *InMemoryCache) evict() {
	c.evictionAlgo.evict(c)
}
