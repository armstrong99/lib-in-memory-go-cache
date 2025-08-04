package main

import (
	"context"
	"log"
	"os"
	"runtime/trace"
	"time"

	"github.com/armstrong99/lib-in-memory-go-cache/cache"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace file: %v", err)
	}
	defer f.Close()

	trace.Start(f)
	defer trace.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lru := cache.NewCache(4, ctx)

	least := time.Now().Add(3 * time.Second)

	lru.Set("apple", "A sweet red fruit", &least)
	lru.Set("cake", "like bread but sweeter", nil)
	lru.Set("car", "An automobile for transportation", nil)
	// lru.Set("girlfriend", "we are dating", nil)

	// lru.RemoveItem("cake")

	time.Sleep(10 * time.Second)

	item := lru.Get("apple")
	if item != nil {
		println("Found item:", item.Key, ": ", item.Value, "\n")
	} else {
		println("Item not found")
	}

}
