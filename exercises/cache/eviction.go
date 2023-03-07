package main

import (
	"log"
	"sort"
)

type evictionAlgo interface {
	evict(c *InMemoryCache)
}

type fifo struct{}

func (f *fifo) evict(c *InMemoryCache) {
	keys := []string{}
	for key, _ := range c.data {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return c.data[keys[i]].createdAt < c.data[keys[j]].createdAt
	})

	log.Printf("Removing: %v\n", c.data[keys[0]].Sprintf())
	c.delete(keys[0])
}

type lru struct{}

func (l *lru) evict(c *InMemoryCache) {
	keys := []string{}
	for key, _ := range c.data {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return c.data[keys[i]].lastUsed > c.data[keys[j]].createdAt
	})

	c.delete(keys[0])
}

type lfu struct{}

func (l *lfu) evict(c *InMemoryCache) {
	keys := []string{}
	for key, _ := range c.data {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return c.data[keys[i]].freq < c.data[keys[j]].freq
	})

	c.delete(keys[0])
}
