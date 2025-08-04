package heap

import (
	"fmt"
)

type MinHeap[T any] struct {
	items []T
	less  func(a, b T) bool
}

// for each node at index i, children are at:
//    left: 2*i + 1
//    right: 2*i + 2
//    parent: (i-1)/2

func (h *MinHeap[T]) getParentIndex(i int) int     { return (i - 1) / 2 }
func (h *MinHeap[T]) getLeftChildIndex(i int) int  { return 2*i + 1 }
func (h *MinHeap[T]) getRightChildIndex(i int) int { return 2*i + 2 }

func (h *MinHeap[T]) hasParent(i int) bool     { return h.getParentIndex(i) >= 0 }
func (h *MinHeap[T]) hasLeftChild(i int) bool  { return h.getLeftChildIndex(i) < len(h.items) }
func (h *MinHeap[T]) hasRightChild(i int) bool { return h.getRightChildIndex(i) < len(h.items) }

func (h *MinHeap[T]) parent(i int) T     { return h.items[h.getParentIndex((i))] }
func (h *MinHeap[T]) leftChild(i int) T  { return h.items[h.getLeftChildIndex(i)] }
func (h *MinHeap[T]) rightChild(i int) T { return h.items[h.getRightChildIndex((i))] }

func (h *MinHeap[T]) swap(i1, i2 int) {
	h.items[i1], h.items[i2] = h.items[i2], h.items[i1]
}

func (h *MinHeap[T]) Insert(value T) (T, error) {
	// print("insert ", value, "\n")
	h.items = append(h.items, value)
	h.heapifyUp()
	return h.PeekMin()
}

func (h *MinHeap[T]) heapifyUp() {
	index := len(h.items) - 1

	for h.hasParent(index) && h.less(h.items[index], h.parent(index)) {
		h.swap(h.getParentIndex(index), index)
		index = h.getParentIndex(index)
	}
}

func (h *MinHeap[T]) ExtractMin() (T, error) {
	var tValue T
	if len(h.items) == 0 {
		return tValue, fmt.Errorf("heap is empty")
	}

	h.items[0] = h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	h.heapifyDown()

	return h.PeekMin()
}

func (h *MinHeap[T]) heapifyDown() {
	index := 0

	for h.hasLeftChild(index) {
		smallerChildIndex := h.getLeftChildIndex(index)
		if h.hasRightChild(index) && h.less(h.rightChild(index), h.leftChild(index)) {
			smallerChildIndex = h.getRightChildIndex(index)
		}

		if h.less(h.items[index], h.items[smallerChildIndex]) {
			break
		}
		h.swap(index, smallerChildIndex)
		index = smallerChildIndex
	}
}

func (h *MinHeap[T]) PeekMin() (T, error) {
	var tValue T

	if len(h.items) == 0 {
		return tValue, fmt.Errorf("heap is empty")
	}

	return h.items[0], nil
}

func NewHeap[T any](less func(a, b T) bool, bufferSize int) *MinHeap[T] {
	minHeap := &MinHeap[T]{
		items: make([]T, 0, bufferSize),
		less:  less,
	}

	return minHeap
}
