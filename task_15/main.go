package main

import (
	"fmt"
)

func main() {
	// Пример использования LRUCache
	cache := NewLRUCache(3)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Get("key1")           // "key1" становится последним использованным
	cache.Set("key4", "value4") // "key2" вытесняется

	fmt.Println(cache.Get("key2")) // nil, false
	fmt.Println(cache.Get("key3")) // "value3", true
}
