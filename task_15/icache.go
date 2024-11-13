package main

type ICache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Delete(key string)
	Clear()
	Size() int
	Keys() []string
}

type DataWrapper struct {
	key  string
	data interface{}
}
