package main

import "fmt"

func main() {
	cache := NewLFUCache(2)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.Set("key2", "value2")
	cache.Set("key2", "value2")
	cache.Set("key2", "value2")

	cache.Set("key1", "value1")
	cache.Set("key3", "value3")

	val1 := cache.Get("key1")
	val2 := cache.Get("key2")
	val3 := cache.Get("key3")

	if val1 != nil || val2 != "value2" || val3 != "value3" {
		fmt.Errorf("Expected nil, got %v", val2)
	}
}
