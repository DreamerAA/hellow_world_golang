package main

import (
	"fmt"
	"sort"
	"sync"

	log "github.com/sirupsen/logrus"
)

type KeyCounter struct {
	key     string
	counter int
}

type CounterContainer struct {
	data []KeyCounter
}

// Add вставляет элемент в контейнер, сохраняя его отсортированным.
func (sc *CounterContainer) Add(key string, counter int) {
	index := sort.Search(len(sc.data), func(i int) bool {
		return sc.data[i].counter >= counter
	})
	sc.data = append(sc.data[:index], append([]KeyCounter{KeyCounter{key, counter}}, sc.data[index:]...)...)
}

func (sc *CounterContainer) Increment(key string) bool {
	index := sc.findIndex(key)
	if index == -1 {
		return false
	}
	sc.data[index].counter++
	for index < len(sc.data)-1 && sc.data[index].counter > sc.data[index+1].counter {
		sc.data[index], sc.data[index+1] = sc.data[index+1], sc.data[index]
		index++
	}
	return true
}

func (sc *CounterContainer) findIndex(key string) int {
	for i, item := range sc.data {
		if item.key == key {
			return i // Возвращаем индекс, если нашли
		}
	}
	return -1 // Возвращаем -1, если не нашли
}

// Remove удаляет элемент из контейнера.
func (sc *CounterContainer) Remove(key string) bool {
	index := sc.findIndex(key)
	if index == -1 {
		return false
	}
	sc.data = append(sc.data[:index], sc.data[index+1:]...)
	return true

}

type LFUCache struct {
	data     map[string]DataWrapper
	counters CounterContainer
	capacity int
	mu       sync.RWMutex
}

func NewCounterContainer() CounterContainer {
	return CounterContainer{data: []KeyCounter{}}
}

func NewLFUCache(capacity int) *LFUCache {
	if capacity <= 0 {
		panic("LFUCache capacity must be greater than 0")
	}
	return &LFUCache{
		data:     make(map[string]DataWrapper),
		counters: NewCounterContainer(),
		capacity: capacity,
	}
}

func (c *LFUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.data[key]; ok {
		log.Info("Cache hit")
		c.counters.Increment(key)
		elem.data = value
		return
	}

	if len(c.data) == c.capacity {
		log.Info("Cache overflow: least used element removed")
		key_least_used_elem := c.counters.data[0].key
		c.counters.Remove(key_least_used_elem)
		delete(c.data, key_least_used_elem)
	}
	log.Info("Element added to cache")
	c.counters.Add(key, 1)
	c.data[key] = DataWrapper{key, value}
}

func (c *LFUCache) Get(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.data[key]
	c.counters.Increment(key)
	if !ok {
		log.Info("Cache miss")
		return nil
	}
	log.Info("Cache hit")
	return elem.data
}

func (c *LFUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	c.counters.Remove(key)
}

func (c *LFUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]DataWrapper)
	c.counters = NewCounterContainer()
}

func (c *LFUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s := len(c.counters.data)
	if s > c.capacity || s != len(c.data) {
		fmt.Println("Error in cache size")
	}
	return s
}

func (c *LFUCache) Keys() []string {
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
