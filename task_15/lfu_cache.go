package main

import (
	"fmt"
	"sort"
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
	index := sort.Search(len(sc.data), func(i int) bool {
		return sc.data[i].key == key
	})

	if index < len(sc.data) {
		sc.data[index].counter++
		for index < len(sc.data)-1 && sc.data[index].counter > sc.data[index+1].counter {
			sc.data[index], sc.data[index+1] = sc.data[index+1], sc.data[index]
			index++
		}
		return true
	}
	return false
}

// Remove удаляет элемент из контейнера.
func (sc *CounterContainer) Remove(key string) bool {
	index := sort.Search(len(sc.data), func(i int) bool {
		return sc.data[i].key == key
	})

	if index < len(sc.data) {
		sc.data = append(sc.data[:index], sc.data[index+1:]...)
		return true
	}
	return false
}

type LFUCache struct {
	data     map[string]DataWrapper
	counters CounterContainer
	capacity int
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
	if elem, ok := c.data[key]; ok {
		c.counters.Increment(key)
		elem.data = value
		return
	}

	if len(c.data) == c.capacity {
		key_least_used_elem := c.counters.data[0].key
		c.counters.Remove(key_least_used_elem)
		delete(c.data, key_least_used_elem)
	}
	c.counters.Add(key, 1)
	c.data[key] = DataWrapper{key, value}
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	elem, ok := c.data[key]
	c.counters.Increment(key)
	if !ok {
		return nil, ok
	}
	return elem.data, ok
}

func (c *LFUCache) Delete(key string) {
	delete(c.data, key)
	c.counters.Remove(key)
}

func (c *LFUCache) Clear() {
	c.data = make(map[string]DataWrapper)
	c.counters = NewCounterContainer()
}

func (c *LFUCache) Size() int {
	s := len(c.counters.data)
	if s > c.capacity || s != len(c.data) {
		fmt.Println("Error in cache size")
	}
	return s
}

func (c *LFUCache) Keys() []string {
	keys := make([]string, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}
