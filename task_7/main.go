package main

import (
	"container/list"
	"fmt"
	"strconv"
)

type ICache interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type SimpleCache struct {
	data map[string]interface{}
}

func (c *SimpleCache) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *SimpleCache) Get(key string) (interface{}, bool) {
	value, ok := c.data[key]
	return value, ok
}

func (c *SimpleCache) Delete(key string) {
	delete(c.data, key)
}

type DataWrapper struct {
	key  string
	data string
}

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

func testPolymorphism(cache ICache) {
	fmt.Println(cache.Get("key1"))
}

func repeatGet(cache ICache) {
	for i := 0; i < 10; i++ {
		for i := 2; i < 4; i++ {
			key := "key" + strconv.Itoa(i)
			_, _ = cache.Get(key)
		}
	}
}

func setValues(cache ICache) {
	for i := 1; i < 4; i++ {
		key := "key" + strconv.Itoa(i)
		val := "value" + strconv.Itoa(i)
		cache.Set(key, val)
	}
}

func main() {
	simple_cache := &SimpleCache{data: make(map[string]interface{})}
	lru_cache := NewLRUCache(3)

	setValues(simple_cache)
	setValues(lru_cache)

	repeatGet(simple_cache)

	fmt.Println(simple_cache.Get("key1"))
	simple_cache.Set("key4", "value4")
	fmt.Println(simple_cache.Get("key1"))
	simple_cache.Delete("key1")
	fmt.Println(simple_cache.Get("key1"))

	repeatGet(lru_cache)

	fmt.Println(len(lru_cache.data), lru_cache.list.Len())
	lru_cache.Set("key4", "value4")
	lru_cache.Set("key5", "value5")
	lru_cache.Set("key6", "value6")
	fmt.Println(len(lru_cache.data), lru_cache.list.Len())
	lru_cache.Delete("key1")
	lru_cache.Delete("key2")
	fmt.Println(len(lru_cache.data), lru_cache.list.Len())

	testPolymorphism(simple_cache)
	testPolymorphism(lru_cache)
}
