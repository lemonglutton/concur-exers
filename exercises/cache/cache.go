package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Cacher interface {
	Update(car Car) error
	Read(id string) (Car, error)
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
	return fmt.Sprintf("Data: %v, createdAt: %v, lastUsed: %v, freq:%d\n", ce.data, time.Unix(ce.createdAt, 0), time.Unix(ce.createdAt, 0), ce.freq)
}

// O(n)
type InMemoryCache struct {
	data         map[string]CachedCar
	evictionAlgo evictionAlgo
	maxCapacity  int
	mu           sync.Mutex
}

func NewInMemoryCache(ev evictionAlgo, maxCap int, initVals []Car) Cacher {
	cache := &InMemoryCache{
		data:         make(map[string]CachedCar),
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

func (c *InMemoryCache) Update(car Car) error {
	now := time.Now().UTC()

	if c == nil {
		return errors.New("Entity is incomplete")
	}

	c.mu.Lock()
	log.Printf("Adding/modifying with: %v\n", c)
	if _, exists := c.data[car.Id()]; c.maxCapacity <= len(c.data) && !exists {
		log.Printf("Starting cleanup...")
		c.evict()
	}
	c.data[car.Id()] = CachedCar{
		data:      car,
		createdAt: now.Unix(),
		lastUsed:  now.Unix(),
	}
	c.mu.Unlock()

	return nil
}

func (c *InMemoryCache) Read(id string) (Car, error) {
	c.mu.Lock()
	e, exists := c.data[id]
	if !exists {
		c.mu.Unlock()
		return Car{}, errors.New("object not in cache")
	}
	e.lastUsed = time.Now().Unix()
	e.freq += 1
	c.mu.Unlock()

	return e.data, nil
}

func (c *InMemoryCache) Purge() {
	c.mu.Lock()
	c.data = make(map[string]CachedCar)
	c.mu.Unlock()
}

func (c *InMemoryCache) warm(ents []Car) {
	now := time.Now().UTC()

	for _, ent := range ents {
		c.data[ent.Id()] = CachedCar{
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
