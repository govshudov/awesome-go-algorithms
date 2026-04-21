// Package heap provides a generic binary heap parameterized by a comparator.
// Pass less(a, b) == (a < b) for a min-heap, or (a > b) for a max-heap.
package heap

type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

func New[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{less: less}
}

func (h *Heap[T]) Len() int { return len(h.data) }

func (h *Heap[T]) Push(v T) {
	h.data = append(h.data, v)
	h.up(len(h.data) - 1)
}

func (h *Heap[T]) Pop() (T, bool) {
	var zero T
	n := len(h.data)
	if n == 0 {
		return zero, false
	}
	top := h.data[0]
	h.data[0] = h.data[n-1]
	h.data[n-1] = zero
	h.data = h.data[:n-1]
	if len(h.data) > 0 {
		h.down(0)
	}
	return top, true
}

func (h *Heap[T]) Peek() (T, bool) {
	var zero T
	if len(h.data) == 0 {
		return zero, false
	}
	return h.data[0], true
}

func (h *Heap[T]) up(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.less(h.data[i], h.data[parent]) {
			return
		}
		h.data[i], h.data[parent] = h.data[parent], h.data[i]
		i = parent
	}
}

func (h *Heap[T]) down(i int) {
	n := len(h.data)
	for {
		left := 2*i + 1
		if left >= n {
			return
		}
		smallest := left
		if right := left + 1; right < n && h.less(h.data[right], h.data[left]) {
			smallest = right
		}
		if !h.less(h.data[smallest], h.data[i]) {
			return
		}
		h.data[i], h.data[smallest] = h.data[smallest], h.data[i]
		i = smallest
	}
}
