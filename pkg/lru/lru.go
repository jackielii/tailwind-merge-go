package lru

import (
	"sync"

	"github.com/jackielii/tailwind-merge-go/pkg/cache"
)

type node struct {
	key  string
	val  string
	prev *node
	next *node
}

type LRU struct {
	maxCapacity int
	size        int
	cache       map[string]*node
	head        *node
	tail        *node
	mu          sync.Mutex
}

func (lru *LRU) Get(key string) string {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	n := lru.cache[key]
	if n == nil {
		return ""
	}

	lru.moveToMostRecentlyUsed(n)
	return n.val
}

func (lru *LRU) Set(key, value string) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if lru.maxCapacity <= 0 {
		return
	}

	if n := lru.cache[key]; n != nil {
		n.val = value
		lru.moveToMostRecentlyUsed(n)
		return
	}

	n := &node{key: key, val: value}
	lru.cache[key] = n
	lru.insertRight(n)
	lru.size++

	if lru.size > lru.maxCapacity {
		evicted := lru.tail.next
		lru.remove(evicted)
		delete(lru.cache, evicted.key)
		lru.size--
	}
}

func (lru *LRU) insertRight(n *node) {
	prev := lru.head.prev
	prev.next = n
	n.prev = prev
	n.next = lru.head
	lru.head.prev = n
}

func (lru *LRU) moveToMostRecentlyUsed(n *node) {
	lru.remove(n)
	lru.insertRight(n)
}

func (lru *LRU) remove(n *node) {
	prev := n.prev
	nxt := n.next
	prev.next = nxt
	nxt.prev = prev
	n.prev = nil
	n.next = nil
}

func Make(maxCapacity int) cache.ICache {
	head := &node{}
	tail := &node{}
	tail.next = head
	head.prev = tail
	return &LRU{
		maxCapacity: maxCapacity,
		size:        0,
		cache:       make(map[string]*node),
		head:        head,
		tail:        tail,
	}
}
