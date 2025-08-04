package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/armstrong99/lib-in-memory-go-cache/heap"
	"github.com/armstrong99/lib-in-memory-go-cache/types"
)

type Cache struct {
	data     map[string]*types.CacheItem
	lru      *LRUList
	capacity int
	ttlHeap  types.ISync[types.CacheItem]
	mu       sync.RWMutex
}

func (c *Cache) RemoveItem(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, exists := c.data[key]; exists {
		fmt.Printf("Removing item: %v \n", item.Node)
		c.lru.remove(&item.Node)
		delete(c.data, item.Key)
	}

}

func (c *Cache) Get(key string) *types.CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if item, exists := c.data[key]; exists {
		c.lru.MoveToFront(&c.data[key].Node)
		return item
	}
	return nil
}

func (c *Cache) Set(key string, value string, ttl *time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.data) >= c.capacity {
		node := c.lru.RemoveLRU()
		delete(c.data, node.Key)
	}

	if ttl == nil {
		defaultx := time.Now().Add(30 * time.Second)
		ttl = &defaultx
	}

	cacheItem := &types.CacheItem{
		Key:       key,
		Value:     value,
		ExpiresAt: *ttl,
		Node: types.Node{
			Key:  key,
			Prev: nil,
			Next: nil,
		},
	}
	c.data[key] = cacheItem
	c.lru.insertAfterHead(&cacheItem.Node)
	c.ttlHeap.Insert(*cacheItem)
	// fmt.Printf("Node info: %+v\n", cacheItem.Node)
}

func NewCache(capacity int, ctx context.Context) *Cache {
	c := &Cache{
		data:     make(map[string]*types.CacheItem),
		lru:      NewLRUList(),
		capacity: capacity,
	}
	c.ttlHeap = heap.InitHeapSync(ctx, c.RemoveItem)
	return c
}
