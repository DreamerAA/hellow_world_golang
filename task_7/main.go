package main

import (
	"fmt"
	"math"
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

type DataWrapper struct {
	data    string
	counter int
}

type LRUCache struct {
	SimpleCache
	count int
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

func (c *LRUCache) getMinimumAndDelete() {
	leastUsedKey := ""
	leastUsedCounter := math.MaxInt
	for key, value := range c.data {
		data := value.(*DataWrapper)
		if data.counter < leastUsedCounter {
			leastUsedCounter = data.counter
			leastUsedKey = key
		}
	}
	delete(c.data, leastUsedKey)
}

func (c *LRUCache) Set(key string, value interface{}) {
	if len(c.data) >= c.count {
		c.getMinimumAndDelete()
	}
	c.data[key] = &DataWrapper{value.(string), 0}
}
func (c *LRUCache) Get(key string) (interface{}, bool) {
	value, ok := c.data[key]
	if !ok {
		return nil, ok
	}

	wrapper := (value).(*DataWrapper)
	wrapper.counter++
	return wrapper.data, ok
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
	simple_cache := SimpleCache{data: make(map[string]interface{})}
	lru_cache := LRUCache{SimpleCache{data: make(map[string]interface{})}, 3}

	setValues(&simple_cache)
	setValues(&lru_cache)

	repeatGet(&simple_cache)

	fmt.Println(simple_cache.Get("key1"))
	simple_cache.Set("key4", "value4")
	fmt.Println(simple_cache.Get("key1"))
	simple_cache.Delete("key1")
	fmt.Println(simple_cache.Get("key1"))

	repeatGet(&lru_cache)

	fmt.Println(lru_cache.Get("key1"))
	lru_cache.Set("key4", "value4")
	fmt.Println(lru_cache.Get("key1"))
	lru_cache.Delete("key1")
	fmt.Println(lru_cache.Get("key1"))

	testPolymorphism(&simple_cache)
	testPolymorphism(&lru_cache)
}
