package main

import (
	"fmt"
	"testing"
)

func keysTest(t *testing.T, cache ICache) {
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	keys := cache.Keys()
	if len(keys) != 2 || keys[0] != "key1" || keys[1] != "key2" {
		fmt.Printf("Length: %v", len(keys))
		fmt.Printf("key1: [%v]", keys[0])
		fmt.Printf("key2: [%v]", keys[1])
		t.Errorf("Expected [key1, key2], got %v", keys)
	}
}

func TestLRUKeys(t *testing.T) {
	cache := NewLRUCache(3)
	keysTest(t, cache)
}
func TestLFUKeys(t *testing.T) {
	cache := NewLFUCache(3)
	keysTest(t, cache)
}
func TestFIFOKeys(t *testing.T) {
	cache := NewFIFOCache(3)
	keysTest(t, cache)
}

func sizeTest(t *testing.T, cache ICache) {
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	if cache.Size() != 3 {
		t.Errorf("Expected1 3, got %v", cache.Size())
	}

	cache.Set("key4", "value4")
	if cache.Size() != 3 {
		t.Errorf("Expected2 3, got %v", cache.Size())
	}

	cache.Delete("key4")
	if cache.Size() != 2 {
		t.Errorf("Expected3 2, got %v", cache.Size())
	}
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Expected4 0, got %v", cache.Size())
	}
}

func TestLRUSize(t *testing.T) {
	cache := NewLRUCache(3)
	sizeTest(t, cache)
}

func TestFIFOSize(t *testing.T) {
	cache := NewFIFOCache(3)
	sizeTest(t, cache)
}

func TestLFUSize(t *testing.T) {
	cache := NewLFUCache(3)
	sizeTest(t, cache)
}

func clearTest(t *testing.T, cache ICache) {
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Clear()

	val1 := cache.Get("key1")
	val2 := cache.Get("key2")
	if val1 != nil || val2 != nil {
		t.Errorf("Expected nil, got %v, %v", val1, val2)
	}
}

func TestLRUClear(t *testing.T) {
	cache := NewLRUCache(3)
	clearTest(t, cache)
}

func TestFIFOClear(t *testing.T) {
	cache := NewFIFOCache(3)
	clearTest(t, cache)
}

func TestLFUClear(t *testing.T) {
	cache := NewLFUCache(3)
	clearTest(t, cache)
}

func deleteTest(t *testing.T, cache ICache) {
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Delete("key1")

	val1 := cache.Get("key1")
	if val1 != nil {
		t.Errorf("Expected nil, got %v", val1)
	}
}

func TestLRUDelete(t *testing.T) {
	cache := NewLRUCache(3)
	deleteTest(t, cache)
}

func TestLFUDelete(t *testing.T) {
	cache := NewLFUCache(3)
	deleteTest(t, cache)
}

func TestFIFODelete(t *testing.T) {
	cache := NewFIFOCache(3)
	deleteTest(t, cache)
}

func TestLRUGetSet(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.Set("key1", "value1")
	cache.Set("key3", "value3")

	val1 := cache.Get("key1")
	val2 := cache.Get("key2")
	val3 := cache.Get("key3")

	if val2 != nil || val3 != "value3" || val1 != "value1" {
		t.Errorf("Expected nil, got %v %v %v", val1, val2, val3)
	}
}

func TestLFUGetSet(t *testing.T) {
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
		t.Errorf("Expected nil, got %v", val2)
	}
}

func TestFIFOGetSet(t *testing.T) {
	cache := NewFIFOCache(2)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	val1 := cache.Get("key1")
	val2 := cache.Get("key2")
	val3 := cache.Get("key3")

	if val1 != nil || val2 != "value2" || val3 != "value3" {
		t.Errorf("Expected nil, got %v", val2)
	}
}
