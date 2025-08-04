package types

import "time"

type CacheItem struct {
	Key       string
	Value     string
	ExpiresAt time.Time
	Node      Node
}

// ---- sync variables -----
type HChanAction uint64

type Mutator[T CacheItem] struct {
	Action HChanAction
	Item   T
}

// ------ Sync Interface ------
type ISync[T CacheItem] interface {
	Insert(value T)
	PeekMin()
	ExtractMin()
	Mutator() chan Mutator[T]
}

type Node struct {
	Key  string
	Prev *Node
	Next *Node
}
