package tests

import (
	"context"
	"testing"
	"time"

	"github.com/armstrong99/lib-in-memory-go-cache/cache"
)

func TestCache_SetAndGet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lru := cache.NewCache(3, ctx)
	lru.Set("fruit", "apple", nil)

	item := lru.Get("fruit")

	if item == nil || item.Value != "apple" {
		t.Errorf("Expected -%s-, got %v", "apple", item)
	}
}

func TestCache_LRUEviction(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	lru := cache.NewCache(2, ctx)
	lru.Set("a", "A", nil)
	lru.Set("b", "B", nil)
	lru.Get("a")           // make "a" recently used
	lru.Set("c", "C", nil) // should evict "b"

	if lru.Get("b") != nil {
		t.Errorf("Expected 'b' to be evicted")
	}

	if lru.Get("a") == nil || lru.Get("c") == nil {
		t.Errorf("Expected 'a' and 'c' to be present")
	}
}

func TestCache_TTLExpiry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lru := cache.NewCache(2, ctx)

	ttl := time.Now().Add(100 * time.Millisecond)
	lru.Set("expiring", "soon", &ttl)

	time.Sleep(200 * time.Millisecond)
	item := lru.Get("expiring")
	if item != nil {
		t.Errorf("Expected 'expiring' to have expired")
	}
}

func TestCache_RemoveItem(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	lru := cache.NewCache(2, ctx)

	lru.Set("cake", "sweet", nil)
	lru.RemoveItem("cake")

	if lru.Get("cake") != nil {
		t.Errorf("Expected 'cake' to be removed")
	}
}
