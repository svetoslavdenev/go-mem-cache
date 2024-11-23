package cache

import (
	"container/list"
	"sync"
)

// LruMemoryCache represents a simple LRU cache.
type LruMemoryCache struct {
	store        sync.Map                 // Concurrent map for storing cache entries
	keyToElement map[string]*list.Element // Map for tracking list elements by key
	order        *list.List               // Doubly linked list for maintaining key order (LRU)
	maxSize      int                      // Maximum number of items
	mu           sync.Mutex               // Mutex for thread safety of the order list and keyToElement map
}

// NewLruMemoryCache creates a new LRU cache with a specified maximum size.
func NewLruMemoryCache(maxSize int) *LruMemoryCache {
	return &LruMemoryCache{
		store:        sync.Map{},
		keyToElement: make(map[string]*list.Element),
		order:        list.New(),
		maxSize:      maxSize,
	}
}

// Set adds an item to the cache.
func (mc *LruMemoryCache) Set(key string, value any) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Check if the key already exists
	if elem, found := mc.keyToElement[key]; found {
		// Update the value
		mc.store.Store(key, value)

		// Move the element to the front of the list
		mc.order.MoveToFront(elem)
		return
	}

	// Add a new entry
	mc.store.Store(key, value)
	element := mc.order.PushFront(key)
	mc.keyToElement[key] = element

	// Evict the oldest item if the cache exceeds the max size
	if mc.maxSize > 0 && mc.order.Len() > mc.maxSize {
		mc.evictOldest()
	}
}

// Get retrieves an item from the cache.
func (mc *LruMemoryCache) Get(key string) (any, bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Look for the item in the store
	value, ok := mc.store.Load(key)
	if !ok {
		return nil, false
	}

	// Move the accessed key to the front of the list
	if elem, found := mc.keyToElement[key]; found {
		mc.order.MoveToFront(elem)
	}

	return value, true
}

// Delete removes an item from the cache.
func (mc *LruMemoryCache) Delete(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if _, found := mc.store.Load(key); found {
		mc.store.Delete(key)
		mc.removeKey(key)
	}
}

// evictOldest removes the least recently used item from the cache.
func (mc *LruMemoryCache) evictOldest() {
	oldest := mc.order.Back()
	if oldest != nil {
		key := oldest.Value.(string)
		mc.store.Delete(key)
		mc.order.Remove(oldest)
		delete(mc.keyToElement, key)
	}
}

// removeKey removes a key from the linked list and the keyToElement map.
func (mc *LruMemoryCache) removeKey(key string) {
	if elem, found := mc.keyToElement[key]; found {
		mc.order.Remove(elem)
		delete(mc.keyToElement, key)
	}
}
