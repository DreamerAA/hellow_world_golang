package main

import (
	"container/list"
	"fmt"
)

type LRUCache struct {
	data     map[string]*list.Element
	list     *list.List
	capacity int
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
	if elem, ok := c.data[key]; ok {
		c.list.MoveToFront(elem)
		(*elem).Value.(*DataWrapper).data = value.(string)
		return
	}

	if c.list.Len() == c.capacity {
		least_used_elem := c.list.Back()
		least_used_key := least_used_elem.Value.(*DataWrapper).key
		c.list.Remove(least_used_elem)
		delete(c.data, least_used_key)
	}
	new_elem := c.list.PushBack(&DataWrapper{key, value.(string)})
	c.data[key] = new_elem
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	elem, ok := c.data[key]
	if !ok {
		return nil, ok
	}
	c.list.MoveToFront(elem)
	return elem.Value.(*DataWrapper).data, ok
}

func (c *LRUCache) Delete(key string) {
	elem, ok := c.data[key]
	if !ok {
		return
	}
	c.list.Remove(elem)
	delete(c.data, key)
}

func (c *LRUCache) Clear() {
	c.data = make(map[string]*list.Element)
}

func (c *LRUCache) Size() int {
	if c.list.Len() > c.capacity || c.list.Len() != len(c.data) {
		fmt.Println("Error in cache size")
	}
	return c.list.Len()
}

func (c *LRUCache) Keys() []string {
	keys := make([]string, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}
