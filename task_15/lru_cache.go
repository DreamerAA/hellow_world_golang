package main

import (
	"container/list"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

type LRUCache struct {
	data     map[string]*list.Element
	list     *list.List
	capacity int
	mu       sync.RWMutex
}

func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		panic("LRUCache capacity must be greater than 0")
	}
	return &LRUCache{
		data:     make(map[string]*list.Element),
		list:     list.New(),
		capacity: capacity,
	}
}

func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.data[key]; ok {
		log.Info("Cache hit")
		c.list.MoveToFront(elem)
		(*elem).Value.(*DataWrapper).data = value.(string)
		return
	}

	if c.list.Len() == c.capacity {
		log.Info("Cache overflow: first element removed")
		least_used_elem := c.list.Back()
		least_used_key := least_used_elem.Value.(*DataWrapper).key
		c.list.Remove(least_used_elem)
		delete(c.data, least_used_key)
	}
	log.Info("Element added to cache")
	new_elem := c.list.PushBack(&DataWrapper{key, value.(string)})
	c.data[key] = new_elem
}

func (c *LRUCache) Get(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.data[key]
	if !ok {
		log.Info("Cache miss")
		return nil
	}
	log.Info("Cache hit")
	c.list.MoveToFront(elem)
	return elem.Value.(*DataWrapper).data
}

func (c *LRUCache) Delete(key string) {
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

func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*list.Element)
	c.list = list.New()
}

func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.list.Len() > c.capacity || c.list.Len() != len(c.data) {
		fmt.Println("Error in cache size")
	}
	return c.list.Len()
}

func (c *LRUCache) Keys() []string {
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
