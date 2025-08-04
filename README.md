# 🧠 lib-in-memory-go-cache

A blazing-fast ⚡️, Golang-native in-memory cache library built for performance-critical systems. Featuring heap-based TTL eviction, lock-safe concurrent access, and context-aware goroutines — this library is designed to scale under pressure 🧠💪.

## 🧩 Core Features

- **🚀 Microsecond-Precision TTL** – Eviction is driven by the next soonest expiry using a binary min-heap, not lazy polling.
- **🧵 Goroutine + Context + Channel-Orchestrated Engine** – Combines Go's concurrency primitives to sleep until eviction is truly due, avoiding wasteful cycles.
- **🛡️ Thread-Safe Read/Write Locks** – Efficient use of `sync.RWMutex` for high-read, low-write scenarios.
- **🔄 Heap-Driven Priority Queue** – Automatic min-priority reordering for accurate TTL expiration.
- **🧪 Battle-Tested with Unit Tests** – Lightweight yet robust test coverage to ensure reliability.

## 🛠 Use Cases

- **LRU-style in-memory caches** 🧊
- **Temporary token/session storage for microservices** 🪪
- **Ultra-low-latency, high-throughput systems** (think real-time bidding or IoT data ingestion) ⚙️📡

## 🧠 Designed for Engineers Who...

- Want fine-grained TTL precision without wasting CPU cycles
- Need scalable concurrent cache logic without race conditions
- Believe in idiomatic Go with zero-dependency clarity

An efficient in-memory caching library in Go with support for:

- **LRU (Least Recently Used) eviction**
- **Per-key TTL expiration**
- **Min-Heap based background cleaner** (no need to wait for `Get`/`Set` calls)
- **Optimized idle CPU usage with sleep and wakeup signals**
- **Modular design**, suitable for embedding into distributed services

---

## 📦 Features

- ✅ Constant-time `Get`, `Set`, and `Remove` operations
- ✅ Automatic removal of expired keys using a **heap-based cleaner**
- ✅ Efficient **LRU eviction** policy when capacity is exceeded
- ✅ Optimized **CPU-friendly background cleaner**
- ✅ No external dependencies like Redis
- ✅ Clean and testable design

---

## 🛠️ Installation

```bash
go get github.com/armstrong99/lib-in-memory-go-cache
```

---

## 🚀 Usage

```go
import (
    "context"
    "time"
    "github.com/armstrong99/lib-in-memory-go-cache/cache"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    c := cache.NewCache(3, ctx) // max 3 items

    // Set with TTL
    ttl := time.Now().Add(2 * time.Second)
    c.Set("key1", "value1", &ttl)

    // Set without TTL
    c.Set("key2", "value2", nil)

    val := c.Get("key1")
    if val != nil {
        fmt.Println(val.Value)
    }

    c.RemoveItem("key2")
}
```

---

## 🧠 TTL and Expiry Internals

This cache does **not** wait for `Get()` or `Set()` to trigger key expiry. Instead, it:

- Stores TTL deadlines in a **min-heap**
- Starts a **background goroutine** that pops the earliest expiry
- Removes items from the cache automatically when expired

### ⚡ Optimized for CPU Efficiency

To avoid spinning or busy-waiting when no TTLs are near expiry, the cleaner now:

- Uses **`time.Sleep()`** for precise minimal delay until the next expiry
- Wakes up early via a **channel** if a new item with a closer TTL is inserted
- Sleeps indefinitely when the heap is empty, waking only on signal

This makes the cleaner very lightweight and **non-blocking**, even in large-scale systems.

---

## 📁 Folder Structure

```
lib-in-memory-go-cache/
├── cache/
│   ├── cache.go          // Main cache logic (LRU, Set, Get)
│   ├── lru.go           // LRU logic w/ double linked list
├── heap/
│   ├── init.go          // initialisation folder for the TTL min heap w/ sync logic
│   ├── theap.go         // the actual min heap implementation
├── tests/
│   └── cache_test.go     // All unit tests
├── go.mod
└── README.md
```

---

## ✅ Tests

Located in `tests/cache_test.go`, covering:

### ✔️ Set and Get

```go
func TestCache_SetAndGet(t *testing.T)
```

- Ensures a basic `Set()` followed by `Get()` works as expected.

### ✔️ LRU Eviction

```go
func TestCache_LRUEviction(t *testing.T)
```

- Confirms that least recently used items are evicted when capacity is reached.

### ✔️ TTL Expiry

```go
func TestCache_TTLExpiry(t *testing.T)
```

- Tests that expired items are removed automatically even without explicit access.

### ✔️ Manual Removal

```go
func TestCache_RemoveItem(t *testing.T)
```

- Verifies that a manually removed item is no longer available.

To run all tests:

```bash
go test ./...
```

---

## 📌 How LRU Works

- Doubly linked list keeps track of usage order
- Hashmap provides O(1) access to nodes
- Most recently used node is moved to head
- Least recently used node is removed on eviction

---

## 📌 How TTL Cleaner Works

- TTL deadlines stored in a **min-heap**
- A background goroutine watches the top of the heap
- When the item at the top expires, it’s removed from both heap and cache

### ⏱️ Now optimized for low CPU usage:

- When heap is empty, the cleaner sleeps indefinitely
- When next expiry is far in the future, it sleeps only until needed
- A channel wakes the cleaner early if a new soon-to-expire item is inserted

This balances **precision and performance**, avoiding unnecessary cycles.

---

## 🔮 Future Improvements

- Optional metrics/statistics (hits/misses/evictions)
- Sharded cache for better concurrency
- Custom eviction strategies (LFU, FIFO, etc.)

---

## 👨‍💻 Author

**Armstrong Ndukwe**  
A high-performance, concurrency-loving Go developer.  
🔗 https://www.linkedin.com/in/ndukwearmstrong/

---

## 📜 License

MIT — free for commercial and personal use.
