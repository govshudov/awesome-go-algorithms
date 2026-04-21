// Package lrucache provides a generic least-recently-used cache with O(1)
// Get, Put, and eviction. It combines a doubly-linked list from
// container/list with a map keyed by the cache key.
package lrucache

import "container/list"

type entry[K comparable, V any] struct {
	key K
	val V
}

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*list.Element
	order    *list.List
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity < 1 {
		capacity = 1
	}
	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element, capacity),
		order:    list.New(),
	}
}

func (c *Cache[K, V]) Len() int { return c.order.Len() }

func (c *Cache[K, V]) Get(key K) (V, bool) {
	var zero V
	el, ok := c.items[key]
	if !ok {
		return zero, false
	}
	c.order.MoveToFront(el)
	return el.Value.(*entry[K, V]).val, true
}

func (c *Cache[K, V]) Put(key K, val V) {
	if el, ok := c.items[key]; ok {
		c.order.MoveToFront(el)
		el.Value.(*entry[K, V]).val = val
		return
	}
	if c.order.Len() >= c.capacity {
		c.evictOldest()
	}
	el := c.order.PushFront(&entry[K, V]{key: key, val: val})
	c.items[key] = el
}

func (c *Cache[K, V]) evictOldest() {
	el := c.order.Back()
	if el == nil {
		return
	}
	c.order.Remove(el)
	delete(c.items, el.Value.(*entry[K, V]).key)
}
