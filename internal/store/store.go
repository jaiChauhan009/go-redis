package store

import (
	"sync"
	"time"
)

var cache sync.Map

func Get(key string) (string, bool) {
	val, ok := cache.Load(key)
	if !ok {
		return "", false
	}
	return val.(string), true
}

func Set(key string, value string, ttl time.Duration) {
	cache.Store(key, value)
	if ttl > 0 {
		go func() {
			time.Sleep(ttl)
			cache.Delete(key)
		}()
	}
}

func Delete(key string) bool {
	_, ok := cache.LoadAndDelete(key)
	return ok
}
