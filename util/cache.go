package util

import (
	"sync"
	"time"
)

// A time key value store
// https://github.com/patrickmn/go-cache
// https://www.alexedwards.net/blog/implementing-an-in-memory-cache-in-go
// https://dev.to/ernesto27/key-value-store-in-golang-52h1

type ValueStore struct {
	mu         sync.RWMutex
	store      map[string]string
	expiration map[string]int64 // Stores expiration times as Unix timestamps, 0 for no expiration
}

func NewValueStore(cleanupInterval time.Duration) *ValueStore {
	vs := &ValueStore{
		store:      make(map[string]string),
		expiration: make(map[string]int64),
	}
	go vs.startCleanup(cleanupInterval)
	return vs
}

// we mark zero as non expirary
func (kv *ValueStore) Set(key, value string, ttl time.Duration) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = value
	if ttl > 0 {
		kv.expiration[key] = time.Now().Add(ttl).UnixNano()
	} else {
		kv.expiration[key] = 0
	}
}

func (kv *ValueStore) Get(key string) (string, bool) {
	kv.mu.RLock()
	exp, ok := kv.expiration[key]
	kv.mu.RUnlock()

	// If the key exists and is expired
	if ok && exp > 0 && time.Now().UnixNano() > exp {
		// Clean up the expired key
		kv.Delete(key)
		return "", false
	}

	kv.mu.RLock()
	value, exists := kv.store[key]
	kv.mu.RUnlock()
	return value, exists
}

func (kv *ValueStore) Delete(key string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.store, key)
	delete(kv.expiration, key)
}

func (kv *ValueStore) startCleanup(interval time.Duration) {
	for {
		time.Sleep(interval)
		now := time.Now().UnixNano()
		kv.mu.Lock()
		for key, exp := range kv.expiration {
			if exp > 0 && now > exp {
				delete(kv.store, key)
				delete(kv.expiration, key)
			}
		}
		kv.mu.Unlock()
	}
}