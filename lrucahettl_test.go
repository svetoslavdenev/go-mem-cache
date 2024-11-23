package cache_test

import (
	"strconv"
	"testing"
	"time"

	cache "github.com/svetoslavdenev/go-mem-cache"
)

func TestLRUCacheWithTTL(t *testing.T) {
	cache := cache.NewLruMemoryCacheWithTtl(3, 2*time.Second)

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

	// Test TTL expiration
	time.Sleep(3 * time.Second) // Wait for keys to expire

	if _, ok := cache.Get("key1"); ok {
		t.Fatalf("expected key1 to be expired, but it was found")
	}

	// Test overwriting keys
	cache.Set("key5", "value5")
	cache.Set("key5", "new_value5")

	if val, ok := cache.Get("key5"); !ok || val.(string) != "new_value5" {
		t.Fatalf("expected key5 to have value 'new_value5', got %v", val)
	}
}

func TestCacheEvictionOrder(t *testing.T) {
	cache := cache.NewLruMemoryCacheWithTtl(2, 5*time.Second)

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

func BenchmarkLRUCacheSet(b *testing.B) {
	cache := cache.NewLruMemoryCacheWithTtl(1000, 5*time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i%1000), i) // Simulate reuse of keys to avoid unbounded growth
	}
}

func BenchmarkLRUCacheGet(b *testing.B) {
	cache := cache.NewLruMemoryCacheWithTtl(1000, 5*time.Second)

	// Pre-fill the cache
	for i := 0; i < 1000; i++ {
		cache.Set(strconv.Itoa(i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(strconv.Itoa(i % 1000)) // Access existing keys
	}
}
