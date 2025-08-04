package cache

import "github.com/armstrong99/lib-in-memory-go-cache/types"

type LRUList struct {
	Head *types.Node
	Tail *types.Node
}

func NewLRUList() *LRUList {
	head := &types.Node{}
	tail := &types.Node{}
	head.Next = tail
	tail.Prev = head
	return &LRUList{Head: head, Tail: tail}
}

func (l *LRUList) MoveToFront(n *types.Node) {
	l.remove(n)
	l.insertAfterHead(n)
}

func (l *LRUList) insertAfterHead(n *types.Node) {
	n.Next = l.Head.Next
	n.Prev = l.Head
	l.Head.Next.Prev = n
	l.Head.Next = n
}

func (l *LRUList) remove(n *types.Node) {
	n.Prev.Next = n.Next
	n.Next.Prev = n.Prev
	n.Next = nil
	n.Prev = nil
}

func (l *LRUList) RemoveLRU() *types.Node {
	if l.Tail.Prev == l.Head {
		return nil
	}

	node := l.Tail.Prev
	l.remove(node)
	return node
}
