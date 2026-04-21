package queue

import "testing"

func TestQueue_FIFO(t *testing.T) {
	q := New[int]()
	for _, v := range []int{1, 2, 3, 4} {
		q.Enqueue(v)
	}
	for _, w := range []int{1, 2, 3, 4} {
		v, ok := q.Dequeue()
		if !ok || v != w {
			t.Fatalf("Dequeue = (%d, %v), want (%d, true)", v, ok, w)
		}
	}
	if q.Len() != 0 {
		t.Fatalf("Len = %d after draining, want 0", q.Len())
	}
}

func TestQueue_GrowsPastInitialCapacity(t *testing.T) {
	q := New[int]()
	const n = 100
	for i := range n {
		q.Enqueue(i)
	}
	if q.Len() != n {
		t.Fatalf("Len = %d, want %d", q.Len(), n)
	}
	for i := range n {
		v, ok := q.Dequeue()
		if !ok || v != i {
			t.Fatalf("iter %d: Dequeue = (%d, %v), want (%d, true)", i, v, ok, i)
		}
	}
}

func TestQueue_WrapsAround(t *testing.T) {
	q := New[int]()
	for i := range 6 {
		q.Enqueue(i)
	}
	for range 4 {
		q.Dequeue()
	}
	for i := 6; i < 10; i++ {
		q.Enqueue(i)
	}
	want := []int{4, 5, 6, 7, 8, 9}
	for _, w := range want {
		v, _ := q.Dequeue()
		if v != w {
			t.Fatalf("Dequeue = %d, want %d", v, w)
		}
	}
}

func TestQueue_EmptyDequeue(t *testing.T) {
	q := New[string]()
	if v, ok := q.Dequeue(); ok || v != "" {
		t.Fatalf("empty Dequeue = (%q, %v), want (\"\", false)", v, ok)
	}
}
