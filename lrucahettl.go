package cache

import (
	"container/list"
	"sync"
	"time"
)

// CacheEntry represents an individual cache item.
type CacheEntry struct {
	Value      any
	Expiration int64 // Unix timestamp in nanoseconds
}

// LruMemoryCacheWithTtl represents the memory cache.
type LruMemoryCacheWithTtl struct {
	store        sync.Map                 // Concurrent map for storing cache entries
	keyToElement map[string]*list.Element // Map for tracking list elements by key
	order        *list.List               // Doubly linked list for maintaining key order (LRU)
	maxSize      int                      // Maximum number of items
	ttl          time.Duration            // TTL for each cache item
	mu           sync.Mutex               // Mutex for thread safety of the order list and keyToElement map
}

// NewLruMemoryCacheWithTtl creates a new cache with optional TTL and size limit.
func NewLruMemoryCacheWithTtl(maxSize int, ttl time.Duration) *LruMemoryCacheWithTtl {
	return &LruMemoryCacheWithTtl{
		store:        sync.Map{},
		keyToElement: make(map[string]*list.Element),
		order:        list.New(),
		maxSize:      maxSize,
		ttl:          ttl,
	}
}

// Set adds an item to the cache.
func (mc *LruMemoryCacheWithTtl) Set(key string, value any) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Check if the key already exists
	if elem, found := mc.keyToElement[key]; found {
		// Update the value and expiration
		entry := &CacheEntry{
			Value:      value,
			Expiration: mc.getExpirationTime(),
		}
		mc.store.Store(key, entry)

		// Move the element to the front of the list
		mc.order.MoveToFront(elem)
		return
	}

	// Add a new entry
	entry := &CacheEntry{
		Value:      value,
		Expiration: mc.getExpirationTime(),
	}
	mc.store.Store(key, entry)
	element := mc.order.PushFront(key)
	mc.keyToElement[key] = element

	// Evict the oldest item if the cache exceeds the max size
	if mc.maxSize > 0 && mc.order.Len() > mc.maxSize {
		mc.evictOldest()
	}
}

// Get retrieves an item from the cache.
func (mc *LruMemoryCacheWithTtl) Get(key string) (any, bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Look for the item in the store
	value, ok := mc.store.Load(key)
	if !ok {
		return nil, false
	}

	entry := value.(*CacheEntry)

	// Check for expiration
	if entry.Expiration > 0 && time.Now().UnixNano() > entry.Expiration {
		mc.store.Delete(key)
		mc.removeKey(key)
		return nil, false
	}

	// Move the accessed key to the front of the list
	if elem, found := mc.keyToElement[key]; found {
		mc.order.MoveToFront(elem)
	}

	return entry.Value, true
}

// Delete removes an item from the cache.
func (mc *LruMemoryCacheWithTtl) Delete(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if _, found := mc.store.Load(key); found {
		mc.store.Delete(key)
		mc.removeKey(key)
	}
}

// evictOldest removes the least recently used item from the cache.
func (mc *LruMemoryCacheWithTtl) evictOldest() {
	oldest := mc.order.Back()
	if oldest != nil {
		key := oldest.Value.(string)
		mc.store.Delete(key)
		mc.order.Remove(oldest)
		delete(mc.keyToElement, key)
	}
}

// removeKey removes a key from the linked list and the keyToElement map.
func (mc *LruMemoryCacheWithTtl) removeKey(key string) {
	if elem, found := mc.keyToElement[key]; found {
		mc.order.Remove(elem)
		delete(mc.keyToElement, key)
	}
}

// getExpirationTime calculates the expiration time based on TTL.
func (mc *LruMemoryCacheWithTtl) getExpirationTime() int64 {
	if mc.ttl > 0 {
		return time.Now().Add(mc.ttl).UnixNano()
	}
	return 0
}
