package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Cacher interface {
	Update(car Car) error
	Read(id int) (Car, error)
	Purge()
	Print()
}

type CachedCar struct {
	data      Car
	createdAt int64
	lastUsed  int64
	freq      int64
}

func (ce CachedCar) Sprintf() string {
	return fmt.Sprintf("Data: %v, createdAt: %v, lastUsed: %v, freq:%d\n", ce.data, ce.createdAt, ce.createdAt, ce.freq)
}

// O(n)
type InMemoryCache struct {
	data         map[int]CachedCar
	evictionAlgo evictionAlgo
	maxCapacity  int
	mu           sync.RWMutex
}

func NewInMemoryCache(ev evictionAlgo, maxCap int, initVals []Car) Cacher {
	cache := &InMemoryCache{
		data:         make(map[int]CachedCar),
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

func (c *InMemoryCache) Update(car Car) error {
	now := time.Now().UTC()

	if c == nil {
		return errors.New("car entity is incomplete")
	}
	c.mu.Lock()
	// log.Printf("Adding/modifying with: %v\n", car)
	if _, exists := c.data[car.Id()]; c.maxCapacity <= len(c.data) && !exists {
		// log.Printf("Starting cleanup...\n")
		c.evict()
	}
	c.data[car.Id()] = CachedCar{
		data:      car,
		createdAt: now.UnixMicro(),
		lastUsed:  now.UnixMicro(),
	}
	c.mu.Unlock()

	return nil
}

func (c *InMemoryCache) Read(id int) (Car, error) {
	c.mu.Lock()
	e, exists := c.data[id]
	if !exists {
		c.mu.Unlock()
		return Car{}, errors.New("object not in cache")
	}
	e.lastUsed = time.Now().UnixMicro()
	e.freq += 1
	c.mu.Unlock()

	return e.data, nil
}

func (c *InMemoryCache) Purge() {
	c.mu.Lock()
	c.data = make(map[int]CachedCar)
	c.mu.Unlock()
}

func (c *InMemoryCache) warm(ents []Car) {
	now := time.Now().UTC()

	for _, ent := range ents {
		c.data[ent.Id()] = CachedCar{
			data:      ent,
			createdAt: now.UnixMicro(),
			lastUsed:  now.UnixMicro(),
		}
	}
}

func (c *InMemoryCache) Print() {
	c.mu.RLock()
	for _, val := range c.data {
		fmt.Printf("%v", val.Sprintf())
	}
	c.mu.RUnlock()
}

func (c *InMemoryCache) delete(id int) {
	delete(c.data, id)
}

func (c *InMemoryCache) evict() {
	c.evictionAlgo.evict(c)
}
