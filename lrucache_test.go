package cache_test

import (
	"strconv"
	"testing"

	cache "github.com/svetoslavdenev/go-mem-cache"
)

func TestSimpleLRUCache(t *testing.T) {
	cache := cache.NewLruMemoryCache(3) // Create an LRU cache with a max size of 3

	// Test Set and Get
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	if val, ok := cache.Get("key1"); !ok || val.(string) != "value1" {
		t.Fatalf("expected key1 to have value 'value1', got %v", val)
	}

	// Test LRU eviction
	cache.Set("key4", "value4") // Evicts key2 (oldest)

	if _, ok := cache.Get("key2"); ok {
		t.Fatalf("expected key2 to be evicted, but it was found")
	}

	// Test overwriting keys
	cache.Set("key3", "new_value3")

	if val, ok := cache.Get("key3"); !ok || val.(string) != "new_value3" {
		t.Fatalf("expected key3 to have value 'new_value3', got %v", val)
	}
}

func TestSimpleLRUEvictionOrder(t *testing.T) {
	cache := cache.NewLruMemoryCache(2) // Create an LRU cache with a max size of 2

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Get("key1")           // Access key1 to make it recently used
	cache.Set("key3", "value3") // Evicts key2 (oldest, least recently used)

	if _, ok := cache.Get("key2"); ok {
		t.Fatalf("expected key2 to be evicted, but it was found")
	}

	if val, ok := cache.Get("key1"); !ok || val.(string) != "value1" {
		t.Fatalf("expected key1 to be present with value 'value1', got %v", val)
	}
}

func TestSimpleLRUDelete(t *testing.T) {
	cache := cache.NewLruMemoryCache(3) // Create an LRU cache with a max size of 3

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Delete("key1") // Remove key1

	if _, ok := cache.Get("key1"); ok {
		t.Fatalf("expected key1 to be deleted, but it was found")
	}

	if val, ok := cache.Get("key2"); !ok || val.(string) != "value2" {
		t.Fatalf("expected key2 to be present with value 'value2', got %v", val)
	}
}

func BenchmarkSimpleLRUCacheSet(b *testing.B) {
	cache := cache.NewLruMemoryCache(1000) // Create an LRU cache with a max size of 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i%1000), i) // Simulate reuse of keys to avoid unbounded growth
	}
}

func BenchmarkSimpleLRUCacheGet(b *testing.B) {
	cache := cache.NewLruMemoryCache(1000) // Create an LRU cache with a max size of 1000

	// Pre-fill the cache
	for i := 0; i < 1000; i++ {
		cache.Set(strconv.Itoa(i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(strconv.Itoa(i % 1000)) // Access existing keys
	}
}
