// Package queue provides a generic FIFO queue backed by a growable ring buffer.
// Enqueue and Dequeue are amortized O(1).
package queue

type Queue[T any] struct {
	buf        []T
	head, tail int
	size       int
}

func New[T any]() *Queue[T] {
	return &Queue[T]{buf: make([]T, 8)}
}

func (q *Queue[T]) Len() int { return q.size }

func (q *Queue[T]) Enqueue(v T) {
	if q.size == len(q.buf) {
		q.grow()
	}
	q.buf[q.tail] = v
	q.tail = (q.tail + 1) % len(q.buf)
	q.size++
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	v := q.buf[q.head]
	q.buf[q.head] = zero
	q.head = (q.head + 1) % len(q.buf)
	q.size--
	return v, true
}

func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	return q.buf[q.head], true
}

func (q *Queue[T]) grow() {
	newBuf := make([]T, len(q.buf)*2)
	for i := 0; i < q.size; i++ {
		newBuf[i] = q.buf[(q.head+i)%len(q.buf)]
	}
	q.buf = newBuf
	q.head = 0
	q.tail = q.size
}
