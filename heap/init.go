package heap

import (
	"context"
	"sync"
	"time"

	"github.com/armstrong99/lib-in-memory-go-cache/types"
)

func less(a, b types.CacheItem) bool {
	return a.ExpiresAt.Before(b.ExpiresAt)
}

var (
	coreHeap = NewHeap(less, 100)
	// numOfWorkers = 32
	rMap     = make(map[string]types.CacheItem)
	rMapChan = make(chan uint8, 1)
	rmLocks  sync.RWMutex
	lruMU    sync.Mutex
)

const (
	InsertAsync = 1 << iota
	PeekMinAsync
	ExtractMinAsync
)

// ------ Sync Implementation ------
type HeapSync[T types.CacheItem] struct {
	mutator chan types.Mutator[T]
}

func (hs *HeapSync[T]) Insert(value T) {
	hs.mutator <- types.Mutator[T]{
		Action: InsertAsync,
		Item:   value,
	}
}

func (hs *HeapSync[T]) PeekMin() {
	hs.mutator <- types.Mutator[T]{
		Action: PeekMinAsync,
	}
}

func (hs *HeapSync[T]) ExtractMin() {
	hs.mutator <- types.Mutator[T]{
		Action: ExtractMinAsync,
	}
}

func (hs *HeapSync[T]) Mutator() chan types.Mutator[T] {
	return hs.mutator
}

func UpdateMinMap(item types.CacheItem) {
	rmLocks.Lock()
	rMap["min"] = item
	select {
	case rMapChan <- 1:
	default:
	}
	rmLocks.Unlock()
}

func DeleteMinMap() {
	rmLocks.Lock()
	delete(rMap, "min")
	select {
	case rMapChan <- 1:
	default:
	}
	rmLocks.Unlock()
}

// ---- Export Sync -----
func InitHeapSync(ctx context.Context, delLRUCacheItem func(key string)) types.ISync[types.CacheItem] {

	syncEngine := HeapSync[types.CacheItem]{
		mutator: make(chan types.Mutator[types.CacheItem], 10),
	}

	go func(rCtx context.Context) {
		for {
			select {
			case <-rCtx.Done():
				return

			case dto := <-syncEngine.mutator:
				switch dto.Action {
				case InsertAsync:
					lruMU.Lock()
					if min, err := coreHeap.Insert(dto.Item); err == nil {
						UpdateMinMap(min)
					}
					lruMU.Unlock()
				}
			}
		}
	}(ctx)

	go func(rCtx context.Context) {
		for {

			rmLocks.RLock()
			item, ok := rMap["min"]
			rmLocks.RUnlock()

			if !ok {
				select {

				case <-rCtx.Done():
					return

				case <-time.After(500 * time.Millisecond):
					continue

				case <-rMapChan:
					continue
				}
			}

			now := time.Now()

			if item.ExpiresAt.Before(now) || item.ExpiresAt.Equal(now) {
				lruMU.Lock()
				delLRUCacheItem(item.Key)
				DeleteMinMap()
				if min, err := coreHeap.ExtractMin(); err == nil {
					UpdateMinMap(min)
				}
				lruMU.Unlock()
				continue
			}

			sleepDuration := time.Until(item.ExpiresAt)

			select {

			case <-rCtx.Done():
				return

			case <-time.After(sleepDuration):
				continue

			case <-rMapChan:
				continue

			}
		}
	}(ctx)

	return &syncEngine
}
