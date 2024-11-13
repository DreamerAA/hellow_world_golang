package main

import (
	"container/list"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

type FIFOCache struct {
	data     map[string]*list.Element
	list     *list.List
	capacity int
	mu       sync.RWMutex
}

func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
		panic("FIFOCache capacity must be greater than 0")
	}
	return &FIFOCache{
		data:     make(map[string]*list.Element),
		list:     list.New(),
		capacity: capacity,
	}
}

func (c *FIFOCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.data[key]; ok {
		log.Info("Cache hit")
		(*elem).Value.(*DataWrapper).data = value.(string)
		return
	}

	if c.list.Len() == c.capacity {
		log.Info("Cache overflow: first element removed")
		oldest_used_elem := c.list.Front()
		oldest_key := oldest_used_elem.Value.(*DataWrapper).key
		c.list.Remove(oldest_used_elem)
		delete(c.data, oldest_key)
	}
	log.Info("Element added to cache")
	new_elem := c.list.PushBack(&DataWrapper{key, value.(string)})
	c.data[key] = new_elem
}

func (c *FIFOCache) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	elem, ok := c.data[key]
	if !ok {
		log.Info("Cache miss")
		return nil
	}
	log.Info("Cache hit")
	return elem.Value.(*DataWrapper).data
}

func (c *FIFOCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.data[key]
	if !ok {
		log.Info("Cache miss")
		return
	}
	log.Info("Element removed from cache")
	c.list.Remove(elem)
	delete(c.data, key)
}

func (c *FIFOCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*list.Element)
	c.list = list.New()
}

func (c *FIFOCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.list.Len() > c.capacity || c.list.Len() != len(c.data) {
		fmt.Println("Error in cache size")
	}
	return c.list.Len()
}

func (c *FIFOCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, len(c.data))
	index := 0
	for key := range c.data {
		keys[index] = key
		index++
	}
	return keys
}
