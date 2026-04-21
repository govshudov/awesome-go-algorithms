// Package linkedlist provides a generic doubly-linked list with O(1)
// insertion and removal at both ends.
package linkedlist

type node[T any] struct {
	value      T
	prev, next *node[T]
}

type List[T any] struct {
	head, tail *node[T]
	size       int
}

func New[T any]() *List[T] { return &List[T]{} }

func (l *List[T]) Len() int { return l.size }

func (l *List[T]) PushFront(v T) {
	n := &node[T]{value: v, next: l.head}
	if l.head != nil {
		l.head.prev = n
	} else {
		l.tail = n
	}
	l.head = n
	l.size++
}

func (l *List[T]) PushBack(v T) {
	n := &node[T]{value: v, prev: l.tail}
	if l.tail != nil {
		l.tail.next = n
	} else {
		l.head = n
	}
	l.tail = n
	l.size++
}

func (l *List[T]) PopFront() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}
	n := l.head
	l.head = n.next
	if l.head != nil {
		l.head.prev = nil
	} else {
		l.tail = nil
	}
	l.size--
	return n.value, true
}

func (l *List[T]) PopBack() (T, bool) {
	var zero T
	if l.tail == nil {
		return zero, false
	}
	n := l.tail
	l.tail = n.prev
	if l.tail != nil {
		l.tail.next = nil
	} else {
		l.head = nil
	}
	l.size--
	return n.value, true
}

func (l *List[T]) ToSlice() []T {
	out := make([]T, 0, l.size)
	for n := l.head; n != nil; n = n.next {
		out = append(out, n.value)
	}
	return out
}
